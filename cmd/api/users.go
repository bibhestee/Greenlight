package main

import (
	"errors"
	"net/http"

	"github.com/bibhestee/Greenlight/internal/data"
	"github.com/bibhestee/Greenlight/internal/validator"
)

func (app *application) registerUserHandler(res http.ResponseWriter, req *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(res, req, &input)
	if err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(res, req, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(res, req, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(res, req, v.Errors)
		default:
			app.serverErrorResponse(res, req, err)
		}
		return
	}

  app.background(func() {
	  err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	  if err != nil {
	  	app.logger.PrintError(err, nil)
  	}
	})

	err = app.writeJSON(res, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}
