package test

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCampaign(t *testing.T) {
	ClearAll()

	requestBody := model.CreateCampaignRequest{
		Name:             "sebuah campaign yang sangat biasa",
		ShortDescription: "sebuah deskripsi singkat biasa",
		Description:      "penjelasan yang pendek",
		GoalAmount:       10000000,
		Perks:            "keuntungan satu, keuntungan dua, dan keuntungan  ketiga",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/campaigns", strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.CampaignResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.ShortDescription, responseBody.Data.ShortDescription)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
	assert.Equal(t, requestBody.GoalAmount, responseBody.Data.GoalAmount)
	assert.Equal(t, 0, responseBody.Data.CurrentAmount)
	assert.Equal(t, requestBody.Perks, responseBody.Data.Perks)
	assert.Equal(t, 0, responseBody.Data.BackerCount)
	assert.NotNil(t, responseBody.Data.Slug)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestSearchCampaign(t *testing.T) {
	ClearAll()

	CreateCampaigns(&entity.Campaign{}, 20)

	request := httptest.NewRequest(http.MethodGet, "/api/campaigns", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.CampaignResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
}
