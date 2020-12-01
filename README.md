# temperature
测温系统 通过测温设备结合人脸识别系统，实现公司每日测温打卡。
此为设备Api。

## 系统
### 测温APP
包括测温App、测温Api

### 平台后台
前端、后台Api

### 企业后台
前端、后台Api

### 用户H5页
前端、后台Api

## 
人脸识别

## Flow
[flow](./spec/temperature_flow.png)
### 百度人脸注册
用户在百度H5页面，注册。上传人脸照片，及身份证、手机号等信息。
百度注册成功后，百度系统同步用户的faceID及注册信息到系统。
支持H5或者微信小程序注册
H5注册页 https://ai.baidu.com/facekit/page/form/E9F6EC9BDB8A4364BEA3DB864FE6B443

### 用户测量体温
用户在设备前，扫描人脸。设备获取到用户头像及体温数据，将数据通过此系统Api上报服务器。
服务器获取用户头像，通过调用百度人脸算法，获取用户人脸faceID。
服务器识别新老用户，返回相关注册信息。

### 平台后台
统一平台，独立企业后台，可以服务多个企业

### 企业后台
企业后台可以查看所属企业员工每天体温信息

## Api docs
spec/openapi.yml

## Server
### Config
使用Redis集群
阿里云oss对象存储，存储图片
使用百度人脸识别Api

### Docker build
```
# run
go run cmd/main.go

# test
make test

# build
make all

#push
make push
```
