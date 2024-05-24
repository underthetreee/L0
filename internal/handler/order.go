package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/underthetreee/L0/internal/model"
)

type OrderGetter interface {
	Get(context.Context, string) (model.Order, error)
}

type OrderHandler struct {
	og OrderGetter
}

func NewOrderHandler(og OrderGetter) *OrderHandler {
	return &OrderHandler{
		og: og,
	}
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	order, err := h.og.Get(context.Background(), id)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		log.Printf("invalid id: %s", id)
		return
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(orderJSON)
}
