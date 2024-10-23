package auth

import (
	"be-rs-lampung-user/entity"

	"github.com/gin-gonic/gin"
)

type RepoAuth interface {
	Login(string) (*entity.Login, error)
}

type UsecaseAuth interface {
	Login(ctx *gin.Context) error
}
