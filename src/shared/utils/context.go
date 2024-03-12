package utils

import (
	"context"
	"github.com/gin-gonic/gin"
)

func OverrideGinRequest(ctx *gin.Context, key string, value interface{}) {
	c := context.WithValue(ctx.Request.Context(), key, value)
	ctx.Request = ctx.Request.WithContext(c)
}
