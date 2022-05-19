package utils

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	mess := NewNotFound("no items")
	assert.Equal(t, http.StatusNotFound, mess.Code)

	mess = NewBadRequest("err")
	assert.Equal(t, http.StatusBadRequest, mess.Code)

	mess = NewInternalError("err")
	assert.Equal(t, http.StatusInternalServerError, mess.Code)
}
