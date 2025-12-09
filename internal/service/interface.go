package service

import "github.com/gin-gonic/gin"

type RegisterService interface {
	Register(e *gin.Engine)
}
