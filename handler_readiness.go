package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	writeJson(w, 200, struct{}{})
}
