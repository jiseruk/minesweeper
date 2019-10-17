package handler

import (
	"net/http"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/apis"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := apis.GetRouter()
	router.ServeHTTP(w, r)
}
