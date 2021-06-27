package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	//IndexHandler(w, r)

	handler := http.HandlerFunc(IndexHandler)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got status %d but wanted %d", w.Code, http.StatusOK)
	}
}
