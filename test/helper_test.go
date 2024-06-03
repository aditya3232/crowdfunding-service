package test

import (
	"crowdfunding-service/internal/entity"
	"strconv"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

/*
- dihelper tidak dapat ditest, dia hanya berisi fungsi-fungsi yang membantu dalam testing (tidak berawalan Testxxx)
- helper langsung menjalankan fungsi-fungsi yang dibutuhkan dalam testing, bukan http request
*/

func ClearAll() {
	ClearUsers()
	ClearCampaigns()
}

func CreateDefaultUser() {
	password := "password"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %+v", err)
	}

	user1 := &entity.User{
		ID:           uuid.NewString(),
		Name:         "Muhammad Aditya",
		Occupation:   "Programmer",
		Email:        "m.aditya3232@gmail.com",
		PasswordHash: string(hashedPassword),
		Role:         "admin",
	}

	user2 := &entity.User{
		ID:           uuid.NewString(),
		Name:         "Ichsan Ashiddiqi",
		Occupation:   "Programmer",
		Email:        "iashiddiqi13@gmail.com",
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	existingUser1 := &entity.User{}
	err = db.Where("email = ?", user1.Email).First(existingUser1).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(user1).Error
			if err != nil {
				log.Fatalf("Failed to create user 1: %+v", err)
			}
		} else {
			log.Fatalf("Failed to check existing user 1: %+v", err)
		}
	}

	existingUser2 := &entity.User{}
	err = db.Where("email = ?", user2.Email).First(existingUser2).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = db.Create(user2).Error
			if err != nil {
				log.Fatalf("Failed to create user 2: %+v", err)
			}
		} else {
			log.Fatalf("Failed to check existing user 2: %+v", err)
		}
	}
}

func ClearUsers() {
	err := db.Where("id is not null").Not("email = ? OR email = ?", "m.aditya3232@gmail.com", "iashiddiqi13@gmail.com").Delete(&entity.User{}).Error
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

func GetDefaultUser() *entity.User {
	user := new(entity.User)
	err := db.Where("email = ?", "m.aditya3232@gmail.com").First(user).Error
	if err != nil {
		log.Fatalf("Failed to get default user : %+v", err)
	}

	return user
}

func ClearCampaigns() {
	err := db.Where("id is not null").Delete(&entity.Campaign{}).Error
	if err != nil {
		log.Fatalf("Failed clear campaign data : %+v", err)
	}
}

func CreateCampaigns(campaign *entity.Campaign, total int) {
	for i := 0; i < total; i++ {
		campaign := &entity.Campaign{
			ID:               uuid.NewString(),
			Name:             "sebuah campaign yang sangat biasa " + strconv.Itoa(i),
			ShortDescription: "sebuah deskripsi singkat biasa",
			Description:      "penjelasan yang pendek",
			GoalAmount:       10000000,
			Perks:            "keuntungan satu, keuntungan dua, dan keuntungan  ketiga",
			UserID:           GetDefaultUser().ID,
		}
		err := db.Create(campaign).Error
		if err != nil {
			log.Fatalf("Failed to create campaign : %+v", err)
		}
	}
}
