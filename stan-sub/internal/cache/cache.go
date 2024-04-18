package cache

import (
	"fmt"
	"github.com/BatoBudaev/WB-L0/internal/database"
	"github.com/BatoBudaev/WB-L0/internal/model"
	"log"
	"sync"
)

var orderCache = make(map[int]model.Order)
var idCache = 1
var cacheMutex = &sync.Mutex{}

func InitCacheFromDb(db *database.DB) {
	data, err := db.GetData()
	if err != nil {
		log.Fatalf("Не удалось загрузить данныые из базы данных: %v", err)
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	for _, d := range data {
		orderCache[d.ID] = d.Order
		idCache = d.ID
	}
}

func UpdateCache(data model.Data) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	orderCache[idCache] = data.Order
	idCache++
}

func GetOrderById(id int) (model.Order, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	order, ok := orderCache[id]
	if !ok {
		return model.Order{}, fmt.Errorf("заказ с ID %d не найден", id)
	}

	return order, nil
}
