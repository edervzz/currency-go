package tech

import (
	"currency-go/service"
)

type MapperNormal struct {
	Result service.CurrencyResponse
}

func (m MapperNormal) MapResponse() service.CurrencyResponse {
	response := service.CurrencyResponse{}
	for _, v := range m.Result.Items {
		item := service.CurrencyItem{
			CurrencyId: v.CurrencyId,
			UpdatedAt:  v.UpdatedAt,
			Value:      v.Value,
		}
		response.Items = append(response.Items, item)
	}
	return response
}
