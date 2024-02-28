package middlewares

import "github.com/gin-gonic/gin"

func NewCors(ctx *gin.Engine) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,PUT,DELETE")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Access-Control-Allow-Origin , Origin,Accept, X-Requested-With, Content-Type, Authorization, Access-Control-Request-Method, Access-Control-Request-Headers, X-CSRF-Token")
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Origin")

		// Handle preflight request
		if ctx.Request.Method == "OPTIONS" {
			ctx.Writer.WriteHeader(204)
			return
		}
		ctx.Next()
	}
}
