package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"temperature/config"
	"time"
)

type SMS struct {
	Mobile string
}

func (this *SMS) Send() (err error) {
	if this.Mobile == "" {
		err = fmt.Errorf("mobile is empty")
	}

	expire := 1
	verify_code := this.randInt(6)
	fmt.Printf("sms verify_code: %s, mobile: %s", verify_code, this.Mobile)
	msg := fmt.Sprintf("您的验证码为%s,有效时间%d分钟", verify_code, expire)

	err = this.SendToSMS(msg, this.Mobile)
	if err != nil {
		log.Printf("SendToSMS err=%s", err)
		err = fmt.Errorf("验证码发送失败")
		return
	}

	err = RedisCluster.Set(this.Mobile, verify_code, 2*time.Minute).Err()
	if err != nil {
		log.Println("redis set err ", err)
		return
	}

	return
}

func (this *SMS) Verify(code string) (bool) {
	verifyCode, err := RedisCluster.Get(this.Mobile).Result()
	if err != nil {
		log.Printf("RedisCluster.Get err=%s", err)
	}
	if verifyCode != "" && verifyCode == code {
		return true
	}

	return false

}

func (this *SMS) SendToSMS(msg, mobile string) (err error) {
	params := make(map[string]interface{})

	fmt.Println(msg)
	//请登录zz.253.com获取API账号、密码以及短信发送的URL
	params["account"] = config.SMS.Account    //创蓝API账号
	params["password"] = config.SMS.Password  //创蓝API密码
	params["phone"] = mobile              //手机号码

	//设置您要发送的内容：其中“【】”中括号为运营商签名符号，多签名内容前置添加提交
	params["msg"] = url.QueryEscape(msg)
	params["report"] = "true"
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	host := "http://smssh1.253.com/msg/send/json" //短信发送URL
	request, err := http.NewRequest("POST", host, reader)
	if err != nil {
		log.Printf("sms http.NewRequest err=%s", err)
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("sms send client err=%s", err)
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("sms ioutil.ReadAll err=%s", err)
		return
	}
	type smsRespStr struct {
		Code     string `json:"code"`
		MsgID    string `json:"msgId"`
		Time     string `json:"time"`
		ErrorMsg string `json:"errorMsg"`
	}

	var smsResp smsRespStr
	_ = json.Unmarshal(respBytes, &smsResp)
	if smsResp.Code != "0" {
		//失败
		log.Printf("sms send err,resp=%s", respBytes)
		err = fmt.Errorf("发送失败")
		return
	}

	return
}

func (this *SMS) randInt(size int) string {
	result := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		result[i] = uint8(48 + rand.Intn(10))
	}
	return string(result)
}
