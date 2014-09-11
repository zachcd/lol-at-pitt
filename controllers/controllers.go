package controllers

type ResponseType int

const (
	JSON ResponseType = iota
	HTML
	// More to be added later?
)

type Response interface{}
type Input interface{}

type Controller struct {
	Response ResponseType
	Route    string
	Input    Input
	Output   Response
	Handler  func(Input) Response
}
