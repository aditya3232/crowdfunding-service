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

func TestCreateTransaction(t *testing.T) {
	ClearAll()

	TestCreateCampaign(t)

	campaign := new(entity.Campaign)
	err := db.Where("name = ?", "sebuah campaign yang sangat biasa").First(campaign).Error
	assert.Nil(t, err)

	requestBody := model.CreateTransactionRequest{
		CampaignID: campaign.ID,
		Amount:     1000000,
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/transactions", strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.TransactionResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, campaign.ID, responseBody.Data.CampaignID)
	assert.NotNil(t, responseBody.Data.UserID)
	assert.Equal(t, requestBody.Amount, responseBody.Data.Amount)
	assert.Equal(t, "pending", responseBody.Data.Status) // status pending karena belum ada notifikasi dari midtrans
	assert.NotNil(t, responseBody.Data.PaymentURL)
	assert.NotNil(t, responseBody.Data.CreatedAt)
}

func TestCreateNotifcationFromMidtrans(t *testing.T) {
	ClearAll()

	TestCreateTransaction(t)

	transaction := new(entity.Transaction)
	err := db.Where("status = ?", "pending").First(transaction).Error
	assert.Nil(t, err)

	requestBody := model.CreateTransactionNotificationRequest{
		TransactionStatus: "capture",
		TransactionID:     transaction.ID,
		PaymentType:       "credit_card",
		FraudStatus:       "accept",
	}

	jsonByte, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/transactions/notification", strings.NewReader(string(jsonByte)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[interface{}])
	responseBody.Data = make(map[string]interface{})
	responseBody.Data = "transaction notification created"
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	// check campaign amount and backer count after transaction notification
	campaign := new(entity.Campaign)
	err = db.Where("id = ?", transaction.CampaignID).First(campaign).Error
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "transaction notification created", responseBody.Data)
	assert.Equal(t, 1000000, campaign.CurrentAmount)
	assert.Equal(t, 1, campaign.BackerCount)
}

// daftar transaksi yg dilakukan oleh user yg login
func TestGetTransactionsByUserID(t *testing.T) {
	ClearAll()

	TestCreateTransaction(t)

	request := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.GetTransactionByUserIDResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, int64(1), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(1), responseBody.Paging.TotalPage)
	assert.Equal(t, 10, responseBody.Paging.Size)
	assert.Equal(t, "pending", responseBody.Data[0].Status)
	assert.Equal(t, 1000000, responseBody.Data[0].Amount)
}

// daftar transaksi pada campaign yg dipilih berdasarkan campaign_id & user_id pemilik campaign
func TestGetTransactionsByCampaignID(t *testing.T) {
	ClearAll()

	TestCreateTransaction(t)

	campaign := new(entity.Campaign)
	err := db.Where("name = ?", "sebuah campaign yang sangat biasa").First(campaign).Error
	assert.Nil(t, err)

	user := new(entity.User)
	err = db.Where("id = ?", campaign.UserID).First(user).Error
	assert.Nil(t, err)

	encodedUserID := url.QueryEscape(user.ID) // user id pemilik campaign

	request := httptest.NewRequest(http.MethodGet, "/api/transactions/campaigns/"+campaign.ID+"?user_id="+encodedUserID, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("OAuth2-token", viperConfig.GetString("test.oauth2.google.accessToken"))

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[[]model.GetTransactionByCampaignIDResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(responseBody.Data))
	assert.Equal(t, int64(1), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(1), responseBody.Paging.TotalPage)
	assert.Equal(t, 10, responseBody.Paging.Size)
	assert.NotNil(t, responseBody.Data[0].ID)
	assert.NotNil(t, responseBody.Data[0].UserName)
	assert.NotNil(t, responseBody.Data[0].Amount)
	assert.NotNil(t, responseBody.Data[0].CreatedAt)
}
