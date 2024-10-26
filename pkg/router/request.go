package router

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/CommunityCharts/CCModels/shared"
	"github.com/classcharts-oss/server/pkg/db"
	"github.com/golang-jwt/jwt"
)

func RequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Log the request for debugging purposes
			fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		}()

		next.ServeHTTP(w, r)
	})
}

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				response := shared.NewErrorfulResponse(fmt.Sprintf("%v", err))

				w.WriteHeader(http.StatusInternalServerError)
				response.Write(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func returnAuthErr(msg string, w http.ResponseWriter) {
	response := shared.NewExpiredResponse(msg)

	w.WriteHeader(http.StatusOK)
	response.Write(w)
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has a valid JWT token
		// If not, panic with an error
		// If it does, continue to the next handler
		// For now, we'll just assume the token is valid
		auth := r.Header.Get("Authorization")
		if auth == "" {
			returnAuthErr("You are not logged in. [Token missing.]", w)
			return
		}

		tokenString := strings.TrimPrefix(auth, "Basic ")
		if tokenString == "" {
			returnAuthErr("You are not logged in. [Token missing.]", w)
			return
		}

		claims := &db.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			returnAuthErr("You are not logged in. ["+err.Error()+"]", w)
			return
		}
		if !token.Valid {
			returnAuthErr("You are not logged in. [Token invalid.]", w)
			return
		}

		// Set the student in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "student", db.GetStudentByID(claims.StudentID))

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
