package middlewares

import (
	"fmt"
	"lb-api/models"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func JWT() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("goFuckYourSelf!"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userId string, password string, c *gin.Context) (interface{}, bool) {
			// if (userId == "admin" && password == "admin") || (userId == "test" && password == "test") {
			// 	return &models.User{
			// 		UserName:  userId,
			// 		LastName:  "Bo-Yi",
			// 		FirstName: "Wu",
			// 	}, true
			// }
			fmt.Println("Authenticator", userId, password)
			var user = &models.User{Phone: userId, Passwd: password}
			if err := user.ValidUser(); nil == err {

				return user, true
			}

			return nil, false
		},
		Authorizator: func(user interface{}, c *gin.Context) bool {
			fmt.Println("Authorizator\n\n", user)
			if v, ok := user.(string); ok && v != "" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:         time.Now,
		SigningAlgorithm: "HS256",
	}

}

func GetUserIDFromToken(c *gin.Context) string {
	fmt.Println("JWT---", jwt.ExtractClaims(c))
	return jwt.ExtractClaims(c)["id"].(string)
}
func GetPalyloadFromToken(c *gin.Context) map[string]interface{} {
	fmt.Println("JWT---", jwt.ExtractClaims(c))

	return jwt.ExtractClaims(c)
}
func GetJWTToken(userID string, data map[string]interface{}) (string, time.Time, error) {
	token, expt, err := JWT().TokenGenerator(userID, data)
	return token, expt, err
}
