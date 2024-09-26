package services

import (
	"github.com/robboworld/scratch_olympiad_platform/internal/consts"
	"github.com/robboworld/scratch_olympiad_platform/internal/gateways"
	"github.com/robboworld/scratch_olympiad_platform/internal/models"
	"github.com/robboworld/scratch_olympiad_platform/pkg/utils"
	"github.com/spf13/viper"
	"github.com/thanhpk/randstr"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SolutionService interface {
	CreateSolution(userID uint, fileHeader *multipart.FileHeader) (string, error)
	DownloadSolution(filename, accessLink string) ([]byte, error)
}

type SolutionServiceImpl struct {
	solutionGateway gateways.SolutionGateway
	userGateway     gateways.UserGateway
}

func (s SolutionServiceImpl) CreateSolution(userID uint, fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size >= consts.MaxSolutionFileSize {
		return "", utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "File is too large. The maximum allowed size is 100 MB.",
		}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	defer file.Close()

	filenameParts := strings.Split(fileHeader.Filename, ".")
	if len(filenameParts) != 2 {
		return "", utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "incorrect filename",
		}
	}

	allowedFormats := map[string]bool{"mp4": true, "sb3": true}
	if !allowedFormats[filenameParts[1]] {
		return "", utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "incorrect file format",
		}
	}

	user, err := s.userGateway.GetUserById(userID)
	if err != nil {
		return "", err
	}
	filenamePrefix := strings.Replace(user.FullName, " ", "_", -1)
	tempFile, err := os.CreateTemp("./internal/tmp_upload", filenamePrefix+"_*."+filenameParts[1])
	if err != nil {
		return "", utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "failed to create temporary file",
		}
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "failed to write file",
		}
	}

	err = os.Chmod(tempFile.Name(), 0644)
	if err != nil {
		return "", utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: "failed to chmod file",
		}
	}

	accessLink := randstr.String(12)
	accessLinkHash := utils.GetHashString(accessLink)
	solution := models.SolutionCore{
		Name:       filepath.Base(tempFile.Name()),
		AccessLink: accessLinkHash,
	}

	newSolution, err := s.solutionGateway.CreateSolution(solution)
	if err != nil {
		return "", err
	}

	fileLink := viper.GetString("http_server_address") + "/solution/download?filename=" + newSolution.Name +
		"&access_link=" + accessLink

	return fileLink, nil
}

func (s SolutionServiceImpl) DownloadSolution(filename, accessLink string) ([]byte, error) {
	solution, err := s.solutionGateway.GetSolutionByName(filename)
	if err != nil {
		return nil, err
	}

	accessLinkHash := utils.GetHashString(accessLink)
	if accessLinkHash != solution.AccessLink {
		return nil, utils.ResponseError{
			Code:    http.StatusForbidden,
			Message: consts.ErrAccessDenied,
		}
	}

	filePath := "./internal/tmp_upload/" + filename
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []byte{}, utils.ResponseError{
			Code:    http.StatusBadRequest,
			Message: "file not found",
		}
	}
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, utils.ResponseError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return buf, nil
}
