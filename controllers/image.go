package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"temperature/config"
	"temperature/models"
	"temperature/views"
)

type Image struct {
	Application
}

func (this *Image) Init(rw http.ResponseWriter, r *http.Request, p httprouter.Params) (err error) {
	if err = this.Base.Init(rw, r, p); err != nil {
		return
	}

	this.skipAuthentication = true
	return
}

func (this *Image) Create() (err error) {
	params := &struct {
		SN string `json:"sn"`
	}{}
	if err = this.extractParams(&params); err != nil {
		return
	}
	if params.SN == "" {
		err = fmt.Errorf("sn is empty")
		return
	}

	imageMod := models.Image{
		SN: params.SN,
	}

	err = imageMod.Save()
	if err != nil {
		return
	}
	url := fmt.Sprintf("temperature/%s/%d/preview.jpg", config.Env, imageMod.ID)

	imageMod.Url = url
	if err = imageMod.Save(); err != nil {
		return
	}

	putUrl := models.PutUrlFor(url, 3000, "image/jpeg")

	resp, err := views.NewCreate(imageMod.ID, putUrl)
	return this.Success(resp)
}
