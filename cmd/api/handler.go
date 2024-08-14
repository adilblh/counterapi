package main

import (
	"net/http"
)

func (app *Application) HandleRequest(w http.ResponseWriter, r *http.Request) {
	app.wc.IncrementCount()
	reqCount, err := app.wc.Count()
	if err != nil {
		app.serverErrorResponse(w, err)
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "request rceived!", "requests count": reqCount}, nil)
}

func (app *Application) HandleCount(w http.ResponseWriter, r *http.Request) {
	reqCount, err := app.wc.Count()
	if err != nil {
		app.serverErrorResponse(w, err)
	}

	app.writeJSON(w, http.StatusOK, envelope{"requests count": reqCount}, nil)
}
