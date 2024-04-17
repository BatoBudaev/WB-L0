package cache

import (
	"github.com/BatoBudaev/WB-L0/internal/database"
	"github.com/BatoBudaev/WB-L0/internal/model"
	"log"
	"sync"
)

var orderCache = make(map[int]model.Order)
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
	}
}

func UpdateCache(data model.Data) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	orderCache[data.ID] = data.Order
}
