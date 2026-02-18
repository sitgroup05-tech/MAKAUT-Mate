package user

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sitgroup05-tech/MakautMate/backend/internal/model"
	"gorm.io/gorm"
)

type ValidationClaim struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Mail   string `json:"mail"`
	jwt.RegisteredClaims
}

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionManager struct {
	secretKey string
	sessionService
}

func NewSessionManager(key string, db *gorm.DB) SessionManager {
	return SessionManager{
		secretKey: key,
		sessionService: sessionService{
			db: db,
		},
	}
}

func (s SessionManager) Verify(role string) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(
				http.StatusUnauthorized, gin.H{
					"error": "Auth Not Found",
				},
			)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, _ := jwt.ParseWithClaims(tokenString, &ValidationClaim{}, func(t *jwt.Token) (interface{}, error) {
			return s.secretKey, nil
		})

		if claims, ok := token.Claims.(*ValidationClaim); ok && token.Valid && claims.Role == role {
			ctx.Set("user-id", claims.UserID)
			ctx.Set("role", claims.Role)
			ctx.Set("mail", claims.Mail)
		} else {
			ctx.JSON(
				http.StatusUnauthorized, gin.H{
					"error": "Invalid Auth or Not Have Permission",
				},
			)
			return
		}
		ctx.Next()
	})
}

// Login String
// abc
func (s SessionManager) Login(ctx *gin.Context) {
	var login_credentials LoginDetails
	if err := ctx.ShouldBindJSON(&login_credentials); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Invalid Credentials",
			},
		)
		return
	}

	var user model.User

	if err := s.GetUser(ctx.Request.Context(), login_credentials.Email, &user); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Invalid User-Name or Password", // invalid user-name
			},
		)
		return
	}

	if err := s.ValidatePassword(user.Password, login_credentials.Password); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Invalid User-Name or Password", // invalid password
			},
		)
		return

	}

	my_claim := ValidationClaim{
		user.Id,
		user.Role,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "Makaut-Mate-Server",
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, &my_claim)

	tokenString, err := tok.SignedString(s.secretKey)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "Server Error",
			},
		)
		log.Printf("Error While Signing jwt :: %s", err)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"token": tokenString,
		},
	)

}
func (s SessionManager) Logout(ctx *gin.Context) {

}
func (s SessionManager) Update(ctx *gin.Context) {

}
