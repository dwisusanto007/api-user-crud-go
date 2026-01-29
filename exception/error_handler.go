package exception

import (
	"api-user-crud-go/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler adalah middleware untuk menangani error secara global.
// Middleware ini akan menangkap panic dan error yang tidak tertangani.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lanjutkan ke handler berikutnya
		c.Next()

		// Cek apakah ada error yang terjadi
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Response dengan error yang konsisten
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: err.Error(),
			})
		}
	}
}

// Recovery adalah middleware untuk menangani panic.
// Mencegah aplikasi crash ketika terjadi panic.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: err,
			})
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "Internal Server Error",
				Message: "An unexpected error occurred",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// LoggerMiddleware adalah middleware untuk logging request.
// Middleware ini mencatat setiap request yang masuk.
func LoggerMiddleware() gin.HandlerFunc {
	return gin.Logger()
}
