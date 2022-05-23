package app

import (
	"currency-go/service"
	"currency-go/tech"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CurrencyHandler struct {
	service service.CurrencyService
}

func (h CurrencyHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	sdtfrom := r.URL.Query().Get("finit")
	sdtto := r.URL.Query().Get("fend")
	ctype := r.URL.Query().Get("ctype")

	var dfrom, dto *time.Time

	if sdtfrom != "" {
		df, err := time.Parse("2006-01-02T15:04:05", sdtfrom)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		dfrom = &df
	}

	if sdtto != "" {
		dt, err := time.Parse("2006-01-02T15:04:05", sdtto)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		dto = &dt
	}

	if dfrom != nil && dto != nil && dfrom.After(*dto) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("TO, cannot be before FROM")
		return
	}

	req := service.CurrencyRequest{
		CurrencyId: mux.Vars(r)["id"],
		From:       sdtfrom,
		To:         sdtto,
	}

	res, appMess := h.service.Get(req)

	var final interface{}

	switch ctype {
	case "normal":
		normal := tech.MapperNormal{
			Result: res,
		}
		final = normal.MapResponse()
	case "extra":
		extra := tech.MapperExtra{
			Result: res,
		}
		final = extra.MapResponse()
	}

	if appMess != nil {
		w.WriteHeader(appMess.Code)
		json.NewEncoder(w).Encode(appMess.Message)
		return
	}

	json.NewEncoder(w).Encode(final)
}

func (h CurrencyHandler) GetCurrencyAPISimulate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	res := service.CurrencyPayloadAPI{
		Meta: service.LastUpdateAPI{
			LastUpdatedAt: "2022-05-15T23:59:59Z",
		},
		Data: map[string]service.CurrencyItemAPI{
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
	time.Sleep(600 * time.Millisecond)
	json.NewEncoder(w).Encode(res)
}

func NewCurrencyHandler(service service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		service,
	}
}
