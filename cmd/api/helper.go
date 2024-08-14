package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (app *Application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())

	message := "the server encountered an issue and was unable to process your request."
	app.errorResponse(w, http.StatusInternalServerError, message)
}

func (app *Application) errorResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}
