package test

import (
	"bytes"
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultUser(t *testing.T) {
	ClearAll()
	CreateDefaultUser()
}

func TestRegister(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		Name:       "User 1",
		Occupation: "Programmer",
		Email:      "user1@gmail.com",
		Password:   "password",
		Role:       "user",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Occupation, responseBody.Data.Occupation)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()

	requestBody := model.RegisterUserRequest{
		Name:       "",
		Occupation: "",
		Email:      "",
		Password:   "",
		Role:       "",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()

	TestRegister(t)

	requestBody := model.RegisterUserRequest{
		Name:       "User 1",
		Occupation: "Programmer",
		Email:      "user1@gmail.com",
		Password:   "password",
		Role:       "user",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateAvatar(t *testing.T) {
	ClearAll()

	file, err := os.Open("sakurawb.jpg")
	assert.Nil(t, err)
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload_avatar", "sakurawb.jpg")
	assert.Nil(t, err)

	_, err = io.Copy(part, file)
	assert.Nil(t, err)

	err = writer.Close()
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/avatar/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.Occupation)
	assert.NotNil(t, responseBody.Data.Email)
	assert.NotNil(t, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
	assert.NotNil(t, responseBody.Data.Avatar)
}

func TestUpdateAvatarNotValidFile(t *testing.T) {
	ClearAll()

	file, err := os.Open("not_valid_img.pdf")
	assert.Nil(t, err)
	defer file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("upload_avatar", "not_valid_img.pdf")
	assert.Nil(t, err)

	_, err = io.Copy(part, file)
	assert.Nil(t, err)

	err = writer.Close()
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/avatar/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestUpdateUser(t *testing.T) {
	ClearAll()

	CreateUsers(&entity.User{}, 1)

	user := new(entity.User)
	err := db.Where("email = ?", "user0@gmail.com").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Name:       "User 1",
		Occupation: "Programmer",
		Email:      "user0@gmail.com",
		Password:   "password",
		Role:       "tes update",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/"+user.ID, strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Occupation, responseBody.Data.Occupation)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

}

func TestUpdateUserError(t *testing.T) {
	ClearAll()

	CreateUsers(&entity.User{}, 1)

	user := new(entity.User)
	err := db.Where("email = ?", "user0@gmail.com").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Name:       "User 1",
		Occupation: "Programmer",
		Email:      "user0@gmail.com",
		Role:       "tes update",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, "/api/users/"+user.ID, strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.Occupation)
	assert.NotNil(t, responseBody.Data.Email)
	assert.NotNil(t, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserError(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", "invalid-token")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestGetUser(t *testing.T) {
	ClearAll()

	TestRegister(t)

	user := new(entity.User)
	err := db.Where("email = ?", "user1@gmail.com").First(user).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodGet, "/api/users/"+user.ID, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.Occupation)
	assert.Equal(t, user.Email, responseBody.Data.Email)
	assert.NotNil(t, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestGetUserError(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, "/api/users/invalid-uuid", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Errors)
}

func TestSearchUser(t *testing.T) {
	ClearAll()

	CreateUsers(&entity.User{}, 20)

	request := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(22), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(3), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestSearchUserWithPagination(t *testing.T) {
	ClearAll()

	CreateUsers(&entity.User{}, 20)

	request := httptest.NewRequest(http.MethodGet, "/api/users?page=2&size=5", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
	assert.Equal(t, int64(22), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(5), responseBody.Paging.TotalPage)
	assert.Equal(t, 2, responseBody.Paging.Page)
	assert.Equal(t, 5, responseBody.Paging.Size)
}

func TestSearchUserWithFilter(t *testing.T) {
	ClearAll()

	CreateUsers(&entity.User{}, 20)

	name := "Muhammad Aditya"
	encodedName := url.QueryEscape(name)

	request := httptest.NewRequest(http.MethodGet, "/api/users?name="+encodedName, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, int64(1), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(1), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}
