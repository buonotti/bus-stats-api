package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/buonotti/bus-stats-api/surreal"
	"github.com/buonotti/bus-stats-api/util"
)

type SaveUserProfileResponse struct {
	Result string `json:"result"`
}

type GetUserProfileResponse struct {
	FileName string `json:"file_name"`
	FileData string `json:"file_data"`
	FileType string `json:"file_type"`
}

func SaveUserProfile(userId models.UserId, formFile *multipart.FileHeader) (SaveUserProfileResponse, int, error) {
	file, err := formFile.Open()
	if err != nil {
		logging.FsLogger.Error(err)
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			logging.FsLogger.Error(err)
		}
	}(file)

	fileNameSplit := strings.Split(formFile.Filename, ".")
	ext := fileNameSplit[len(fileNameSplit)-1]

	var fileContent = make([]byte, 10_000_000)
	_, err = file.Read(fileContent)
	if err != nil {
		logging.FsLogger.Error(err)
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}

	fileName := util.FileName(string(userId), ext)
	logging.FsLogger.Infof("saving file %s", fileName)
	err = os.WriteFile(util.FileName(string(userId), ext), fileContent, 0644)
	if err != nil {
		logging.FsLogger.Error(err)
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}

	_, err = surreal.Query("UPDATE user:? SET image.name = ?, image.type = ?", userId, string(userId), ext)

	if err != nil {
		logging.DbLogger.Error(err)
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}

	return SaveUserProfileResponse{
		Result: "OK",
	}, http.StatusOK, nil
}

func GetUserProfile(userId models.UserId) (GetUserProfileResponse, int, error) {
	selectResponse, err := surreal.Query("SELECT * FROM user:" + string(userId))

	if err != nil {
		logging.DbLogger.Error(err)
		return GetUserProfileResponse{}, http.StatusUnauthorized, services.CredentialError
	}

	var selectUserResponse models.UserSelectResult
	responseString := surreal.FormatResponse(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		logging.ApiLogger.Error(err)
		return GetUserProfileResponse{}, http.StatusBadRequest, services.FormatError
	}

	if len(selectUserResponse.Result) < 1 {
		return GetUserProfileResponse{}, http.StatusNotFound, fmt.Errorf("user not found")
	}

	if selectUserResponse.Result[0].Image.Name == "" {
		return GetUserProfileResponse{}, http.StatusNotFound, fmt.Errorf("no image found")
	}

	imageData := selectUserResponse.Result[0].Image
	fileName := util.FileName(imageData.Name, imageData.Type)

	fileData, err := os.ReadFile(fileName)
	if err != nil {
		logging.FsLogger.Error(err)
		return GetUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}
	base64data := base64.RawURLEncoding.EncodeToString(fileData)

	return GetUserProfileResponse{
		FileName: fileName,
		FileData: base64data,
		FileType: fmt.Sprintf("image/%s", imageData.Type),
	}, http.StatusOK, nil
}
