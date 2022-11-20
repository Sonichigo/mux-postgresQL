// main.go

package main

import (
	"net/http"
	"os"
	"github.com/keploy/go-sdk/integrations/kmux"

	"github.com/gorilla/mux"
	"github.com/keploy/go-sdk/keploy"
)

func main() {
	a := App{}
	r := mux.NewRouter()
	port := "8080"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "my-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:6789/api",
		},
	})
	
	r.Use(kmux.MuxMiddleware(k))
	http.ListenAndServe(":"+port, r)
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":8010")
}
