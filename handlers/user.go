package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/auth-microservice/models"
	"github.com/auth-microservice/repository"
	"github.com/auth-microservice/server"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SignUpLoginRequest struct {
	Name        string `json:"name"`
	MiddleName  string `json:"middlename"`
	Rut         string `json:"rut"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type SignUpResponse struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	MiddleName  string `json:"middlename"`
	Rut         string `json:"rut"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// Verificar si el usuario ya existe en la base de datos
		existingUser, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
			// El usuario ya existe, puedes manejar esto como desees
			http.Error(w, "User alredy exists", http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var user = models.User{
			Name: request.Name,
			MiddleName: request.MiddleName,
			Rut: request.Rut,
			PhoneNumber: request.PhoneNumber,
			Email:    request.Email,
			Password: string(hashedPassword),
		}
		id,err := repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "aplication/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id:    id,
			Name:user.Name,
			MiddleName:user.MiddleName,
			Rut:user.Rut,
			PhoneNumber:user.PhoneNumber,
			Email:user.Email,
		})
	}
}
func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// verifica que usuario exista
		if user == nil {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			log.Println(err)
			http.Error(w, "Incorrect Password", http.StatusUnauthorized)
			return
		}
		//  generar token
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}
func MeHandler(s server.Server) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))	

		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})

		if err != nil {
			log.Println("Error parsing token:", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid{
			user, err := repository.GetUserById(r.Context(), claims.UserId)
			if err!= nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type","application/json")
			json.NewEncoder(w).Encode(user)
		}else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}