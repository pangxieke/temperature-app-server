package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"net/url"
	"strings"
	"temperature/models"
	"time"
)

type Report struct {
	Application
}

func (this *Report) Init(rw http.ResponseWriter, r *http.Request, p httprouter.Params) (err error) {
	if err = this.Base.Init(rw, r, p); err != nil {
		return
	}

	this.skipAuthentication = true
	return
}

func (this *Report) Create() (err error) {
	params := &struct {
		FaceImage string  `json:"face_image"`
		Num       float64 `json:"num"`
		SN        string  `json:"sn"`
	}{}
	if err = this.extractParams(&params); err != nil {
		return
	}
	if params.FaceImage == "" {
		err = fmt.Errorf("face_image is empty")
		return
	}
	if params.Num == 0 {
		err = fmt.Errorf("num is empty")
		return
	}

	var companyId uint
	devMod := new(models.Device)
	devInfo, err := devMod.FindDevBySN(params.SN)
	if err != nil {
		log.Printf("models.FindDevBySN err=%s, sn=%s", err, params.SN)
	}
	if devInfo == nil {
		log.Printf("sn is not found, sn=%s", params.SN)
	} else {
		companyId = devInfo.StoreId
	}

	//image为putUrl路径，需要解析从新加密
	newUrl, err := url.Parse(params.FaceImage)
	if err != nil {
		fmt.Printf("url parse err=%s, url=%s", err, params.FaceImage)
	}
	faceImagePath := newUrl.Path
	faceImagePath = strings.Trim(faceImagePath, "/")
	faceImageNew := models.GetUrlFor(faceImagePath, 3600)

	// 查询用户信息
	baiDuClient := new(models.BaiDu)
	baiDuClient.Init()
	baiDuUserId, err := baiDuClient.Search(faceImageNew)

	if err != nil {
		log.Printf("baidu search err=%s, image=%s", err, params.FaceImage)
	}

	// 注册百度人脸库
	fmt.Println("baiDuUserId:", baiDuUserId)
	//未注册
	if baiDuUserId == "" {
		baiDuUserId = models.CreatRandBaiDuId()
		//注册百度人脸库
		baiDuClient.Add(faceImageNew, baiDuUserId)
	}

	//查询是否face表注册
	faceMod, err := models.FindFaceByBaiDuId(baiDuUserId)
	//fmt.Println("faceMod", faceMod)
	if err != nil {
		log.Printf("FindFaceByBaiDuId search err=%s, image=%s", err, params.FaceImage)
		return
	}
	if faceMod.ID == 0 {
		//face表 注册
		faceMod.BaiDuID = baiDuUserId
		faceMod.UserInfo = ""
		faceMod, err = faceMod.Register()
		if err != nil {
			log.Printf("face Register, err=%s, data=%v", err, faceMod)
		}
	}

	var userInfo models.User
	var userType uint
	var username string
	var userId uint

	var recordList []models.Record
	if faceMod.ID != 0 {
		// 查询用户
		userInfo, err = models.FindUserByFaceIdAndCompany(faceMod.ID, companyId)
		if err != nil {
			log.Printf("models.FindUserByNO err=%s, face_id=%d", err, faceMod.ID)
		}
		if userInfo.ID != 0 {
			//已注册
			userType = 1
			username = userInfo.Name
			userId = userInfo.ID
		} else {
			userInfo, err = models.FindUserByFaceId(faceMod.ID)
			if userInfo.ID != 0 {
				username = userInfo.Name
				userId = userInfo.ID
			}
		}

		//查询最近记录
		//lastRecord, err := models.FindLastByFaceId(userInfo.FaceId)
		//if err != nil {
		//	log.Printf("models.FindLastByUserId err=%s, id=%d", err, userInfo.ID)
		//}
		//查询历史记录
		recordList, err = models.FindListByFaceId(faceMod.ID)
		if err != nil {
			log.Printf("models.FindLastByUserId err=%s, id=%d", err, userInfo.ID)
		}
		fmt.Println(recordList)
	}

	record := new(models.Record)
	record.FaceImage = faceImagePath
	record.Num = params.Num
	record.SN = params.SN
	if devInfo != nil && devInfo.StoreId != 0 {
		record.StoreId = devInfo.StoreId // check sn
	}

	record.UserId = userId //check face image
	record.Type = userType // 0为访客，1为注册用户
	record.FaceId = faceMod.ID

	err = record.Save()
	if err != nil {
		return
	}

	res := make(map[string]interface{})
	res["id"] = record.ID
	res["now_time"] = time.Now().Unix()
	res["name"] = username
	res["type"] = userType //0为访客，1为注册用户
	//res["last_num"] = lastNum
	//res["last_time"] = lastTime

	list := make([]interface{}, 0)
	if len(recordList) > 0 {
		for _, record := range recordList {
			temp := make(map[string]interface{})
			temp["time"] = record.CreatedAt
			temp["num"] = record.Num
			temp["id"] = record.ID
			list = append(list, temp)
		}
	}

	res["list"] = list

	return this.Success(res)

}
