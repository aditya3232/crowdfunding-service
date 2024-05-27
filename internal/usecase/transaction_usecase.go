package usecase

import (
	"context"
	"crowdfunding-service/internal/entity"
	payment_gateway "crowdfunding-service/internal/gateway/payment-gateway"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/model/converter"
	"crowdfunding-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionUseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	TransactionRepository  *repository.TransactionRepository
	CampaignRepository     *repository.CampaignRepository
	UserRepository         *repository.UserRepository
	MidtransPaymentGateway *payment_gateway.MidtransPaymentGateway
}

func NewTransactionUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	transactionRepository *repository.TransactionRepository, campaignRepository *repository.CampaignRepository,
	userRepository *repository.UserRepository, midtransPaymentGateway *payment_gateway.MidtransPaymentGateway) *TransactionUseCase {
	return &TransactionUseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		TransactionRepository:  transactionRepository,
		CampaignRepository:     campaignRepository,
		UserRepository:         userRepository,
		MidtransPaymentGateway: midtransPaymentGateway,
	}
}

/*
- kita harus menampilkan transactions pada campaigns, sesuai dengan user_id yg berada di campaigns
- karena user di tb campaign adalah user pemilik campaign, sedangkan user di tb transaction adalah user yg memberikan funding,
- maka kita harus membandingkan user_id yg melakukan request dengan user_id yg ada di tb campaign
*/
func (u *TransactionUseCase) GetTransactionsByCampaignID(ctx context.Context, request *model.GetTransactionByCampaignIDRequest) ([]model.GetTransactionByCampaignIDResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	// check if campaign exists
	// pakai GetByID karena kita hanya butuh 1 data, dan sudah ada preloading user
	campaign := new(entity.Campaign)
	if err := u.CampaignRepository.GetByID(tx, campaign, request.CampaignID); err != nil {
		u.Log.WithError(err).Error("error getting campaign by ID")
		return nil, 0, fiber.ErrNotFound
	}

	// check if user is the owner of the campaign
	if campaign.UserID != request.UserID {
		u.Log.Error("user is not the owner of the campaign")
		return nil, 0, fiber.ErrForbidden
	}

	transactions, total, err := u.TransactionRepository.GetTransactionByCampaignID(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("error getting transactions by campaign ID")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.GetTransactionByCampaignIDResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *converter.GetTransactionByCampaignIDToResponse(&transaction)
	}

	return responses, total, nil
}

func (u *TransactionUseCase) GetTransactionsByUserID(ctx context.Context, request *model.GetTransactionByUserIDRequest) ([]model.GetTransactionByUserIDResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	transactions, total, err := u.TransactionRepository.GetTransactionByUserID(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("error getting transactions by user ID")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.GetTransactionByUserIDResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *converter.GetTransactionByUserIDToResponse(&transaction)
	}

	return responses, total, nil
}

/*
- disini ada penerapan transactional
- Membuat SavePoint sebelum membuat transaksi dengan tx.SavePoint("sp_before_create_transaction")
- Membuat transaksi dan jika terjadi kesalahan, akan melakukan tx.RollbackTo("sp_before_create_transaction")
- Jika semua langkah berjalan dengan baik, maka melakukan tx.Commit() untuk menyimpan semua perubahan dalam transaksi
- fungsinya adalah jika terjadi kesalahan di tengah-tengah proses, kita bisa melakukan rollback ke savepoint yang sudah dibuat
- disini perlu savpoint karena create harus berhasil, karena id transaction yang di create akan digunakan untuk membuat payment URL
- baru kemudian update transaction dengan payment URL
*/
func (u *TransactionUseCase) CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*model.TransactionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// check if campaign exists
	campaign := new(entity.Campaign)
	if err := u.CampaignRepository.GetByID(tx, campaign, request.CampaignID); err != nil {
		u.Log.WithError(err).Error("error getting campaign by ID")
		return nil, fiber.ErrNotFound
	}

	// check if user exists
	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.UserID); err != nil {
		u.Log.WithError(err).Error("error getting user by ID")
		return nil, fiber.ErrNotFound
	}

	transaction := &entity.Transaction{
		ID:         uuid.New().String(),
		CampaignID: request.CampaignID,
		UserID:     request.UserID,
		Amount:     request.Amount,
		Status:     "pending",
	}

	// Create a savepoint before creating the transaction
	if err := tx.SavePoint("sp_before_create_transaction").Error; err != nil {
		u.Log.WithError(err).Error("error creating savepoint")
		return nil, fiber.ErrInternalServerError
	}

	// disini pakai create repository custom, karena kita ingin mengembalikan data id transaction yang sudah di create
	// create the transaction
	NewTransaction, err := u.TransactionRepository.CreateTransaction(tx, transaction)
	if err != nil {
		u.Log.WithError(err).Error("error creating transaction")
		tx.RollbackTo("sp_before_create_transaction")
		return nil, fiber.ErrInternalServerError
	}

	// log data NewTransaction
	u.Log.WithField("NewTransaction", NewTransaction).Info("New Transaction")

	// generate payment URL
	payment := &model.PaymentRequest{
		Email:         user.Email,
		FullName:      user.Name,
		TransactionID: NewTransaction.ID,
		Amount:        NewTransaction.Amount,
	}

	// minta payment URL ke midtrans
	paymentURL, err := u.MidtransPaymentGateway.GetPaymentURL(payment)
	if err != nil {
		u.Log.WithError(err).Error("error getting payment URL")
		return nil, fiber.ErrInternalServerError
	}

	// update data payment URL di transaction
	transaction.PaymentURL = paymentURL
	if err := u.TransactionRepository.Update(tx, transaction); err != nil {
		u.Log.WithError(err).Error("error updating transaction")
		tx.RollbackTo("sp_before_create_transaction")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.TransactionToResponse(transaction), nil

}

// usecase ini yg akan digunakan midtrans mengirim notification status pembayaran ke service kita
func (u *TransactionUseCase) CreateTransactionNotification(ctx context.Context, request *model.CreateTransactionNotificationRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	transaction := new(entity.Transaction)
	if err := u.TransactionRepository.FindById(tx, transaction, request.TransactionID); err != nil {
		u.Log.WithError(err).Error("error getting transaction by ID")
		return fiber.ErrNotFound
	}

	// update status transaction
	if request.PaymentType == "credit_card" && request.TransactionStatus == "capture" && request.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if request.TransactionStatus == "cancel" || request.TransactionStatus == "deny" || request.TransactionStatus == "expire" {
		transaction.Status = "cancelled"
	}

	if err := u.TransactionRepository.Update(tx, transaction); err != nil {
		u.Log.WithError(err).Error("error updating transaction")
		return fiber.ErrInternalServerError
	}

	// check campaign id di transaction
	campaign := new(entity.Campaign)
	if err := u.CampaignRepository.GetByID(tx, campaign, transaction.CampaignID); err != nil {
		u.Log.WithError(err).Error("error getting campaign by ID")
		return fiber.ErrNotFound
	}

	// update status campaign
	if transaction.Status == "paid" {
		campaign.BackerCount += 1                    // tambah backer count, biar fleksibel pakai += 1, bisa diganti dengan jumlah backer
		campaign.CurrentAmount += transaction.Amount // tambah current amount dengan amount dari transaction

		if err := u.CampaignRepository.Update(tx, campaign); err != nil {
			u.Log.WithError(err).Error("error updating campaign")
			return fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return fiber.ErrInternalServerError
	}

	return nil

}
