package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"temperature/models"
)

type Notice struct {
	Application
}

type BaiDuNotice struct {
	Item2        string `json:"item2"`
	ProjectToken string `json:"projectToken"`
	Item1        string `json:"item1"`
	GroupID      string `json:"groupId"`
	Telephone    string `json:"telephone"`
	FaceImage    string `json:"faceImage"`
	OptType      int    `json:"optType"`
	ApplyID      int    `json:"applyId"`
	UUserID      string `json:"uUserId"`
	AppID        int    `json:"appId"`
	Name         string `json:"name"`
	LogID        int64  `json:"logId"`
	Item4        string `json:"item4"` //单位名称
	Item6        string `json:"item6"` //单位编码
	Item7        string `json:"item7"` //部门
}

func (this *Notice) Init(rw http.ResponseWriter, r *http.Request, p httprouter.Params) (err error) {
	if err = this.Base.Init(rw, r, p); err != nil {
		return
	}

	this.skipAuthentication = true
	return
}

func (this *Notice) BaiDuNotice() (err error) {
	var params BaiDuNotice
	if err = this.extractParams(&params); err != nil {
		return
	}
	if params.UUserID == "" {
		return newInvalidParameterError("uUserId", "empty")
	}

	baiDuUserId := strconv.Itoa(params.ApplyID)
	//查询是否face表注册
	faceMod, err := models.FindFaceByBaiDuId(baiDuUserId)
	fmt.Println("faceMod", faceMod)
	if err != nil {
		log.Printf("FindFaceByBaiDuId search err=%s, image=%s", err, params.FaceImage)
		return
	}
	if faceMod.ID == 0 {
		//face表 注册
		faceMod.BaiDuID = baiDuUserId
		faceMod.UserInfo = params.UUserID
		faceMod, err = faceMod.Register()
		if err != nil {
			log.Printf("face Register, err=%s, data=%v", err, faceMod)
		}
	}

	companyId, _ := strconv.Atoi(params.Item6)
	//同一个公司，多次注册时，更新
	//check user_no
	user, err := models.FindUserByNOAndCompany(params.UUserID, uint(companyId))
	if err != nil {
		log.Printf("models.FindUserByNOAndCompany, err=%s, phone=%s, company_id=%d", err, params.UUserID, companyId)
		return
	}

	// check AppId 18666841
	// Item1 省份证号码
	// Item2 年龄
	// item6 公司编号
	// Item4 公司名
	log.Printf("baidu notice: %+v", params)

	user.UserNO = params.UUserID
	user.Tel = params.Telephone
	user.IdNum = params.Item1
	age, _ := strconv.Atoi(params.Item2)
	user.Age = uint(age)
	user.CompanyId = uint(companyId) //Check
	user.Company = params.Item4
	user.FaceImage = params.FaceImage
	user.AppId = uint(params.AppID)
	user.Group = params.GroupID
	user.Name = params.Name
	user.FaceId = faceMod.ID
	user.Department = params.Item7

	err = user.Save()
	if err != nil {
		log.Printf("user save, err=%s, data=%v", err, user)
		return
	}

	return this.Success(nil)
}
