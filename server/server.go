package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/auth-microservice/models"
	"github.com/auth-microservice/repository"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	repo repository.Repository
}

func NewAuthServer(repository repository.Repository) *Server {
	return &Server{repo: repository}
}
func (s *Server) SignUp(ctx context.Context, newUser *models.SignUpUser) (*models.User, error) {
	existingUser, err := repository.GetUserByEmail(ctx, newUser.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("User Already Exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := models.User{
		Name:        newUser.Name,
		MiddleName:  newUser.MiddleName,
		Rut:         newUser.Rut,
		PhoneNumber: newUser.PhoneNumber,
		Email:       newUser.Email,
		Password:    string(hashedPassword),
	}
	fmt.Println(user.Email,"antes de insertar")
	id, err := repository.InsertUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Id:          id,
		Name:        user.Name,
		MiddleName:  user.MiddleName,
		Rut:         user.Rut,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}, nil
}
func (s *Server) Login(ctx context.Context, logedUser *models.LoginUser) (string, error) {
	
	user, err := repository.GetUserByEmail(ctx, logedUser.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return "", err
	}
	if user == nil {
		log.Println(err)
		return "", errors.New("Invalid Credentiasl")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logedUser.Password)); err != nil {
		log.Println(err)
		return "",errors.New("Incorrect Password")
	}
	err = godotenv.Load(".env")
	if err != nil {
		 log.Fatal("Error loading .env file")
	}
	JWT_SECRET := os.Getenv("JWT_SECRET")
	claims := models.AppClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
									ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
	},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		log.Fatal(err)
		return "",errors.New("Erro Con token")
	}
	return tokenString,nil
}

// func (b *Broker) Start(binder func (s Server, r *mux.Router)){
// 	b.router=mux.NewRouter()
// 	binder(b,b.router)
// 	repo, err := database.NewMySQLRepository(b.config.DatabaseUrl)
// 	if err!= nil {
// 		log.Fatal(err)
// 	}
// 	repository.SetRepository(repo)
// 	log.Println("Starting server on port", b.Config().Port)
// 	if err:= http.ListenAndServe(b.config.Port,b.router) ; err!=nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }