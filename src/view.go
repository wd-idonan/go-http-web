package main


import (
	"net/http"
)

func describeUserInfo(w http.ResponseWriter, r *http.Request) int {
	return writeResponseMsg(w, StatusSuccess, nil)
}