package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/buonotti/bus-stats-api/errors"
	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/models"
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
		return SaveUserProfileResponse{}, http.StatusBadRequest, errors.CannotReadFileError.New("cannot read form file")
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
		return SaveUserProfileResponse{}, http.StatusBadRequest, errors.CannotReadFileError.New("cannot read form file")
	}

	mimeType := http.DetectContentType(fileContent)
	if mimeType != "image/jpeg" && mimeType != "image/png" && mimeType != "image/svg+xml" {
		return SaveUserProfileResponse{}, http.StatusBadRequest, errors.MimeTypeError.New("invalid file type")
	}

	if len(fileContent) > 10_000_000 {
		return SaveUserProfileResponse{}, http.StatusBadRequest, errors.FileSizeError.New("file too big")
	}

	fileName := util.FileName(string(userId))
	logging.FsLogger.Infof("saving file %s", fileName)
	err = os.WriteFile(util.FileName(string(userId)), fileContent, 0644)
	if err != nil {
		logging.FsLogger.Error(err)
		return SaveUserProfileResponse{}, http.StatusBadRequest, errors.CannotWriteFileError.New("cannot write file to disk")
	}

	_, err = surreal.Query("UPDATE user:? SET image.name = ?, image.type = ?", userId, string(userId), ext)

	if err != nil {
		logging.DbLogger.Error(err)
		return SaveUserProfileResponse{}, http.StatusBadRequest, errors.SurrealQueryError.New("cannot update user image")
	}

	return SaveUserProfileResponse{
		Result: "OK",
	}, http.StatusOK, nil
}

func GetUserProfile(userId models.UserId) (GetUserProfileResponse, int, error) {
	selectResponse, err := surreal.Query("SELECT * FROM user:" + string(userId))

	if err != nil {
		logging.DbLogger.Error(err)
		return GetUserProfileResponse{}, http.StatusUnauthorized, errors.SurrealQueryError.New("cannot get user image data")
	}

	var selectUserResponse models.UserSelectResult
	responseString := surreal.FormatResponse(selectResponse)
	err = json.Unmarshal([]byte(responseString), &selectUserResponse)
	if err != nil {
		logging.ApiLogger.Error(err)
		return GetUserProfileResponse{}, http.StatusBadRequest, errors.SurrealDeserializaError.New("cannot deserialize user image data")
	}

	if len(selectUserResponse.Result) < 1 {
		return GetUserProfileResponse{}, http.StatusNotFound, errors.UserNotFoundError.New("user not found")
	}

	if selectUserResponse.Result[0].Image.Name == "" {
		return GetUserProfileResponse{}, http.StatusNotFound, errors.UserNoProfileError.New("user has no profile image")
	}

	imageData := selectUserResponse.Result[0].Image
	fileName := util.FileName(imageData.Name)

	fileData, err := os.ReadFile(fileName)
	if err != nil {
		logging.FsLogger.Error(err)
		return GetUserProfileResponse{}, http.StatusBadRequest, errors.CannotReadFileError.New("cannot read file from disk")
	}
	base64data := base64.StdEncoding.EncodeToString(fileData)

	return GetUserProfileResponse{
		FileName: fileName,
		FileData: base64data,
		FileType: fmt.Sprintf("image/%s", imageData.Type),
	}, http.StatusOK, nil
}
