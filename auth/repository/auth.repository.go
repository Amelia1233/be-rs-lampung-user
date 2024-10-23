package repository

import (
	"be-rs-lampung-user/entity"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateOrUpdateUser(user *entity.User) error {
	return r.DB.Save(user).Error
}

func (r *AuthRepository) SaveRefreshToken(username, refreshToken string) error {
	return r.DB.Model(&entity.User{}).Where("username = ?", username).Update("refresh_token", refreshToken).Error
}

func (r *AuthRepository) GetUserByRefreshToken(refreshToken string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("refresh_token = ?", refreshToken).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(user *entity.User) error {
	return r.DB.Save(user).Error
}

func (r *AuthRepository) CreateUser(user *entity.User) error {
	return r.DB.Create(user).Error
}
