package usecase

import (
	"context"
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/model/converter"
	"crowdfunding-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		u.Log.WithError(err).Error("error hashing password")
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		ID:           uuid.New().String(),
		Name:         request.Name,
		Occupation:   request.Occupation,
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Role:         request.Role,
	}

	totalName, err := u.UserRepository.CountByName(tx, user)
	if err != nil {
		u.Log.WithError(err).Error("error checking name availability")
		return nil, fiber.ErrInternalServerError
	}

	if totalName > 0 {
		u.Log.Warnf("Name already taken : %+v", err)
		return nil, fiber.ErrConflict
	}

	totalEmail, err := u.UserRepository.CountByEmail(tx, user)
	if err != nil {
		u.Log.WithError(err).Error("error checking email availability")
		return nil, fiber.ErrInternalServerError
	}

	if totalEmail > 0 {
		u.Log.Warnf("Email already taken : %+v", err)
		return nil, fiber.ErrConflict
	}

	if err := u.UserRepository.Create(tx, user); err != nil {
		u.Log.WithError(err).Error("error creating user")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		u.Log.WithError(err).Error("error finding user")
		return nil, fiber.ErrNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		u.Log.WithError(err).Error("error hashing password")
		return nil, fiber.ErrInternalServerError
	}

	user.Name = request.Name
	user.Occupation = request.Occupation
	user.Email = request.Email
	user.PasswordHash = string(hashedPassword)
	user.Role = request.Role

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.WithError(err).Error("error updating user")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) Get(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		u.Log.WithError(err).Error("error finding user")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (u *UserUseCase) Delete(ctx context.Context, request *model.DeleteUserRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		u.Log.WithError(err).Error("error finding user")
		return fiber.ErrNotFound
	}

	if err := u.UserRepository.Delete(tx, user); err != nil {
		u.Log.WithError(err).Error("error deleting user")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *UserUseCase) Search(ctx context.Context, request *model.SearchUserRequest) ([]model.UserResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	users, total, err := u.UserRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("error searching user")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *converter.UserToResponse(&user)
	}

	return responses, total, nil
}
