package main

import (
	"net/http"
)

func (app *Application) HandleRequest(w http.ResponseWriter, r *http.Request) {
	app.wc.IncrementCount()
	app.writeJSON(w, http.StatusOK, envelope{"message": "request rceived!"}, nil)
}

func (app *Application) HandleCount(w http.ResponseWriter, r *http.Request) {
	reqCount, err := app.wc.Count()
	if err != nil {
		app.logger.Error(err.Error())
		app.errorResponse(w, http.StatusInternalServerError, err.Error())
	}

	app.writeJSON(w, http.StatusOK, envelope{"requests count": reqCount}, nil)
}
