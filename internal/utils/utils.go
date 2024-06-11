package utils

import (
	"time"
    // "strings"
    // "github.com/dgrijalva/jwt-go"
	// "github.com/gin-gonic/gin"
)


// CalculateAge calculates the age of a person based on their birthdate.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age
}


// func checkTokenAndGetUserId(tokenString string) (string, error) {
//     // Extract the token from the Authorization header
//     tokenParts := strings.Split(tokenString, " ")
//     if len(tokenParts) != 2 {
//         return "", errors.New("Invalid token format")
//     }
//     tokenString := tokenParts[1]

//     // Parse the token
//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         // Verify the token signature using the key provided by Keycloak
//         return []byte("your-keycloak-public-key"), nil
//     })
//     if err != nil {
//         return "", err
//     }

//     // Extract the user ID from the token claims
//     if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//         if userID, ok := claims["sub"].(string); ok {
//             return userID, nil
//         }
//     }

//     return "", errors.New("Failed to extract user ID from token")
// }


// func checkTokenAndGetUserId(c *gin.Context) (string, error) {
// 	// get the token from request header
//     token := c.GetHeader("Authorization")
//     if token == "" {
//         c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Authorization token is required"))
//         return
//     }
	
//     // Extract the token from the Authorization header
//     tokenParts := strings.Split(tokenString, " ")
//     if len(tokenParts) != 2 {
//         return "", errors.New("Invalid token format")
//     }
//     tokenString := tokenParts[1]

//     // Parse the token
//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         // Verify the token signature using the key provided by Keycloak
//         return []byte("your-keycloak-public-key"), nil
//     })
//     if err != nil {
//         return "", err
//     }

//     // Extract the user ID from the token claims
//     if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//         if userID, ok := claims["sub"].(string); ok {
//             return userID, nil
//         }
//     }

//     return "", errors.New("Failed to extract user ID from token")
// } 