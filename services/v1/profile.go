package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/buonotti/bus-stats-api/models"
	"github.com/buonotti/bus-stats-api/util"
)

type SaveUserProfileResponse struct {
	Result string `json:"result"`
}

type GetUserProfileResponse struct {
	FileName string `json:"file_name"`
}

func SaveUserProfile(userToken string, userId models.UserId, formFile *multipart.FileHeader) (SaveUserProfileResponse, error, int) {
	file, err := formFile.Open()
	if err != nil {
		return SaveUserProfileResponse{}, fmt.Errorf("could not open sent file"), http.StatusInternalServerError
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			// Ignored
		}
	}(file)
	uidStr, err := util.ExtractTokenIdString(userToken)
	uid := models.UserId(uidStr)
	if uid != userId {
		return SaveUserProfileResponse{}, fmt.Errorf("user id and token id do not match"), http.StatusUnauthorized
	}

	fileNameSplit := strings.Split(formFile.Filename, ".")
	ext := fileNameSplit[len(fileNameSplit)-2]

	var fileContent = make([]byte, 10_000_000)
	_, err = file.Read(fileContent)
	if err != nil {
		return SaveUserProfileResponse{}, fmt.Errorf("could not read file content"), http.StatusInternalServerError
	}
	fileName := util.BuildFileName(string(uid), ext)
	util.FsLogger.Infof("saving file %s", fileName)
	err = ioutil.WriteFile(fileName, fileContent, os.ModePerm)
	if err != nil {
		return SaveUserProfileResponse{}, fmt.Errorf("could not save file to disk"), http.StatusInternalServerError
	}

	_, err = util.RestClient.R().
		SetBody(util.Query("UPDATE user:? SET image.name = ?, image.type = ?", uid, uid, ext)).
		Post(util.DatabaseUrl())

	return SaveUserProfileResponse{
		Result: "OK",
	}, nil, http.StatusOK
}

func GetUserProfile(userId models.UserId) (GetUserProfileResponse, error, int) {
	selectResponse, err := util.RestClient.R().
		SetBody("SELECT * FROM user:" + userId).
		Post(util.DatabaseUrl())

	if err != nil {
		return GetUserProfileResponse{}, fmt.Errorf("could not get user information"), http.StatusInternalServerError
	}

	var selectUserResponse models.UserSelectResult
	responseString := util.FormatResponseString(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		return GetUserProfileResponse{}, fmt.Errorf("unexpected database response for userSelectResult %s", err.Error()), http.StatusBadRequest
	}

	if len(selectUserResponse.Result) < 1 {
		return GetUserProfileResponse{}, fmt.Errorf("no user with such id present"), http.StatusInternalServerError
	}

	imageData := selectUserResponse.Result[0].Image
	fileName := util.BuildFileName(imageData.Name, imageData.Type)

	return GetUserProfileResponse{FileName: fileName}, nil, http.StatusOK
}
