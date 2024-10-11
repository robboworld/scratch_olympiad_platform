package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type ApplicationAPI interface {
	CreateApplication(application models.ApplicationCore) error
	ExportAllApplications(applications []models.ApplicationCore) error
}

type ApplicationAPIImpl struct {
}

func NewApplicationAPIImpl() ApplicationAPIImpl {
	return ApplicationAPIImpl{}
}

type ApplicationPayloadHTTP struct {
	Author                        *models.UserHTTP `json:"author"`
	Nomination                    string           `json:"nomination"`
	AlgorithmicTaskLink           string           `json:"algorithmicTaskLink"`
	AlgorithmicTaskFile           string           `json:"algorithmicTaskFile"`
	CreativeTaskLink              string           `json:"creativeTaskLink"`
	CreativeTaskFile              string           `json:"creativeTaskFile"`
	EngineeringTaskFile           string           `json:"engineeringTaskFile"`
	EngineeringTaskCloudLink      string           `json:"engineeringTaskCloudLink"`
	EngineeringTaskVideo          string           `json:"engineeringTaskVideo"`
	EngineeringTaskVideoCloudLink string           `json:"engineeringTaskVideoCloudLink"`
	Note                          string           `json:"note"`
}

func (a ApplicationAPIImpl) CreateApplication(application models.ApplicationCore) error {
	authorHttp := &models.UserHTTP{}
	authorHttp.FromCore(application.Author)
	applicationPayload := ApplicationPayloadHTTP{
		Author:                        authorHttp,
		Nomination:                    application.Nomination,
		AlgorithmicTaskLink:           application.AlgorithmicTaskLink,
		AlgorithmicTaskFile:           application.AlgorithmicTaskFile,
		CreativeTaskLink:              application.CreativeTaskLink,
		CreativeTaskFile:              application.CreativeTaskFile,
		EngineeringTaskFile:           application.EngineeringTaskFile,
		EngineeringTaskCloudLink:      application.EngineeringTaskCloudLink,
		EngineeringTaskVideo:          application.EngineeringTaskVideo,
		EngineeringTaskVideoCloudLink: application.EngineeringTaskVideoCloudLink,
		Note:                          application.Note,
	}

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

type ApplicationHTTPListPayload struct {
	Applications []ApplicationPayloadHTTP `json:"applications"`
}

func (a ApplicationAPIImpl) ExportAllApplications(applications []models.ApplicationCore) error {
	var applicationsPayload []ApplicationPayloadHTTP
	for _, application := range applications {
		fmt.Println(application.Author)
		authorHttp := &models.UserHTTP{}
		authorHttp.FromCore(application.Author)
		applicationPayload := ApplicationPayloadHTTP{
			Author:                        authorHttp,
			Nomination:                    application.Nomination,
			AlgorithmicTaskLink:           application.AlgorithmicTaskLink,
			AlgorithmicTaskFile:           application.AlgorithmicTaskFile,
			CreativeTaskLink:              application.CreativeTaskLink,
			CreativeTaskFile:              application.CreativeTaskFile,
			EngineeringTaskFile:           application.EngineeringTaskFile,
			EngineeringTaskCloudLink:      application.EngineeringTaskCloudLink,
			EngineeringTaskVideo:          application.EngineeringTaskVideo,
			EngineeringTaskVideoCloudLink: application.EngineeringTaskVideoCloudLink,
			Note:                          application.Note,
		}
		applicationsPayload = append(applicationsPayload, applicationPayload)
	}

	applicationsListPayload := ApplicationHTTPListPayload{
		Applications: applicationsPayload,
	}
	applicationsListPayloadBytes, err := json.Marshal(applicationsListPayload)
	if err != nil {
		return utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	urlAddr := viper.GetString("api_urls.export_all_applications")
	request, err := http.NewRequest("POST", urlAddr, bytes.NewBuffer(applicationsListPayloadBytes))
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
