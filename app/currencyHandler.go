package app

import (
	"currency-go/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type CurrencyHandler struct {
	service service.CurrencyService
}

func (h CurrencyHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	req := service.CurrencyRequest{
		CurrencyId: mux.Vars(r)["id"],
		From:       r.URL.Query().Get("finit"),
		To:         r.URL.Query().Get("fend"),
	}
	res, appMess := h.service.Get(req)
	if appMess != nil {
		json.NewEncoder(w).Encode(appMess.Message)
		w.WriteHeader(appMess.Code)
		return
	}
	fmt.Println(res)
	json.NewEncoder(w).Encode(res)
}

func NewCurrencyHandler(service service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		service,
	}
}
