package handlers

import (
	"be-rs-lampung-user/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHandlers(us auth.UsecaseAuth, r *gin.RouterGroup) {
	eng := &Usecase{
		usAuth: us,
	}

	v2 := r.Group("auth")
	v2.POST("/login", eng.Login)
}

type Usecase struct {
	usAuth auth.UsecaseAuth
}

func (us Usecase) Login(ctx *gin.Context) {
	err := us.usAuth.Login(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"Message": "Succes Login"})
}
