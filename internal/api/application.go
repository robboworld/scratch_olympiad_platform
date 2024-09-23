package api

import (
	"bytes"
	"encoding/json"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type ApplicationAPI interface {
	CreateApplication(application models.ApplicationCore) error
}

type ApplicationAPIImpl struct {
}

func NewApplicationAPIImpl() ApplicationAPIImpl {
	return ApplicationAPIImpl{}
}

func (a ApplicationAPIImpl) CreateApplication(application models.ApplicationCore) error {
	applicationPayload := models.ApplicationPayloadHTTP{}
	applicationPayload.FromCore(application)
	applicationPayloadBytes, err := json.Marshal(applicationPayload)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	urlAddr := viper.GetString("api_urls.create_application")
	request, err := http.NewRequest("POST", urlAddr, bytes.NewBuffer(applicationPayloadBytes))
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	request.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	if response.StatusCode != http.StatusOK {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: string(body),
		}
	}
	return nil
}
