package test

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.CampaignResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestSearchCampaignWithPagination(t *testing.T) {
	ClearAll()

	CreateCampaigns(&entity.Campaign{}, 20)

	request := httptest.NewRequest(http.MethodGet, "/api/campaigns?page=2&size=5", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.CampaignResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(4), responseBody.Paging.TotalPage)
	assert.Equal(t, 2, responseBody.Paging.Page)
	assert.Equal(t, 5, responseBody.Paging.Size)

}

func TestSearchCampaignWithFilter(t *testing.T) {
	ClearAll()

	CreateCampaigns(&entity.Campaign{}, 20)

	campaignName := "sebuah campaign yang sangat biasa 0"
	encodedCampaignName := url.QueryEscape(campaignName)

	request := httptest.NewRequest(http.MethodGet, "/api/campaigns?campaign_name="+encodedCampaignName, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.CampaignResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, int64(1), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(1), responseBody.Paging.TotalPage)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestGetCampaign(t *testing.T) {
	ClearAll()

	TestCreateCampaign(t)

	campaign := new(entity.Campaign)
	err := db.Where("name = ?", "sebuah campaign yang sangat biasa").First(campaign).Error
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodGet, "/api/campaigns/"+campaign.ID, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.CampaignResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, campaign.ID, responseBody.Data.ID)
	assert.Equal(t, campaign.UserID, responseBody.Data.UserID)
	assert.Equal(t, campaign.Name, responseBody.Data.Name)
	assert.Equal(t, campaign.ShortDescription, responseBody.Data.ShortDescription)
	assert.Equal(t, campaign.Description, responseBody.Data.Description)
	assert.Equal(t, campaign.GoalAmount, responseBody.Data.GoalAmount)
	assert.Equal(t, campaign.CurrentAmount, responseBody.Data.CurrentAmount)
	assert.Equal(t, campaign.Perks, responseBody.Data.Perks)
	assert.Equal(t, campaign.BackerCount, responseBody.Data.BackerCount)
	assert.Equal(t, campaign.Slug, responseBody.Data.Slug)
	assert.Equal(t, campaign.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, campaign.UpdatedAt, responseBody.Data.UpdatedAt)
}
