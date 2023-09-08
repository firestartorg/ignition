package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/firestart/ignition/pkg/maintenance"
	"net/http"
)

func main() {
	srv := maintenance.NewServer(nil, nil)
	srv.SetLivenessProbe(func(ctx context.Context) (bool, error) {
		return false, fmt.Errorf("error")
	})

	go srv.ListenAndServe(":8080")

	// do other stuff
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		_, _ = w.Write([]byte("Hello, World!"))
	})

	_ = http.ListenAndServe(":8081", router)
}
