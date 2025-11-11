package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sprint6/internal/service"
	"time"
)

func Ind(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed"))

		return
	}

	http.ServeFile(w, r, "index.html")
}

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("method not allowed"))
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("parse form error: %v", err), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, fmt.Sprintf("read form file error: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("read file error: %v", err), http.StatusInternalServerError)
		return
	}

	converted, err := service.ConvertAuto(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("convert error: %v", err), http.StatusInternalServerError)
		return
	}

	ts := time.Now().UTC().Format("2006-01-02T15-04-05Z")
	ext := filepath.Ext(header.Filename)
	outName := fmt.Sprintf("%s%s", ts, ext)

	if err := os.WriteFile(outName, []byte(converted), 0o644); err != nil {
		http.Error(w, fmt.Sprintf("write result file error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(converted))
}
