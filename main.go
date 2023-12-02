package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/auth-microservice/handlers"
	"github.com/auth-microservice/middleware"
	"github.com/auth-microservice/server"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func BindRoutes(s server.Server, r *mux.Router) {
	// middleware
	r.Use(middleware.CheckAuthMiddleWare(s))

	// Users
	// r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)


}
func main() {
	// obtner el archivo .env
	err := godotenv.Load(".env")
	if err != nil {
		 log.Fatal("Error loading .env file")
	}
	// declarar e inicializar variablas en base al .env
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_USER := os.Getenv("DATABASE_USER")

	s, err:= server.NewServer(context.Background(), &server.Config{
			Port: PORT,        
			JWTSecret: JWT_SECRET, 
			DatabaseUrl: DATABASE_USER, 
	})
	if err!= nil {
		log.Fatal(err)
	}
	s.Start(BindRoutes)
	
}

//  ******************************+
// falta hacer la conexión con mysql y no con postgre, probrar todos los metodos primero 