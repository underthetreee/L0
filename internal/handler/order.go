package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/underthetreee/L0/internal/model"
)

const (
	ErrInternalServer = "internal server error"
	ErrInvalidInput   = "invalid input"
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
		http.Error(w, ErrInvalidInput, http.StatusBadRequest)
		log.Println(err)
		return
	}

	orderBytes, err := json.Marshal(order)
	if err != nil {
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(orderBytes)
	if err != nil {
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
