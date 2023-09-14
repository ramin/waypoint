package server

import (
	"net/http"
)

// Status handles health checking the http server
// and reporting if the service is operating, without
// validating anything else (db connections etc)
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	resp := NewResponder(w)
	resp.Data = nil

	resp.Data = Status{
		Message: "ok",
	}

	resp.Success()
}
