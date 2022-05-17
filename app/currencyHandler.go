package app

import (
	"currency-go/service"
	"encoding/json"
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

	json.NewEncoder(w).Encode(res)
}

func (h CurrencyHandler) GetCurrencyAPISimulate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := service.CurrencyAPIPayload{
		Meta: service.APILastUpdate{
			LastUpdatedAt: "2022-05-15T23:59:59Z",
		},
		Data: map[string]service.CurrencyAPIItem{
			"AED": {
				Code:  "AED",
				Value: 3.67311,
			},
			"AFN": {
				Code:  "AFN",
				Value: 88.00199,
			},
			"ALL": {
				Code:  "ALL",
				Value: 15.85339,
			},
		},
	}

	json.NewEncoder(w).Encode(res)
}

func NewCurrencyHandler(service service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		service,
	}
}
