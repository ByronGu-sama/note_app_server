package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"note_app_server/global"
	"note_app_server/model/appModel"
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
	// 绑定消息
	token := ctx.Query("t")
	uid, checkResult := verifyToken(token)
	if !checkResult {
		response.RespondWithStatusBadRequest(ctx, "校验错误")
		return
	}

	// 升级为ws连接
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 消息处理
	go func() {
		for msg := range messageQueue {
			RWLocker.Lock()
			err = msg.Conn.WriteMessage(websocket.TextMessage, msg.Msg.EncodeMessage())
			RWLocker.Unlock()
			if err != nil {
				log.Println(err)
			}
		}
	}()

	go connectionProc(conn, uid)
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
		for _, i := range strings {
			rawMsg := Msg{}
			if err = rawMsg.ParseMsg([]byte(i)); err == nil {
				messageQueue <- Message{Conn: conn, Msg: rawMsg}
			}
		}
		global.MsgRdb.LPop(context.TODO(), strconv.Itoa(int(uid)))
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			delConn(uid)
			break
		}
		tempMsg := &Msg{}
		if err = tempMsg.ParseMsg(message); err != nil {
			RWLocker.Lock()
			conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			RWLocker.Unlock()
			log.Println(err.Error())
		}
		tempMsg.FromId = uid
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
	var ctx = context.Background()
	if conn != nil {
		messageQueue <- Message{Conn: conn, Msg: *msg}
		msg.Read = true
	} else {
		msg.Read = false
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
	fromToTarget := fmt.Sprintf("%d-%d", msg.FromId, msg.ToId)
	targetToFrom := fmt.Sprintf("%d-%d", msg.ToId, msg.FromId)

	key := fromToTarget
	if fromToTarget > targetToFrom {
		key = targetToFrom
	}

	global.MsgRdb.LPush(ctx, key, msg.EncodeMessage())
	global.MsgRdb.Expire(ctx, key, 24*time.Hour)
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

// 删除连接
func delConn(fromId uint) {
	delete(websocketConn, fromId)
}

// 获取连接
func getConn(uid uint) *websocket.Conn {
	RWLocker.RLock()
	defer RWLocker.RUnlock()
	return websocketConn[uid]
}

// 检查token有效性
func verifyToken(token string) (uint, bool) {
	temp, err := ParseJWT(token)
	// 校验token有效性
	if token == "" || err != nil {
		return 0, false
	}
	claims := temp.(*appModel.JWT)
	// 验证uid
	uid := claims.Uid
	if uid == 0 {
		return 0, false
	}

	rCtx := context.Background()
	_, err = global.TokenRdb.Get(rCtx, strconv.Itoa(int(uid))).Result()
	if errors.Is(err, redis.Nil) {
		return 0, false
	} else if err != nil {
		return 0, false
	}

	var user *userModel.UserInfo
	if err = global.Db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return 0, false
	}

	return uid, true
}
