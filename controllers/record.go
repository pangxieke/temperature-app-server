package controllers

import (
	"fmt"
	"log"
	"temperature/models"
	"time"
)

type Record struct {
	Application
}

func (this *Record) UserInfo()(err error){
	if this.session.Mobile == "" {
		err = fmt.Errorf("请登录")
	}

	user, err := models.FindUserByMobile(this.session.Mobile)
	if err != nil {
		err = fmt.Errorf("用户不存在")
		return
	}


	data := make(map[string]interface{})
	data["username"] = user.Name
	data["company"] = user.Company
	data["face_image"] = user.FaceImage

	return this.Success(data)
}

func (this *Record) List() (err error) {
	params := &struct {
		StartTime int `json:"start_time"`
		EndTime   int `json:"end_time"`
	}{}
	if err = this.extractParams(&params); err != nil {
		return
	}
	user, err := models.FindUserByMobile(this.session.Mobile)
	if err != nil {
		return
	}
	startTime := params.StartTime

	endTime := int(time.Now().Unix())
	if params.EndTime != 0 {
		endTime = params.EndTime
	}
	faceId := user.FaceId
	list, err := models.FindListSN(uint(faceId), startTime, endTime)
	if err != nil {
		return
	}

	s := make(map[uint]interface{})
	if len(list) > 0 {
		for _, val := range list {
			if val.SN == "" {
				continue
			}
			recordList, err := models.FindListBySN(uint(faceId), val.City, startTime, endTime)
			if err != nil {
				log.Printf("FindListBySN err=%s", err)
			}
			data := make([]interface{}, 0)
			if len(recordList) > 0 {
				for _, v := range recordList {
					temp := make(map[string]interface{})
					temp["id"] = v.ID
					temp["time"] = v.CreatedAt
					temp["num"] = v.Num

					data = append(data, temp)
				}
			}

			if len(data) > 0 {
				s[val.City] = data
			}

		}
	}
	return this.Success(s)
}
