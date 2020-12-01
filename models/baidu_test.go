package models_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"temperature/models"
	"temperature/test"
	"testing"
)

func TestBaiDu_GetToken(t *testing.T) {
	test.Prepare(t)
	assert := assert.New(t)

	baiDuClient := new(models.BaiDu)
	baiDuClient.Init()
	token := baiDuClient.Token
	assert.NotEmpty(token)
}

func TestBaiDu_Detect(t *testing.T) {
	test.Prepare(t)
	assert := assert.New(t)

	baiDuClient := new(models.BaiDu)
	baiDuClient.Init()
	imgUrl := "http://bj.bcebos.com/v1/aip-web/34547299248C4BCFA592E15DCB3BDD60?authorization=bce-auth-v1%2Ff86a2044998643b5abc89b59158bad6d%2F2020-03-05T08%3A15%3A31Z%2F-1%2F%2F880c44a1527aaa612422ed5016bf44bfb4631b3ae9c0cb43d1ae5b141dd468cf"
	respData, err := baiDuClient.Detect(imgUrl)

	assert.Nil(err)
	assert.NotEmpty(respData)
}


func TestBaiDu_Search(t *testing.T) {
	test.Prepare(t)
	assert := assert.New(t)

	baiDuClient := new(models.BaiDu)
	baiDuClient.Init()
	imgUrl := "https://image.baidu.com/search/detail?z=0&word=%E6%B8%B8%E4%BE%A0%E4%BD%9C%E5%93%81&hs=0&pn=0&spn=0&di=0&pi=5737220180451235231&tn=baiduimagedetail&is=0%2C0&ie=utf-8&oe=utf-8&cs=950652634%2C417824920&os=&simid=&adpicid=0&lpn=0&fm=&sme=&cg=&bdtype=-1&oriquery=&objurl=http%3A%2F%2Ft7.baidu.com%2Fit%2Fu%3D378254553%2C3884800361%26fm%3D79%26app%3D86%26f%3DJPEG%3Fw%3D1280%26h%3D2030&fromurl=&gsm=10000000001&catename=pcindexhot&islist=&querylist="
	respData, err := baiDuClient.Search(imgUrl)

	assert.Nil(err)
	fmt.Println(respData)
	assert.NotEmpty(respData)
}

func TestBaiDu_Add(t *testing.T) {
	test.Prepare(t)
	assert := assert.New(t)

	baiDuClient := new(models.BaiDu)
	baiDuClient.Init()
	imgUrl := "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1583754793816&di=2b5ba90bd29aef246a04f474544907af&imgtype=0&src=http%3A%2F%2Ft7.baidu.com%2Fit%2Fu%3D3616242789%2C1098670747%26fm%3D79%26app%3D86%26f%3DJPEG%3Fw%3D900%26h%3D1350"
	respData, err := baiDuClient.Add(imgUrl, "r02")

	assert.Nil(err)
	fmt.Println(respData)
	assert.NotEmpty(respData)
}
