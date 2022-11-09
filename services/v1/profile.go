package v1

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/buonotti/bus-stats-api/util"
)

type SaveUserProfileResponse struct {
	Result string `json:"result"`
}

type GetUserProfileResponse struct {
	FileName string `json:"file_name"`
}

func SaveUserProfile(userId models.UserId, formFile *multipart.FileHeader) (SaveUserProfileResponse, int, error) {
	file, err := formFile.Open()
	if err != nil {
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	fileNameSplit := strings.Split(formFile.Filename, ".")
	ext := fileNameSplit[len(fileNameSplit)-2]

	var fileContent = make([]byte, 10_000_000)
	_, err = file.Read(fileContent)
	if err != nil {
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}
	fileName := util.FileName(string(userId), ext)
	util.FsLogger.Infof("saving file %s", fileName)
	err = ioutil.WriteFile(fileName, fileContent, os.ModePerm)
	if err != nil {
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}

	_, err = util.RestClient.R().
		SetBody(util.Query("UPDATE user:? SET image.name = ?, image.type = ?", userId, userId, ext)).
		Post(util.DatabaseUrl())

	if err != nil {
		return SaveUserProfileResponse{}, http.StatusBadRequest, services.FileError
	}

	return SaveUserProfileResponse{
		Result: "OK",
	}, http.StatusOK, nil
}

func GetUserProfile(userId models.UserId) (GetUserProfileResponse, int, error) {
	selectResponse, err := util.RestClient.R().
		SetBody("SELECT * FROM user:" + userId).
		Post(util.DatabaseUrl())

	if err != nil {
		return GetUserProfileResponse{}, http.StatusUnauthorized, services.CredentialError
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		return GetUserProfileResponse{}, http.StatusBadRequest, services.FormatError
	}

	if len(selectUserResponse.Result) < 1 {
		return GetUserProfileResponse{}, http.StatusUnauthorized, services.CredentialError
	}

	imageData := selectUserResponse.Result[0].Image
	fileName := util.FileName(imageData.Name, imageData.Type)

	return GetUserProfileResponse{FileName: fileName}, http.StatusOK, nil
}
