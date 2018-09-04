package main

import "net/http"

func respondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {}
