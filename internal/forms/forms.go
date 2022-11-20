package forms

import "net/url"

type From struct {
	url.Values
	Errors errors
}

//New init a form struct
func New(data url.Values) *From {
	return &From{
		data,
		errors(map[string][]string{}),
	}
}
