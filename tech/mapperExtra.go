package tech

import (
	"currency-go/service"
)

type MapperExtra struct {
	Result service.CurrencyResponse
}

func (m MapperExtra) MapResponse() service.CurrencyExtraResponse {
	response := service.CurrencyExtraResponse{}
	for _, v := range m.Result.Items {
		item := service.CurrencyItemExtra{
			CurrencyId: v.CurrencyId,
			Value:      v.Value,
		}
		response.Items = append(response.Items, item)
	}
	return response
}
