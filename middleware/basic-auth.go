package middleware

import "github.com/gin-gonic/gin"

func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"iman":  "123456",
		"admin": "admin123",
		"user":  "user123",
		"guest": "guest123",
		"test":  "test123",
		"test2": "test2123",
		"test3": "test3123",
		"test4": "test4123",
		"test5": "test5123",
		"test6": "test6123",
		"test7": "test7123",
	})
}
