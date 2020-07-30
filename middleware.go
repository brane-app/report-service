package main

import (
	"github.com/imonke/monkebase"

	"net/http"
)

func ReaderMustModerator(request *http.Request) (_ *http.Request, ok bool, code int, _ map[string]interface{}, err error) {
	if request.Method == "POST" {
		ok = true
		return
	}

	var who string = request.Context().Value("requester").(string)
	ok, err = monkebase.IsModerator(who)
	code = 403
	return
}
