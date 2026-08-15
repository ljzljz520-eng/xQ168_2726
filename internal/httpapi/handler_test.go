package httpapi

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"craftmaterials/internal/app"
	"craftmaterials/internal/fixture"
)

func TestRecommendUpload(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("csv", "learners.csv")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	if _, err := part.Write([]byte(fixture.LearnerCSV)); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	service := app.NewService(fixture.Courses(), fixture.Materials(), fixture.ExpectedMaterials())
	request := httptest.NewRequest(http.MethodPost, "/recommend?limit=2", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	NewHandler(service).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var decoded app.Response
	if err := json.Unmarshal(response.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if len(decoded.Results) != 2 || len(decoded.Results[0].Recommendations) != 2 {
		t.Fatalf("results = %#v", decoded.Results)
	}
	if decoded.Results[0].Recommendations[0].Reasons[0] == "" {
		t.Fatalf("reasons = %#v", decoded.Results[0].Recommendations[0].Reasons)
	}
}
