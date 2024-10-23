package handler

import (
	"be-rs-lampung-user/auth/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

func NewAuthHandler(usecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func (h *AuthHandler) Login(c *gin.Context) {
	token, refreshToken, err := h.usecase.Login(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, int((15 * time.Minute).Seconds()), "/", "", false, true)

	c.SetCookie("refresh_token", refreshToken, int((30 * 24 * time.Hour).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil"})
}

func (h *AuthHandler) GetUserData(c *gin.Context) {
	username := c.GetString("username")
	userData, err := h.usecase.GetUserData(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userData)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token tidak ditemukan"})
		return
	}

	newToken, err := h.usecase.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set token baru sebagai HTTP-only cookie
	c.SetCookie("token", newToken, int(time.Hour.Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token berhasil diperbarui"})
}
