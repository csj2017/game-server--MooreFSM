package main

import (
	"./teleport"
	"./teleport/debug"
	"./handlers"
	"log"
	//"time"
	"github.com/golang/protobuf/proto"
	"./proto"
	"bufio"
	"os"
	"strconv"
	"fmt"
)

var table_id string

// 有标识符UID的demo，保证了客户端链接唯一性
func main() {

	debug.Debug = true

	uid := "C1"

	//注册请求处理函数
	clientHandlers := teleport.API{
		"CreateRoomReturn" : new(ClientHeartBeat),
		"ActionResponse" : new(ActionResponse),
		"ActionPrompt" : new(ActionPrompt),
		//teleport.HEARTBEAT : new(ClientHeartBeat),
		teleport.IDENTITY : new(handlers.Identity),
	}

	//启动客户端
	tp := teleport.New().SetUID(uid, "abc").SetAPI( clientHandlers )
	tp.Client("127.0.0.1", ":20125")

	req := &server_proto.CreateRoomRequest{uid,4}
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatal("create room request error: ", err)
	}

	tp.Request(data, "CreateRoom", "create_room_flag")

	//指令
	running := true
	var inp []byte
	var data2 []byte
	var order string
	for running {
		fmt.Println("please input :")
		reader := bufio.NewReader(os.Stdin)
		inp, _, _ = reader.ReadLine()
		order = string(inp)
		if order == "stop"{
			running = false
		} else if order == "discard" {
			inp, _, _ = reader.ReadLine()
			card_input := string(inp)
			card, _ := strconv.Atoi(card_input)
			request := &server_proto.DiscardRequest{
				int32(card),
			}
			data2 = server_proto.MessageEncode( request )
			tp.Request(data2, "Discard", "discard_flag")
		} else if order == "action" {
			inp, _, _ = reader.ReadLine()
			select_id_input := string(inp)
			select_id, _ := strconv.Atoi(select_id_input)
			request := &server_proto.ActionSelectRequest{
				int32(select_id),
			}
			data2 = server_proto.MessageEncode( request )
			tp.Request(data2, "ActionSelect", "action_select")
		} else if order == "ready" {
			tp.Request(nil, "Ready", "ready_flag")
		}
	}

	select {}
}

type ClientHeartBeat struct{}
func (*ClientHeartBeat) Process(receive *teleport.NetData) *teleport.NetData {

	log.Println("=============room create return===============")

	// 进行解码
	response := &server_proto.CreateRoomResponse{}
	server_proto.MessageDecode( receive.Body, response )
	log.Println(response.RoomId)

	request := &server_proto.EnterRoomRequest{response.RoomId}
	data := server_proto.MessageEncode(request)

	return teleport.ReturnData( data, "EnterRoom" )
}

type ActionResponse struct{}
func (*ActionResponse) Process(receive *teleport.NetData) *teleport.NetData {

	log.Println("=============ActionResponse===============")

	// 进行解码
	response := &server_proto.ActionResponse{}
	server_proto.MessageDecode( receive.Body, response )
	log.Println(response.Uuid," ",response.ActionName," ",response.Card)
	return nil
}

type ActionPrompt struct {}
func (*ActionPrompt) Process(receive *teleport.NetData) *teleport.NetData {

	log.Println("=============ActionPrompt===============")

	// 进行解码
	response := &server_proto.ActionPrompt{}
	server_proto.MessageDecode( receive.Body, response )
	for _,v := range response.Action {
		log.Println( v )
	}
	return nil
}