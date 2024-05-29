package test

import (
	"crowdfunding-service/internal/entity"
	"mime/multipart"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// dihelper tidak dapat ditest, dia hanya berisi fungsi-fungsi yang membantu dalam testing (tidak berawalan Testxxx)

func ClearAll() {
	ClearUsers()
}

func IsImage(file *multipart.FileHeader) bool {
	switch file.Header.Get("Content-Type") {
	case "image/jpeg", "image/jpg", "image/png":
		return true
	default:
		return false
	}
}

func ClearUsers() {
	err := db.Where("id is not null").Not("email = ?", "iashiddiqi13@gmail.com").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func CreateUsers(user *entity.User, total int) {
	password := "password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatalf("Failed to hash password : %+v", err)
	}

	for i := 0; i < total; i++ {
		user := &entity.User{
			ID:           uuid.NewString(),
			Name:         "User " + strconv.Itoa(i),
			Occupation:   "Programmer",
			Email:        "user" + strconv.Itoa(i) + "@gmail.com",
			PasswordHash: string(hashedPassword),
			Role:         "user",
		}
		err := db.Create(user).Error
		if err != nil {
			log.Fatalf("Failed to create user : %+v", err)
		}
	}
}

func GetFirstUser() *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	if err != nil {
		log.Fatalf("Failed to get first user : %+v", err)
	}
	return user
}
