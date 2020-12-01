package models

import (
	"encoding/json"
	"fmt"
	"log"
	"temperature/config"
)

type TokenResponse struct {
	Error       string `json:"error,omitempty"`
	ErrorDesc   string `json:"error_description,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}

var BaiDuFaceGroup = "temperature"

type BaiDu struct {
	Token string
}

func (this *BaiDu) Init() {
	this.Token, _ = this.GetToken()
}

func (this *BaiDu) GetToken() (token string, err error) {
	reqUrl := fmt.Sprintf(
		"%s/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		config.BaiDu.Host, config.BaiDu.ApiKey, config.BaiDu.SecretKey)
	resp, err := HttpPost(reqUrl, nil, nil)
	log.Printf("baidu get token, url=%s, resp=%s, err=%s", reqUrl, resp, err)
	if err != nil {
		return
	}
	var tokenResp TokenResponse
	err = json.Unmarshal(resp, &tokenResp)
	if err != nil {
		return
	}

	token = tokenResp.AccessToken
	return
}

func (this *BaiDu) Detect(imgUrl string) (respData interface{}, err error) {

	reqUrl := fmt.Sprintf("%s/rest/2.0/face/v3/detect?access_token=%s", config.BaiDu.Host, this.Token)
	body := make(map[string]interface{})
	body["image"] = imgUrl
	body["image_type"] = "URL"
	body["face_field"] = "age,beauty,expression,face_shape,gender,glasses,landmark,landmark150,race,quality,eye_status,emotion,face_type"
	bodyData, _ := json.Marshal(body)
	resp, err := HttpPost(reqUrl, bodyData, nil)
	log.Printf("baidu detect, url=%s, resp=%s, err=%s", reqUrl, resp, err)
	if err != nil {
		return
	}
	type detectStr struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		LogID     int64  `json:"log_id"`
		Timestamp int    `json:"timestamp"`
		Cached    int    `json:"cached"`
		Result    struct {
			FaceNum  int `json:"face_num"`
			FaceList interface {
			} `json:"face_list"`
		} `json:"result"`
	}
	var detectRes detectStr
	err = json.Unmarshal(resp, &detectRes)
	if err != nil {
		return
	}
	if detectRes.ErrorCode == 0 {
		respData = detectRes.Result.FaceList
	}
	return
}

func (this *BaiDu) Search(imgUrl string) (baiDuUserId string, err error) {
	reqUrl := fmt.Sprintf("%s/rest/2.0/face/v3/search?access_token=%s", config.BaiDu.Host, this.Token)
	body := make(map[string]interface{})
	body["image"] = imgUrl
	body["image_type"] = "URL"
	body["group_id_list"] = BaiDuFaceGroup
	body["max_user_num"] = 5

	bodyData, _ := json.Marshal(body)
	resp, err := HttpPost(reqUrl, bodyData, nil)
	log.Printf("baidu detect, url=%s, resp=%s, err=%s", reqUrl, resp, err)
	if err != nil {
		return
	}

	type searchStr struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		LogID     int64  `json:"log_id"`
		Timestamp int    `json:"timestamp"`
		Cached    int    `json:"cached"`
		Result    struct {
			FaceToken string `json:"face_token"`
			UserList  []struct {
				GroupID  string  `json:"group_id"`
				UserID   string  `json:"user_id"`
				UserInfo string  `json:"user_info"`
				Score    float64 `json:"score"`
			} `json:"user_list"`
		} `json:"result"`
	}

	var searchRes searchStr
	err = json.Unmarshal(resp, &searchRes)
	if err != nil {
		return
	}
	if searchRes.ErrorCode == 0 {
		userList := searchRes.Result.UserList
		//Score 用户的匹配得分，推荐阈值80分
		if len(userList) > 0 && userList[0].UserID != "" && userList[0].Score > 80 {
			baiDuUserId = userList[0].UserID
		}
		// 优先匹配user已注册用户
		if len(userList) > 0 {
			for _, val := range userList {
				if val.UserID != "" && val.Score > 80 {
					user, err := FindUserByBaiDuId(val.UserID)
					if err != nil {
						log.Printf("FindUserByBaiDuId err=%s, baidu_id=%s", err, val.UserID)
					}
					if user.ID != 0 {
						baiDuUserId = val.UserID
						break
					}
				}
			}
		}
	}
	return
}

func (this *BaiDu) Add(imgUrl string, userIdStr string) (tel string, err error) {
	reqUrl := fmt.Sprintf("%s/rest/2.0/face/v3/faceset/user/add?access_token=%s", config.BaiDu.Host, this.Token)
	body := make(map[string]interface{})
	body["image"] = imgUrl
	body["image_type"] = "URL"
	body["group_id"] = BaiDuFaceGroup
	body["user_id"] = userIdStr

	bodyData, _ := json.Marshal(body)
	resp, err := HttpPost(reqUrl, bodyData, nil)
	log.Printf("baidu detect, url=%s, resp=%s, err=%s", reqUrl, resp, err)
	if err != nil {
		return
	}

	type searchStr struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		LogID     int64  `json:"log_id"`
		Timestamp int    `json:"timestamp"`
		Cached    int    `json:"cached"`
		Result    struct {
			FaceToken string `json:"face_token"`
			UserList  []struct {
				GroupID  string  `json:"group_id"`
				UserID   string  `json:"user_id"`
				UserInfo string  `json:"user_info"`
				Score    float64 `json:"score"`
			} `json:"user_list"`
		} `json:"result"`
	}

	var searchRes searchStr
	err = json.Unmarshal(resp, &searchRes)
	if err != nil {
		return
	}
	if searchRes.ErrorCode == 0 {
		userList := searchRes.Result.UserList
		if len(userList) > 0 && userList[0].UserInfo != "" {
			tel = userList[0].UserInfo
		}
	}
	return
}
