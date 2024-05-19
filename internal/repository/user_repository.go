package repository

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) CountByName(db *gorm.DB, user *entity.User) (int64, error) {
	var total int64
	err := db.Model(user).Where("name = ?", user.Name).Count(&total).Error
	return total, err
}

func (r *UserRepository) CountByEmail(db *gorm.DB, user *entity.User) (int64, error) {
	var total int64
	err := db.Model(user).Where("email = ?", user.Email).Count(&total).Error
	return total, err
}

func (r *UserRepository) FindByName(db *gorm.DB, user *entity.User, name string) error {
	return db.Where("name = ?", name).Take(user).Error
}

func (r *UserRepository) Search(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, int64, error) {
	var users []entity.User
	if err := db.Scopes(r.FilterUser(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.User{}).Scopes(r.FilterUser(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FilterUser(request *model.SearchUserRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if email := request.Email; email != "" {
			email = "%" + email + "%"
			tx = tx.Where("email LIKE ?", email)
		}

		if role := request.Role; role != "" {
			tx = tx.Where("role = ?", role)
		}

		return tx
	}
}
