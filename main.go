package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/auth-microservice/config"
	"github.com/auth-microservice/database"
	"github.com/auth-microservice/models"
	"github.com/auth-microservice/repository"
	"github.com/auth-microservice/server"
	"github.com/auth-microservice/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)
type MessagePattern struct {
	Pattern struct {
		Cmd string `json:"cmd"`
	} `json:"pattern"`
	Data interface{} `json:"data"`
	ID   string      `json:"id"`
}
type TokenResponse struct {
	Token string `json:"token"`
}
type ErrorResponse struct {
	Status string `json:"status"`
}
// func BindRoutes(s server.Server, r *mux.Router) {
// 	// middleware
// 	r.Use(middleware.CheckAuthMiddleWare(s))

// 	// Users
// 	// r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
// 	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
// 	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
// 	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)


// }
func main() {
	// obtner el archivo .env
	err := godotenv.Load(".env")
	if err != nil {
		 log.Fatal("Error loading .env file")
	}
	// declarar e inicializar variablas en base al .env
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	repo, err := database.NewMySQLRepository(DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepository(repo)
	server:= server.NewAuthServer(repo)
	
	log.Println("Starting server on port", PORT)
	// logedUser:=&models.LoginUser{
	// 		Email:       "claudio.doe@example.com",
	// 		Password:    "securepassword",
	// }

	// token,err:=server.Login(context.Background(),logedUser)
	// if err != nil{
	// 	log.Fatal(err)
	// }
	// fmt.Println("Token Del Usuario",token)
	config.SetupRabbitMQ()
	forever := make(chan bool)
	//consumir los mensajes de la cola
	messages, err := config.RMQChannel.Consume(
		config.RMQQueue.Name, // Nombre de la cola
		"",     // Nombre del consumidor
		true,   // Auto-ack (auto acknowledge)
		false,  // Exclusiva
		false,  // No-local
		false,  // No-wait
		nil,    // Argumentos adicionales
	)
	if err != nil {
		log.Fatal(err)
	}	
go func() {
		for msg := range messages {
			// Convertir el cuerpo del mensaje a la estructura Mensaje
			var message MessagePattern
			if err := json.Unmarshal(msg.Body, &message); err != nil {
				log.Printf("Error al decodificar el mensaje: %v", err)
				continue
			}
			switch message.Pattern.Cmd {
			case "SignUp":		
				data, ok := message.Data.(map[string]interface{})["newUser"].(interface{})
				if !ok {
					fmt.Println("Error: Atributo newUser no se encuentra")
					return
				}
				userJSON, err := json.Marshal(data)
				if err != nil {
					fmt.Println("Error al convertir el mapa a JSON:", err)
					return
				}
				var signUpUser models.SignUpUser
				err = json.Unmarshal(userJSON, &signUpUser)
				if err != nil {
					fmt.Println("Error al convertir JSON a SignUpUser:", err)
					return
				}
				
				user,err:=server.SignUp(context.Background(),&signUpUser)
					if err != nil {
					fmt.Println("Error al insertar usuario:", err)
					return
				}
				resultJSON, err := json.MarshalIndent(user, "", "  ")
				if err != nil {
					fmt.Println("Error al convertir a JSON:", err)
					return
				}
				utils.SendResponse(config.RMQChannel, msg, string(resultJSON));
				return 
			case "Login":
				data, ok := message.Data.(map[string]interface{})["logedUser"].(interface{})
				if !ok {
					fmt.Println("Error: Atributo logedUser no se encuentra")
					return
				}
				userJSON, err := json.Marshal(data)
				if err != nil {
					fmt.Println("Error al convertir el mapa a JSON:", err)
					return
				}
				var loginUser models.LoginUser
				err = json.Unmarshal(userJSON, &loginUser)
				if err != nil {
					fmt.Println("Error al convertir JSON a LoginUser:", err)
					return
				}
				token,err:=server.Login(context.Background(),&loginUser)
				if err != nil {
					fmt.Println("Error obtener el token:", err)
					return
				}
				response := TokenResponse{
					Token: token,
				}
				resultJSON, err := json.MarshalIndent(response, "", "  ")
				if err != nil {
					fmt.Println("Error al convertir a JSON:", err)
					return
				}
				utils.SendResponse(config.RMQChannel, msg, string(resultJSON));
			case "Validate":
				tokenString, ok := message.Data.(map[string]interface{})["token"].(string)
				if !ok {
					fmt.Println("Error: Atributo token no se encuentra")
					return
				}
				token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(JWT_SECRET), nil
				})
				if err != nil {
					errorResponse := ErrorResponse{
						Status: "Error parsing token",
					}
					errorJSON, err := json.MarshalIndent(errorResponse, "", "  ")
					if err != nil {
						fmt.Println("Error al convertir a JSON:", err)
						return
					}
					utils.SendResponse(config.RMQChannel, msg, string(errorJSON));
					log.Fatal("Error parsing token:", err)
				}
				if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid{
					user, err := repository.GetUserById(context.Background(), claims.UserId)
					if err!= nil {
						log.Fatal("Error luego de obtener el usuario:", err)
					}
					if err != nil {
						fmt.Println("Error al insertar usuario:", err)
					}
					resultJSON, err := json.MarshalIndent(user, "", "  ")
					if err != nil {
						fmt.Println("Error al convertir a JSON:", err)
					}
					utils.SendResponse(config.RMQChannel, msg, string(resultJSON));
				}else{
					errorResponse := ErrorResponse{
						Status: "invalid credentials",
					}
					errorJSON, err := json.MarshalIndent(errorResponse, "", "  ")
					if err != nil {
						fmt.Println("Error al convertir a JSON:", err)
					}
					utils.SendResponse(config.RMQChannel, msg, string(errorJSON));
				}
				// default:
			// 	fmt.Println("No existe ese método")
			}
		}
	}()
	fmt.Println("Microservicio escuchando. Presiona CTRL+C para salir.")
	<-forever
}
//  ******************************+
// falta hacer la conexión con mysql y no con postgre, probrar todos los metodos primero 