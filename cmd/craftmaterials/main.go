package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"craftmaterials/internal/app"
	"craftmaterials/internal/fixture"
	"craftmaterials/internal/httpapi"
)

func main() {
	inputPath := flag.String("input", "", "path to learner CSV; uses the fixture when empty")
	listenAddress := flag.String("listen", "", "HTTP listen address, for example :8080")
	limit := flag.Int("limit", 3, "maximum recommendations per learner")
	flag.Parse()

	service := app.NewService(fixture.Courses(), fixture.Materials(), fixture.ExpectedMaterials())
	if *listenAddress != "" {
		if err := http.ListenAndServe(*listenAddress, httpapi.NewHandler(service)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	reader, closeReader, err := openInput(*inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if closeReader != nil {
		defer closeReader()
	}
	response, err := service.FromCSV(reader, *limit)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(response); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func openInput(path string) (io.Reader, func() error, error) {
	if strings.TrimSpace(path) == "" {
		return strings.NewReader(fixture.LearnerCSV), nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open learner CSV: %w", err)
	}
	return file, file.Close, nil
}
