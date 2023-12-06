package main

import (
	"github.com/julienschmidt/httprouter"
	"gitlab.com/firestart/ignition/x/application/extensions/http"
	"gitlab.com/firestart/ignition/x/application/presets"
	nethttp "net/http"
)

func main() {
	app := presets.NewHttpApp("example")

	http.MustAddGetRoute(app, "/", func(w nethttp.ResponseWriter, r *nethttp.Request, params httprouter.Params) {
		// Do something
		w.WriteHeader(nethttp.StatusOK)
	})

	app.Run()
}
