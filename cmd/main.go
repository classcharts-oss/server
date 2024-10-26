package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/classcharts-oss/server/pkg/db"
	"github.com/classcharts-oss/server/pkg/router"
	muxHandlers "github.com/gorilla/handlers"
)

const PORT = 4000

func main() {
	r := router.CreateMuxRouter()

	headersOk := muxHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := muxHandlers.AllowedOrigins([]string{"*"})
	methodsOk := muxHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	handler := muxHandlers.CORS(headersOk, originsOk, methodsOk)(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + strconv.Itoa(PORT),
	}

	db.InitDB()

	fmt.Printf("Binding to :%v\n", PORT)

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
