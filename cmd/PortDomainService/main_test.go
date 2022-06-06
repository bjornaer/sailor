package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	p "github.com/bjornaer/sailor/cmd/PortDomainService"
	"github.com/stretchr/testify/assert"
)

func TestHelloRoute(t *testing.T) {
	router := p.SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello Sailor! Welcome to the Port Domain Service!", w.Body.String())
}
