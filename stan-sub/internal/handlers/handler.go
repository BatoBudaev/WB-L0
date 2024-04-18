package handlers

import (
	"fmt"
	"github.com/BatoBudaev/WB-L0/internal/cache"
	"github.com/BatoBudaev/WB-L0/internal/model"
	"html/template"
	"net/http"
	"strconv"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	strId := r.URL.Query().Get("id")
	if strId == "" {
		http.Error(w, "ID не указан", http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(strId)
	order, err := cache.GetOrderById(id)
	if err != nil {
		http.Error(w, "Не удалось получить данные", http.StatusInternalServerError)
		return
	}

	data := model.Data{
		ID:    id,
		Order: order,
	}

	t, err := template.ParseFiles("../internal/html/order.html")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	}
}
