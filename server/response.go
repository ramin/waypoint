package server

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Response struct {
	Writer http.ResponseWriter `json:"-"`
	Data   interface{}         `json:"data"`

	Error *Error `json:"error,omitempty"`
	Meta  Meta   `json:"meta"`
}

type Meta struct {
	Version string `json:"version,omitempty"`
	Token   string `json:"token,omitempty"`
}

type Error struct {
	Message string

	Info []string
}

type TimeStartKey struct{}

func NewResponder(w http.ResponseWriter) Response {
	return Response{
		Writer: w,
		Error:  nil,
		Meta:   Meta{},
	}
}

func (r *Response) AddError(err Error) *Response {
	r.Error = &err
	return r
}

func (r *Response) Empty() {
	r.Writer.WriteHeader(http.StatusAccepted)
	_, err := r.Writer.Write(nil)
	if err != nil {
		log.Error(err)
	}
}

func (r *Response) NotFound() {
	r.respond(http.StatusNotFound)
}

func (r *Response) Success() {
	r.respond(http.StatusOK)
}

func (r *Response) UnexpectedError() {
	r.respond(http.StatusInternalServerError)
}

func (r *Response) BadRequest() {
	r.respond(http.StatusBadRequest)
}

func (r *Response) NoContent() {
	r.respond(http.StatusNoContent)
}

func (r *Response) NoContentWithError(err error) {
	log.Error(err)
	r.Error = &Error{Message: err.Error()}
	r.respond(http.StatusNoContent)
}

func (r *Response) ForbiddenWithError(err error) {
	log.Error(err)
	r.Error = &Error{Message: err.Error()}
	r.respond(http.StatusForbidden)
}

func (r *Response) BadRequestWithError(err error) {
	log.Error(err)
	r.Error = &Error{Message: err.Error()}
	r.respond(http.StatusBadRequest)
}

func (r *Response) UnexpectedWithError(err error) {
	log.Error(err)
	r.Error = &Error{Message: err.Error()}
	r.respond(http.StatusInternalServerError)
}

func (r *Response) Unauthorized() {
	r.Error = &Error{Message: "Unauthorized"}
	r.respond(http.StatusUnauthorized)
}

func (r *Response) WithHeader(header int) {
	r.respond(header)
}

func (r *Response) respond(header int) {
	data, err := json.Marshal(r)

	if err != nil {
		log.Error(err)
	}

	r.Writer.Header().Set("Content-Type", "application/json")

	r.Writer.WriteHeader(header)
	_, err = r.Writer.Write(data)
	if err != nil {
		log.Error(err)
	}
}

type Status struct {
	Additional string `json:"additional,omitempty"`
	Message    string `json:"message"`
}
