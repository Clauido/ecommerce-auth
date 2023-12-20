package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RMQConnection *amqp.Connection
var RMQChannel *amqp.Channel
var RMQQueue amqp.Queue

func SetupRabbitMQ() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	rabbitMQUser := os.Getenv("RABBITMQ_USERNAME")
	rabbitMQPass := os.Getenv("RABBITMQ_PASSWORD")
	rabbitMQHost := os.Getenv("RABBITMQ_HOST")
	rabbitMQPort := os.Getenv("RABBITMQ_PORT")
	rabbitMQQueue := os.Getenv("RABBITMQ_QUEUE")

	if rabbitMQUser == "" || rabbitMQPass == "" || rabbitMQHost == "" || rabbitMQPort == "" || rabbitMQQueue == "" {
		log.Fatal("Ninguna variable de entorno puede estar vacía")
	}

	RMQConnection, err = amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPass, rabbitMQHost, rabbitMQPort))
	if err != nil {
		log.Fatal(err)
	}

	RMQChannel, err = RMQConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	RMQQueue, err = RMQChannel.QueueDeclare(
		rabbitMQQueue, // Nombre de la cola
		false,         // Durable
		false,         // Eliminar cuando no haya consumidores
		false,       // Exclusiva
		false,         // No esperar confirmación de servidor
		nil,           // Argumentos adicionales
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexion a RabbitMQ completada con exito.")
}