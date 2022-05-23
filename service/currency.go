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

type CurrencyExtraResponse struct {
	Items []CurrencyItemExtra `json:"items"`
}
type CurrencyItemExtra struct {
	CurrencyId string  `json:"currencyId"`
	Value      float32 `json:"value"`
}

type CurrencyService interface {
	Get(req CurrencyRequest) (CurrencyResponse, *utils.AppMess)
}

////////////////////////////////
type CurrencyItemAPI struct {
	Code  string  `json:"code"`
	Value float32 `json:"value"`
}

type LastUpdateAPI struct {
	LastUpdatedAt string `json:"last_updated_at"`
}

type CurrencyPayloadAPI struct {
	Meta LastUpdateAPI              `json:"meta"`
	Data map[string]CurrencyItemAPI `json:"data"`
}
