package controllers

import (
	"fmt"
	"net/http"
)

type Error struct {
	error
	Status int
	Code   *int
}

func errorStatus(err error) int {
	e, ok := err.(Error)
	if !ok {
		return http.StatusInternalServerError
	}
	return e.Status
}

func errorBody(err error) interface{} {
	r := map[string]interface{}{"message": err.Error()}
	if e, ok := err.(Error); ok && e.Code != nil {
		r["code"] = e.Code
	}
	return r
}

func errorf(status int, code *int, format string, a ...interface{}) error {
	return Error{
		fmt.Errorf(format, a...),
		status,
		code,
	}
}

func newInvalidInputErrorf(format string, a ...interface{}) error {
	return errorf(http.StatusBadRequest, nil, format, a...)
}

func newUnauthorizedErrorf(format string, a ...interface{}) error {
	return errorf(http.StatusUnauthorized, nil, format, a...)
}

func newNotFoundError(name, id interface{}) error {
	return errorf(http.StatusNotFound, nil, "%s not found, id: %v", name, id)
}

func newInvalidParameterError(name, value interface{}) error {
	return newInvalidInputErrorf("Invalid parameter [%s]='%v'", name, value)
}

func newCannotProcessErrorf(format string, a ...interface{}) error {
	return errorf(http.StatusUnprocessableEntity, nil, format, a...)
}

func newCannotProcessErrorfWithCode(code int, format string, a ...interface{}) error {
	return errorf(http.StatusUnprocessableEntity, &code, format, a...)
}
