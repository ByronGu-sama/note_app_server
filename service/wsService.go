package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"note_app_server/global"
	"note_app_server/model/msgModel"
	"note_app_server/model/userModel"
	"note_app_server/response"
	"strconv"
	"sync"
	"time"
)

type Msg msgModel.Message
type Message struct {
	Conn *websocket.Conn
	Msg  Msg
}

// ParseMsg 解码消息
func (msg *Msg) ParseMsg(message []byte) error {
	fmt.Println(message)
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return err
	}
	return nil
}

// EncodeMessage 编码消息
func (msg *Msg) EncodeMessage() []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	return data
}

var (
	users      = make(map[uint]userModel.UserInfo)
	privateMsg = make([]string, 0)
	publicMsg  = make([]string, 0)
	upGrader   = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	websocketConn = make(map[uint]*websocket.Conn)
	messageQueue  = make(chan Message, 200)
	RWLocker      = new(sync.RWMutex)
)

func InitWS(ctx *gin.Context) {
	// 消息处理
	go func() {
		for msg := range messageQueue {
			if msg.Conn != nil {
				err := msg.Conn.WriteMessage(websocket.TextMessage, msg.Msg.EncodeMessage())
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	// 绑定消息
	tempUid, ok := ctx.Get("uid")
	if !ok {
		response.RespondWithStatusBadRequest(ctx, "缺少必要信息")
		return
	}
	fromId := tempUid.(uint)

	// 升级为ws连接
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	go connectionProc(conn, fromId)
}

// 处理连接
func connectionProc(conn *websocket.Conn, uid uint) {
	defer conn.Close()
	// 已存在连接
	if getConn(uid) != nil {
		return
	}

	addConn(uid, conn)

	// 处理待发送消息
	go func() {
		result := global.MsgRdb.LRange(context.TODO(), strconv.Itoa(int(uid)), 0, -1)
		strings, err := result.Result()
		if err != nil {
			return
		}
		for i := range strings {
			rawMsg := Msg{}
			err = rawMsg.ParseMsg([]byte(strings[i]))
			if err != nil {
				continue
			}
			messageQueue <- Message{Conn: conn, Msg: rawMsg}
		}
		global.MsgRdb.LPop(context.TODO(), strconv.Itoa(int(uid)))
	}()

	for {
		_, message, err := conn.ReadMessage()
		tempMsg := &Msg{}
		err = tempMsg.ParseMsg(message)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			log.Println(err.Error())
		}

		if tempMsg.Type == 1 {
			privateMsgProc(tempMsg)
		} else if tempMsg.Type == 2 {
			groupMsgProc(tempMsg)
		}
	}
}

// 私聊消息
func privateMsgProc(msg *Msg) {
	conn := getConn(msg.ToId)
	if conn != nil {
		messageQueue <- Message{Conn: conn, Msg: *msg}
	} else {
		var ctx = context.Background()
		// 接收人不在线暂存至redis列表保存30min，超时归档至mongodb
		global.MsgRdb.LPush(ctx, strconv.Itoa(int(msg.ToId)), msg.EncodeMessage())
		global.MsgRdb.Expire(ctx, strconv.Itoa(int(msg.ToId)), 30*time.Minute)
		// 消息存储优化点
		// 1、redis中的消息设置过期时间
		// 2、消息队列加入确认机制
		// 3、批处理推送数据到mq
		// 4、消息重试机制
		// 延伸操作：
		// 归档到冷存储
		// 数据可视化
		// 可以加入到数据分析系统
	}
}

// 群发消息
func groupMsgProc(msg *Msg) {
	//defer RWLocker.RUnlock()
	//RWLocker.RLock()
	//for uid, conn := range websocketConn {
	//
	//}
}

// 添加连接
func addConn(fromId uint, conn *websocket.Conn) {
	RWLocker.RLock()
	defer RWLocker.RUnlock()
	websocketConn[fromId] = conn
}

// 获取连接
func getConn(uid uint) *websocket.Conn {
	RWLocker.RLock()
	defer RWLocker.RUnlock()
	return websocketConn[uid]
}
