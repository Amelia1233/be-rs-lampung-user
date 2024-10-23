package usecase

import (
	"be-rs-lampung-user/auth/repository"
	"be-rs-lampung-user/entity"
	"be-rs-lampung-user/jwt"
	"encoding/json"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	repo *repository.AuthRepository
}

func NewAuthUsecase(repo *repository.AuthRepository) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (u *AuthUsecase) Login(ctx *gin.Context) (string, string, error) {
	var login entity.Login
	if err := ctx.ShouldBindJSON(&login); err != nil {
		return "", "", err
	}

	user, err := u.repo.GetUserByUsername(login.Username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return "", "", errors.New("password salah")
	}

	var roles []string
	if err := json.Unmarshal([]byte(user.Roles), &roles); err != nil {
		return "", "", err
	}

	token, err := jwt.GenerateToken(user.Username, roles)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := u.repo.SaveRefreshToken(user.Username, refreshToken); err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (u *AuthUsecase) GetUserData(username string) (*entity.UserResponse, error) {
	user, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	var roles []string
	if err := json.Unmarshal([]byte(user.Roles), &roles); err != nil {
		return nil, err
	}

	permissions := u.getPermissionsForRoles(roles)

	return &entity.UserResponse{
		Name:        user.Name,
		Email:       user.Email,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (u *AuthUsecase) RefreshToken(refreshToken string) (string, error) {
	user, err := u.repo.GetUserByRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("refresh token tidak valid")
	}

	var roles []string
	if err := json.Unmarshal([]byte(user.Roles), &roles); err != nil {
		return "", err
	}

	newToken, err := jwt.GenerateToken(user.Username, roles)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (u *AuthUsecase) GenerateShortLivedToken(username string) (string, error) {
	user, err := u.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	var roles []string
	if err := json.Unmarshal([]byte(user.Roles), &roles); err != nil {
		return "", err
	}

	return jwt.GenerateTokenWithExpiration(username, roles, time.Second)
}

func (u *AuthUsecase) InitStaticAdmins() error {
	admins := []entity.User{
		{
			Username: "admin",
			Password: "adminskynet100",
			Roles:    `["ADMIN_STAFF", "IT_ADMIN"]`,
			Name:     "Adi",
			Email:    "adi@gmil.com.com",
		},
		{
			Username: "dokter",
			Password: "dokter123",
			Roles:    `["DOKTER"]`,
			Name:     "Ani",
			Email:    "ani@gmil.com.com",
		},
	}

	for _, admin := range admins {
		existingUser, err := u.repo.GetUserByUsername(admin.Username)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existingUser != nil {
			existingUser.Roles = admin.Roles
			existingUser.Name = admin.Name
			existingUser.Email = admin.Email

			if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(admin.Password)); err != nil {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
				if err != nil {
					return err
				}
				existingUser.Password = string(hashedPassword)
			}

			admin = *existingUser
		} else {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			admin.Password = string(hashedPassword)
		}

		if err := u.repo.CreateOrUpdateUser(&admin); err != nil {
			return err
		}
	}

	return nil
}

func (u *AuthUsecase) getPermissionsForRoles(roles []string) []string {
	permissionSet := make(map[string]bool)
	for _, role := range roles {
		switch role {
		case "ADMIN_STAFF", "IT_ADMIN":
			permissionSet["VIEW_USR_MGT"] = true
			permissionSet["CREATE_USR_MGT"] = true
			permissionSet["VIEW_RPT_BILL"] = true
			permissionSet["EDIT_USR_MGT"] = true
			permissionSet["DELETE_USR_MGT"] = true
			permissionSet["GEN_RPT_BILL"] = true
		// Tambahkan case untuk role lain jika diperlukan
		
		case "DOKTER":
			permissionSet["MEMERIKSA"] = true
	}

	}

	var permissions []string
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}
	return permissions
}
