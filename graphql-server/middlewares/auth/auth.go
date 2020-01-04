package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"root/models"
	"strings"

	"github.com/hako/branca"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var contxt context.Context
var users []*models.User
var envLoadErr = godotenv.Load()

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if envLoadErr != nil {
				log.Fatal(envLoadErr)
			}


			bearerToken := r.Header.Get("Authorization")

			if bearerToken == "" {
				next.ServeHTTP(w, r)
			} else {

				bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)

				secretKey := os.Getenv("SECRET_KEY")
				b := branca.NewBranca(secretKey)

				message, err := b.DecodeToString(bearerToken)
				if err != nil {
					log.Println(fmt.Errorf("An undecodable or invalid token was entered."))
					return
				}

				userInfo := strings.Split(message, "/")
				userEmail := userInfo[0]
				userPassword := userInfo[1]

				db, err := gorm.Open("postgres", os.Getenv("GORM_SETUP"))
				defer db.Close()

				db.Where("email = ? AND password = ?", userEmail, userPassword).First(&users)
				user := users[0]

				ctx := context.WithValue(r.Context(), "User", user)
				next.ServeHTTP(w, r.WithContext(ctx))
			}

		})
	}
}

func UserContextExtracter(ctx context.Context) *models.User {
	raw, _ := ctx.Value("User").(*models.User)
	return raw
}
