package main

import (
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"time"
)

func main() {
	sc, err := stan.Connect(
		"test-cluster",
		"publisher",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer sc.Close()

	jsonFilePath := "assets/model-1.json"

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatalf("Не удалось прочитать файл: %v", err)
	}

	for {
		err = sc.Publish("channel-1", jsonData)
		if err != nil {
			log.Fatalf("Не удалось отправить сообщение: %v", err)
		}
		log.Println("Сообщение отправлено")
		time.Sleep(10 * time.Second)
	}
}
