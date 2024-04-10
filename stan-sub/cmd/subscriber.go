package main

import (
	"encoding/json"
	"github.com/BatoBudaev/WB-L0/internal/database"
	"github.com/BatoBudaev/WB-L0/internal/model"
	"github.com/nats-io/stan.go"
	"log"
)

func main() {
	db, err := database.InitDB("postgres", "1", "orders_db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к PostgreSQL: %v", err)
	}
	defer db.Close()

	sc, err := stan.Connect(
		"test-cluster",
		"subscriber",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	log.Printf("Подключен к кластеру: test-cluster")
	defer sc.Close()

	cb := func(msg *stan.Msg) {
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("Не удалось разобрать JSON: %v", err)
			return
		}

		err = model.ValidateOrder(order)
		if err != nil {
			log.Printf("Не удалось проверить данные: %v", err)
			return
		}

		err = db.InsertOrder(order)
		if err != nil {
			log.Fatalf("Не удалось вставить данные: %v", err)
		}
		log.Println("Данные успешно вставлены")
	}
	_, err = sc.Subscribe("channel-1", cb, stan.DurableName("durable-1"))
	if err != nil {
		log.Fatalf("Не удалось подписаться: %v", err)
	}

	select {}
}