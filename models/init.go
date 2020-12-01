package models

import (
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"strings"
	"temperature/config"
)

var (
	RedisCluster *redis.ClusterClient
	db           *gorm.DB
	secret       []byte
)

func Init() (err error) {
	err = initMySQL()
	if err != nil {
		return
	}

	if err = initRedis(); err != nil {
		return
	}

	err = initOSS(config.OSS)
	if err != nil {
		return
	}

	//err = initStages("stages.json")
	//if err != nil {
	//	return
	//}
	//err = initWigs("wigs.json")
	//if err != nil {
	//	return
	//}
	secret, _ = hex.DecodeString("temperature")
	return
}

func InitForTest() (err error) {
	err = initRedis()
	if err != nil {
		return err
	}
	err = initOSS(config.OSS)
	if err != nil {
		return err
	}
	if err = initMySQL(); err != nil {
		return err
	}
	//
	//err = initStages("../stages.json")
	//if err != nil {
	//	return
	//}
	//err = initWigs("../wigs.json")
	//if err != nil {
	//	return
	//}
	//
	//db = d
	//if err := Migrate(); err != nil {
	//	return errors.Wrapf(err, "Failed to migrate database")
	//}
	//db.LogMode(true)
	return
}

func initMySQL() (err error) {
	host := config.MySQL.Host
	port := config.MySQL.Port
	username := config.MySQL.User
	password := config.MySQL.Password
	name := config.MySQL.DB
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, name)
	db, err = gorm.Open("mysql", args)
	if err != nil {
		return errors.Wrapf(err, "Failed to open database")
	}

	db.LogMode(config.MySQL.Debug)
	return
}

func initRedis() (err error) {
	opt := redis.ClusterOptions{
		Addrs:    strings.Split(config.Redis.Address, ","),
		Password: "",
		PoolSize: 10,
	}
	RedisCluster = redis.NewClusterClient(&opt)

	//RedisClient = redis.NewClient(&redis.Options{
	//	Addr:     config.Redis.Address,
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})

	pong, err := RedisCluster.Ping().Result()
	if err != nil {
		return errors.Wrapf(err, "Failed to connect redis")
	}
	if pong != "PONG" {
		return errors.Wrapf(err, "Failed to ping redis")
	}
	return
}
