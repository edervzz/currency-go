package utils

import "net/http"

type AppMess struct {
	Code    int
	Message string
}

func NewNotFound(m string) *AppMess {
	return &AppMess{
		Code:    http.StatusNotFound,
		Message: m,
	}
}

func NewInternalError(m string) *AppMess {
	return &AppMess{
		Code:    http.StatusInternalServerError,
		Message: m,
	}
}

func NewBadRequest(m string) *AppMess {
	return &AppMess{
		Code:    http.StatusBadRequest,
		Message: m,
	}
}
