package http

import (
	"encoding/json"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"io"
	"net/http"
	"os"
	"strings"
)

type AvatarHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type AvatarHandlerImpl struct {
	loggers logger.Loggers
}

func (a AvatarHandlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/avatar" {
		switch r.Method {
		case http.MethodPost:
			// TODO rm to avatar service
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			file, header, err := r.FormFile("file")
			if err != nil {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			filename := strings.Split(header.Filename, ".")
			if len(filename) != 2 {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, "incorrect filename", http.StatusBadRequest)
				return
			}
			if filename[1] != "png" && filename[1] != "jpg" {
				a.loggers.Err.Printf("%s", "incorrect file format")
				http.Error(w, "incorrect file format", http.StatusBadRequest)
				return
			}
			tempFile, err := os.CreateTemp("./internal/tmp_upload", "upload-*"+"."+filename[1])
			if err != nil {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, "failed to create temporary file", http.StatusInternalServerError)
				return
			}
			defer tempFile.Close()

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, "failed to read file", http.StatusInternalServerError)
				return
			}
			tempFile.Write(fileBytes)

			err = os.Chmod(tempFile.Name(), 0644)
			if err != nil {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, "failed to chmod file", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			jData, err := json.Marshal(map[string]interface{}{
				"filename": strings.Replace(tempFile.Name(), "tmp_upload/", "", 99),
			})
			if err != nil {
				a.loggers.Err.Printf("%s", err.Error())
				http.Error(w, "failed to marshal http response", http.StatusInternalServerError)
				return
			}
			w.Write(jData)
		}
	}
}
