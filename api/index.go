package handler

import (
	"fmt"
	"net/http"

	"github.com/jiseruk/minesweeper/cmd/minesweeper/apis"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(r.RequestURI)
	router := apis.GetRouter()
	router.ServeHTTP(w, r)
}
