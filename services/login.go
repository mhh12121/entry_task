package service

import (
	"fmt"
	"log"

	dao "entry_task/DAO"
	data "entry_task/Data"
	Util "entry_task/Util"

	"github.com/golang/protobuf/proto"
)

// func LoginHandle(conn net.Conn, ruser data.User) {
// func LoginHandle(conn net.Conn, toServerD *data.ToServerData, wg *sync.WaitGroup) {

func LoginHandle(toServerD *data.ToServerData) (*data.ResponseFromServer, error) {
	// defer wg.Done()
	fmt.Println("login coming!")
	tmpdata := &data.User{}
	tmpErr := proto.Unmarshal(toServerD.Httpdata, tmpdata)
	if tmpErr != nil {
		fmt.Println("login err:", tmpErr)
		panic(tmpErr)
	}
	log.Println("tcp login username:" + tmpdata.GetUsername())
	// log.Println("login tcp decode data", tmpdata)
	//get remote Addr
	// remoteAddr := conn.RemoteAddr().String()
	// fmt.Println("tcp server connect:" + remoteAddr)
	//first go through redis cache
	//check if exists or different
	//what if login in another device?
	exists, errtoken := dao.CheckToken(*tmpdata.Username, *tmpdata.Token)
	if errtoken != nil {
		log.Println("login checktoken cache err", errtoken)
	}

	//todo
	//some problems here(consistency)
	//1.checktoken in redis success then return success msg to http
	//2.http redirect to home
	//3.in the same time, the token in redis expires
	if exists {
		//if exists just take info from cache
		// gob.Register(new(data.ResponseFromServer))

		returnValue := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: nil}

		return returnValue, nil
	}

	//check from mysql
	success, errorcheck := dao.Check(*tmpdata.Username, *tmpdata.Password)

	//login fail
	if !success || errorcheck != nil {
		log.Println("password wrong! usename:", tmpdata.GetUsername())

		returnValue := &data.ResponseFromServer{Success: proto.Bool(false), TcpData: nil}
		//--------gRPC no need to marshal----------------

		return returnValue, nil
	}

	//if mysql check success, it will save it to redis as cache or update cache
	tokenerr := dao.SetToken(*tmpdata.Username, *tmpdata.Token, Util.TokenExpires)
	if tokenerr != nil {
		log.Println("login save cache err", tokenerr)
	}
	//login success
	log.Println("login handle tcp")

	returnValue := &data.ResponseFromServer{Success: proto.Bool(true), TcpData: nil}
	//--------gRPC no need to marshal----------------

	return returnValue, nil
}
