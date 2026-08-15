package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"craftmaterials/internal/app"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) http.Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Path == "/health" {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
		return
	}
	if r.Method != http.MethodPost || r.URL.Path != "/recommend" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "route not found"})
		return
	}

	limit, err := parseLimit(r.URL.Query().Get("limit"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	reader, closeReader, err := csvReader(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if closeReader != nil {
		defer closeReader()
	}

	response, err := h.service.FromCSV(reader, limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, response)
}

func csvReader(r *http.Request) (io.Reader, func() error, error) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid Content-Type: %w", err)
	}
	if mediaType == "text/csv" || mediaType == "application/csv" {
		return io.LimitReader(r.Body, 1<<20), r.Body.Close, nil
	}
	if mediaType != "multipart/form-data" {
		return nil, nil, fmt.Errorf("Content-Type must be text/csv or multipart/form-data")
	}
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, nil, fmt.Errorf("read multipart body: %w", err)
	}
	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			return nil, nil, fmt.Errorf("multipart field csv is required")
		}
		if nextErr != nil {
			return nil, nil, fmt.Errorf("read multipart part: %w", nextErr)
		}
		if part.FormName() == "csv" {
			return io.LimitReader(part, 1<<20), part.Close, nil
		}
		if err := part.Close(); err != nil {
			return nil, nil, fmt.Errorf("close multipart part: %w", err)
		}
	}
}

func parseLimit(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 3, nil
	}
	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 || limit > 20 {
		return 0, fmt.Errorf("limit must be between 1 and 20")
	}
	return limit, nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
