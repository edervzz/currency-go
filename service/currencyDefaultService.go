package service

import (
	"currency-go/domain"
	"currency-go/utils"
)

type DefaultCurrencyService struct {
	repo domain.CurrencyRepository
}

func (s DefaultCurrencyService) Get(req CurrencyRequest) (CurrencyResponse, *utils.AppMess) {

	response := CurrencyResponse{}

	result, err := s.repo.Get(req.CurrencyId, req.From, req.To)
	if err != nil {
		appMess := utils.NewNotFound(err.Error())
		return CurrencyResponse{}, appMess
	}

	if len(result) == 0 {
		appMess := utils.NewNotFound("no items")
		return CurrencyResponse{}, appMess
	}

	for _, v := range result {
		item := CurrencyItem{
			CurrencyId: v.CurrencyId,
			UpdatedAt:  v.UpdatedAt,
			Value:      v.Value,
		}
		response.Items = append(response.Items, item)
	}

	return response, nil

}

func NewDefaultCurrencyService(repo domain.CurrencyRepository) *DefaultCurrencyService {
	return &DefaultCurrencyService{
		repo,
	}
}
