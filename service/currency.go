package service

import "currency-go/utils"

type CurrencyRequest struct {
	CurrencyId string `json:"currencyId"`
	From       string `json:"from"`
	To         string `json:"to"`
}

type CurrencyResponse struct {
	Items []CurrencyItem `json:"items"`
}

type CurrencyItem struct {
	CurrencyId string  `json:"currencyId"`
	UpdatedAt  string  `json:"from"`
	Value      float32 `json:"value"`
}

type CurrencyService interface {
	Get(req CurrencyRequest) (CurrencyResponse, *utils.AppMess)
}
