package main

import (
	"encoding/json"
	"github.com/BatoBudaev/WB-L0/internal/cache"
	"github.com/BatoBudaev/WB-L0/internal/database"
	"github.com/BatoBudaev/WB-L0/internal/handlers"
	"github.com/BatoBudaev/WB-L0/internal/model"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

func main() {
	db, err := database.InitDB("postgres", "1", "orders_db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к PostgreSQL: %v", err)
	}
	defer db.Close()

	cache.InitCacheFromDb(db)

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
		var data model.Data
		err := json.Unmarshal(msg.Data, &data.Order)
		if err != nil {
			log.Printf("Не удалось разобрать JSON: %v", err)
			return
		}

		err = model.ValidateOrder(data.Order)
		if err != nil {
			log.Printf("Не удалось проверить данные: %v", err)
			return
		}

		cache.UpdateCache(data)

		err = db.InsertOrder(data.Order)
		if err != nil {
			log.Fatalf("Не удалось вставить данные: %v", err)
		}
		log.Println("Данные успешно вставлены")
	}
	_, err = sc.Subscribe("channel-1", cb, stan.DurableName("durable-1"))
	if err != nil {
		log.Fatalf("Не удалось подписаться: %v", err)
	}

	http.HandleFunc("/order", handlers.OrderHandler)
	log.Println("Сервер запущен на порту 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
