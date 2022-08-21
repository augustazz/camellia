package web

import (
	"github.com/augustazz/camellia/logger"
	"net/http"
)

func StartWebServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	logger.Fatal(http.ListenAndServe(":8080", mux))
}

//http handle func
func ping(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("pong from camellia http server"))
}
