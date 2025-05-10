package initialize

import (
	"log"

	"github.com/anonystick/go-ecommerce-backend-api/global"
	"github.com/segmentio/kafka-go"
)

// Init kafka Producer
var KafkaProducer *kafka.Writer

func InitKafka() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:19094"),
		Topic:    "otp-auth-topic", // topic
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %v", err)
	}
}
