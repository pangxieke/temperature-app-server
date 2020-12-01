package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"temperature/models"
)

type User struct {
	Application
}

func (this *User) Init(rw http.ResponseWriter, r *http.Request, p httprouter.Params) (err error) {
	if err = this.Base.Init(rw, r, p); err != nil {
		return
	}

	this.skipAuthentication = true
	return
}

func (this *User) Login() (err error) {
	params := &struct {
		Mobile string `json:"mobile"`
		Code   string `json:"code"`
	}{}
	if err = this.extractParams(&params); err != nil {
		return
	}
	if params.Mobile == "" {
		err = fmt.Errorf("mobile is empty")
		return
	}
	if params.Code == "" {
		err = fmt.Errorf("code is empty")
		return
	}

	sms := new(models.SMS)
	sms.Mobile = params.Mobile
	if ok := sms.Verify(params.Code); ok == false {
		err = fmt.Errorf("验证码错误或已过期")
		return
	}
	_, err = models.FindUserByMobile(params.Mobile)
	if err != nil {
		err = fmt.Errorf("用户不存在")
		return
	}

	session, err := models.NewSession(params.Mobile)
	if err != nil {
		return
	}

	data := make(map[string]interface{})
	data["token"] = session.Token.String()

	return this.Success(data)
}

func (this *User) SendSMS() (err error) {
	params := &struct {
		Mobile string `json:"mobile"`
	}{}
	if err = this.extractParams(&params); err != nil {
		return
	}
	if params.Mobile == "" {
		err = fmt.Errorf("mobile is empty")
		return
	}

	sms := new(models.SMS)
	sms.Mobile = params.Mobile

	err = sms.Send()
	if err != nil {
		return
	}

	return this.Success(nil)
}
