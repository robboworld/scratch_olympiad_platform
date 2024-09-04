package http

import (
	"encoding/json"
	"github.com/skinnykaen/rpa_clone/internal/consts"
	"github.com/skinnykaen/rpa_clone/internal/models"
	"github.com/skinnykaen/rpa_clone/internal/services"
	"github.com/skinnykaen/rpa_clone/pkg/logger"
	"io"
	"net/http"
	"strconv"
)

type ProjectHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type ProjectHandlerImpl struct {
	loggers        logger.Loggers
	projectService services.ProjectService
}

func (p ProjectHandlerImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/project" {
		switch r.Method {
		case http.MethodGet:
			projectId := r.URL.Query().Get("id")
			atoi, err := strconv.Atoi(projectId)
			if err != nil {
				p.loggers.Err.Printf("%s", err.Error())
				http.Error(w, consts.ErrAtoi, http.StatusBadRequest)
				return
			}
			project, err := p.projectService.GetProjectById(uint(atoi), r.Context().Value(consts.KeyId).(uint), r.Context().Value(consts.KeyRole).(models.Role))
			if err != nil {
				p.loggers.Err.Printf("%s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(project.Json))
		case http.MethodPost:
			dataBytes, err := io.ReadAll(r.Body)
			if err != nil {
				p.loggers.Err.Printf("%s", err.Error())
				http.Error(w, "incorrect json body", http.StatusBadRequest)
				return
			}
			projectId := r.URL.Query().Get("id")
			atoi, err := strconv.Atoi(projectId)
			if err != nil {
				http.Error(w, consts.ErrAtoi, http.StatusBadRequest)
				return
			}
			project := models.ProjectCore{}
			project.ID = uint(atoi)
			project.Json = string(dataBytes)
			_, err = p.projectService.UpdateProject(project)
			if err != nil {
				p.loggers.Err.Printf("%s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jData, err := json.Marshal(map[string]interface{}{
				"Id": project.ID,
			})
			if err != nil {
				p.loggers.Err.Printf("%s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jData)
		default:
			http.Error(w, "not allowed method", http.StatusMethodNotAllowed)
			return
		}
	}
}
