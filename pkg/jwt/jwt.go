package jwt

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mfachryna/paimon-bank/internal/common/response"
)

type Claim struct {
	jwt.StandardClaims
	UserId string `json:"user_id"`
}

type JWTToken struct {
	Token    string
	Claim    Claim
	ExpireAt time.Time
	Scheme   string
}

func SignedToken(claim Claim) (string, error) {
	exp := time.Now().Add(8 * time.Hour)
	expAt := exp.Unix()
	iat := time.Now().Unix()

	claim.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expAt,
		IssuedAt:  iat,
	}
	secretKey := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			(&response.Response{
				HttpStatus: http.StatusUnauthorized,
				Message:    "token not found",
			}).GenerateResponse(w)
			return
		}

		tokenString := string(authHeader)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Printf("unexpected signing method: %v \n", t.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			validationErr, ok := err.(*jwt.ValidationError)
			if ok {
				if validationErr.Errors == jwt.ValidationErrorExpired {
					(&response.Response{
						HttpStatus: http.StatusUnauthorized,
						Message:    "given security scheme is valid, but the lifetime has been expired or revoked.",
					}).GenerateResponse(w)
					return
				}
			}
			(&response.Response{
				HttpStatus: http.StatusUnauthorized,
				Message:    "token is invalid.",
			}).GenerateResponse(w)
			return
		}

		if !token.Valid {
			(&response.Response{
				HttpStatus: http.StatusUnauthorized,
				Message:    "invalid token claims",
			}).GenerateResponse(w)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", token.Claims.(jwt.MapClaims)["user_id"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func OptionalJwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := string(authHeader)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			validationErr, ok := err.(*jwt.ValidationError)
			if ok {
				if validationErr.Errors == jwt.ValidationErrorExpired {
					next.ServeHTTP(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
			return
		}

		if !token.Valid {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", token.Claims.(jwt.MapClaims)["user_id"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
