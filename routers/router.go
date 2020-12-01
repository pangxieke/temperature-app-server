package routers

import (
	"github.com/julienschmidt/httprouter"
	"temperature/controllers"
)

func Router() *httprouter.Router {

	router := httprouter.New()
	//微信用户授权后重定向
	router.POST("/temperature/v1/data", controllers.Action((*controllers.Report).Create))
	//二维码

	//百度通知
	router.POST("/temperature/v1/notice", controllers.Action((*controllers.Notice).BaiDuNotice))

	//创建OSS授权
	router.POST("/temperature/v1/image", controllers.Action((*controllers.Image).Create))

	//短信验证码
	router.POST("/temperature/v1/send", controllers.Action((*controllers.User).SendSMS))
	//登录
	router.POST("/temperature/v1/login", controllers.Action((*controllers.User).Login))
	//用户信息
	router.GET("/temperature/v1/user", controllers.Action((*controllers.Record).UserInfo))
	//记录查询
	router.POST("/temperature/v1/list", controllers.Action((*controllers.Record).List))
	return router
}
