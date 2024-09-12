package middleware

import (
	"net/http"

	jwtauth "github.com/Lafetz/loyalty_marketplace/internal/web/jwt"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			c.Abort()
			return
		}

		user, err := jwtauth.PareseJwt(jwtToken)
		if err != nil {
			if err == jwtauth.ErrInvalidToken {
				c.JSON(http.StatusUnauthorized, gin.H{
					"Error": "Unauthorized",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"Error": "Internal server error"})
			c.Abort()
			return

		}
		c.Set("user", user.GetUserToken())

		c.Next()
	}
}
