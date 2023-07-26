package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	errorJSON(w, "Something went wrong", 400)
}
