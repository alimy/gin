package ginL

import (
	"github.com/alimy/gin"
	"github.com/unisx/logus"
)

// GinLogger instances a Logger middleware that use logus write the logs to zap
func GinLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func (context *gin.Context) {
		// TODO
		logus.Info("gin handle", logus.String("clientIP", context.ClientIP()))
	})
}

// GinErrorLogger returns a gin.HandlerFunc for any error type.
func GinErrorLogger() gin.HandlerFunc {
	return GinErrorLoggerT(gin.ErrorTypeAny)
}

// GinErrorLoggerT returns a gin.HandlerFunc for a given error type.
func GinErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return gin.HandlerFunc(func (context *gin.Context) {
		// TODO
		logus.Info("gin handle", logus.String("clientIP", context.ClientIP()))
	})
}