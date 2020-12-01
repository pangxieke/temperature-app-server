package controllers

import (
	"strconv"
	"strings"
	"temperature/models"

	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Application struct {
	Base
	skipAuthentication bool
	session            *models.Session
}

func (a *Application) authenticate(r *http.Request) (err error) {
	authorization := strings.TrimSpace(r.Header.Get("Authorization"))
	if authorization == "" {
		return newUnauthorizedErrorf("missing Authorization")
	}
	s := strings.Split(authorization, " ")
	var token string
	if len(s) == 1 {
		token = s[0]
	} else if len(s) == 2 {
		var typ string
		typ, token = s[0], s[1]
		if typ != "Bearer" {
			return newUnauthorizedErrorf("invalid Authorization type: `%s`, only support Bearer", s[0])
		}
	} else {
		return newUnauthorizedErrorf("Authorization error")
	}

	session, err := models.LoadSession(token)
	if err != nil || session == nil {
		return newUnauthorizedErrorf("Invalid Authorization")
	}
	a.session = session
	return
}

func (a *Application) Init(rw http.ResponseWriter, r *http.Request, p httprouter.Params) (err error) {
	if err = a.Base.Init(rw, r, p); err != nil {
		return
	}
	if !a.skipAuthentication {
		if err = a.authenticate(r); err != nil {
			return
		}
	}
	return
}

type CommonResp struct {
	Code      int         `json:"code"`
	RequestId string      `json:"requestId"`
	Message   string      `json:"msg"`
	Data      interface{} `json:"data"`
}

//错误返回
func (a *Application) Error(err error) {
	log.Printf("err = %+v\n", err)
	//var code int
	//if e, ok := err.(Error); ok && e.Code != nil {
	//	code = *e.Code
	//} else {
	//	code = 1
	//}
	r := CommonResp{
		Code:      1,
		Message:   err.Error(),
		RequestId: a.UUID,
		Data:      nil,
	}
	_ = a.respondJson(r, errorStatus(err))
}

func (a *Application) Success(data interface{}) error {
	r := CommonResp{
		Code:      0,
		Message:   "success",
		RequestId: a.UUID,
		Data:      data,
	}
	return a.respondJson(r, http.StatusOK)
}

func (a *Application) respondJson(response interface{}, statusCode ...int) error {
	if len(statusCode) > 0 {
		a.ResponseWriter.WriteHeader(statusCode[0])
	}

	j, err := json.Marshal(response)
	if err != nil {
		return err
	}
	a.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	a.ResponseWriter.Write(j)
	return nil
}

func (a Application) extractParams(params interface{}) (err error) {
	err = json.Unmarshal([]byte(a.RequestBody), params)
	if err != nil {
		log.Printf("invalid body format, err = %+v\n", err)
		return newInvalidInputErrorf("invalid body format")
	}
	return
}

func (a Application) getPathParam(name string) (res string) {
	for _, p := range a.Params {
		if p.Key == name {
			return p.Value
		}
	}
	return
}

func (a Application) getID(name string) (res int, err error) {
	id := a.getPathParam("figure_id")
	if id == "" {
		return 0, newInvalidParameterError("figure_id", "")
	}
	res, _ = strconv.Atoi(id)
	return
}

func (a Application) getQuery(name string) (res string, err error) {
	values, ok := a.Request.URL.Query()[name]
	if !ok || len(values[0]) < 1 {
		err = newInvalidParameterError(name, "missing")
		return
	}
	res = values[0]
	return
}

func (a Application) getQueryInt(name string, defaultValue ...int) (res int, err error) {
	value, err := a.getQuery(name)
	if err != nil {
		return
	}
	if value == "" {
		if len(defaultValue) > 1 {
			return defaultValue[0], nil
		}
		return 0, newInvalidParameterError(name, "")
	}
	res, err = strconv.Atoi(value)
	if err != nil {
		return 0, newInvalidParameterError(name, err.Error())
	}
	return
}

func (a Application) getLimitAndOffset() (limit, offset int, err error) {
	offset, err = a.getQueryInt("offset", 0)
	if err != nil {
		return
	}
	limit, err = a.getQueryInt("limit", 24)
	if err != nil {
		return
	}
	return
}
