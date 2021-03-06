package dao

import (
	"fmt"
	"log"
	"time"

	"entry_task/Conf"
	data "entry_task/Data"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

var client *redis.Client

func init() {
	// _, filepath, _, _ := runtime.Caller(0)
	// p := path.Dir(filepath)
	// p = path.Dir(p)

	// log.Println("log path", p)

	// Conf.LoadConf(p + "/Conf/config.json")
	// log.Println(Conf.Config.Redis.Host+":"+Conf.Config.Redis.Port, Conf.Config.Redis.Password, Conf.Config.Redis.Db, Conf.Config.Redis.Poolsize)
	client = redis.NewClient(&redis.Options{
		Addr:     Conf.Config.Redis.Host + ":" + Conf.Config.Redis.Port,
		Password: Conf.Config.Redis.Password,
		DB:       Conf.Config.Redis.Db,
		PoolSize: Conf.Config.Redis.Poolsize,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		//do nothing
		// panic(err)
	}

	fmt.Println("initialize redis:", pong)
}

// InvaildCache
func InvalidCache(username string, token string) error {
	//todo
	// dont know the return value
	log.Println("invalid redis")
	errinfo := client.HSet(username, "valid", "0").Err()
	if errinfo != nil {

		return errinfo
	}
	client.Del(tokenFormat(username))
	return nil
}

// SetToken
func SetToken(username string, token string, expiration int64) error {
	log.Println(" setTOken username, token:", username, token)
	err := client.Set(tokenFormat(username), token, time.Duration(expiration)).Err()
	if err != nil {
		return err
	}
	return nil
}

// CheckToken
func CheckToken(username string, token string) (bool, error) {
	log.Println("checktoken redis:", username, token)
	val, err := client.Get(tokenFormat(username)).Result()
	if err != nil {
		return false, err
	}

	return token == val, nil
}

//not used
func SaveCacheInfo(username string, nickname string, avatar string) bool {
	tmp := map[string]interface{}{
		"valid":    "1",
		"nickname": nickname,
		"avatar":   avatar,
	}

	err := client.HMSet(username, tmp).Err()
	if err != nil {
		fmt.Println("redis save cache fail:", err)
		return false
	}
	// client.Save()
	return true
}

// CacheInfo
//not used
func GetCacheInfo(username string) (*data.RealUser, bool, error) {
	val, err := client.HGetAll(username).Result()
	log.Println("val", val)
	log.Println("redis val", val)
	if err != nil {
		return nil, false, err
	}
	if val["valid"] != "0" && len(val) != 0 {
		tmpuser := &data.RealUser{Username: proto.String(username), Avatar: proto.String(val["avatar"]), Nickname: proto.String(val["nickname"])}
		return tmpuser, true, nil
	}
	return nil, false, err
}

//not used
func UpdateCacheNickname(username string, nickname string) error {
	row := map[string]interface{}{
		"valid":    "1",
		"nickname": nickname,
	}
	err := client.HMSet(username, row).Err()
	if err != nil {
		return err
	}
	return nil
}

// update avatar
//not used
func UpdateCacheAvatar(username string, avatar string) error {
	row := map[string]interface{}{
		"valid":  "1",
		"avatar": avatar,
	}
	err := client.HMSet(username, row).Err()
	if err != nil {
		return err
	}
	return nil
}

func tokenFormat(username string) string {
	return "token_" + username
}
