package service

import (
	"fmt"
	"log"
	"net"

	dao "entry_task/DAO"
	data "entry_task/Data"

	"github.com/golang/protobuf/proto"
)

func HomeHandle(conn net.Conn, username string, token interface{}) {
	//checktoken first
	exists, errtoken := dao.CheckToken(username, token.(string))
	// data.FailSafeCheckErr("home checktoken cache err", errtoken)
	//1. cookie still exists but token expires
	//---solution: clear cookie first then redirect to login
	//2. cookie expires but token exists
	//---solution: login and refresh the token

	//token not exists or not correct
	if !exists || errtoken != nil {
		log.Println("home checktoken cache err", errtoken)
		// gob.Register(new(data.ResponseFromServer))

		returnValue := &data.ResponseFromServer{Success: proto.Bool(false), TcpData: nil}
		returnValueData, rErr := proto.Marshal(returnValue)
		if rErr != nil {
			panic(rErr)
		}
		conn.Write(returnValueData)
		// encoder := gob.NewEncoder(conn)
		// errreturn := encoder.Encode(returnValue)
		// if errreturn != nil {
		// 	log.Println("home auth encode direct from cache err", errreturn)
		// }
		// data.FailSafeCheckErr("home auth encode direct from cache err", errreturn)
		return
	}

	//First go through the Redis get cache
	user, ok, err := dao.GetCacheInfo(username)
	if err != nil {
		log.Println("redis get cache fail err", err)
	}
	// data.FailSafeCheckErr("redis get cache fail:", err)
	//cache still valid
	//
	if ok && err == nil {
		log.Println("tcp home handle cache get info okay", *user)
		// gob.Register(new(data.RealUser))
		userData, userErr := proto.Marshal(user)
		if userErr != nil {
			panic(userErr)
		}
		tohttp := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: userData}
		tohttpData, toHttpErr := proto.Marshal(tohttp)
		if toHttpErr != nil {
			fmt.Println("tohttperr:", toHttpErr)
			panic(toHttpErr)
		}
		conn.Write(tohttpData)
		// encoder := gob.NewEncoder(conn)
		// errreturn := encoder.Encode(tohttp)
		// if errreturn != nil {
		// 	panic(errreturn)
		// }

		return
	}

	//cache expires or not exists then go to mysql
	userdb, okdb := dao.AllInfo(username)
	//retrieve from mysql success
	if okdb {
		//it will also save it to cache
		successCache := dao.SaveCacheInfo(username, *userdb.Nickname, *userdb.Avatar)
		if !successCache {
			fmt.Println("update redis homne cache fail")
			//do nothing
		}

		//save cache success
		//here how
		// if successCache {

		userdbData, userdbErr := proto.Marshal(userdb)
		if userdbErr != nil {
			panic(userdbErr)
		}
		tohttp := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: userdbData}
		tohttpData, toHttpErr := proto.Marshal(tohttp)
		if toHttpErr != nil {
			fmt.Println("tohttperr:", toHttpErr)
			panic(toHttpErr)
		}
		_, writeErr := conn.Write(tohttpData)
		if writeErr != nil {
			fmt.Println("home write conn err,", writeErr)
		}

		return
		// }

	}

	return
}
