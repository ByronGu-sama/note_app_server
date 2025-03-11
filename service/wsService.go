package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"note_app_server/global"
	"note_app_server/model/appModel"
	"note_app_server/model/msgModel"
	"note_app_server/model/userModel"
	"note_app_server/producer"
	"note_app_server/response"
	"strconv"
	"sync"
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
	ctx := context.TODO()
	defer conn.Close()
	// 已存在连接
	if getConn(uid) != nil {
		return
	}

	addConn(uid, conn)

	// 处理待发送消息
	go rePushMsg(uid, conn, ctx)

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
	// 存储id格式遵循小id在前，大id在后
	var firstKey uint
	var secondKey uint
	if msg.FromId < msg.ToId {
		firstKey = msg.FromId
		secondKey = msg.ToId
	} else {
		firstKey = msg.ToId
		secondKey = msg.FromId
	}

	// 存入mongoDB持久化存储
	err := producer.SyncMessageToMongo(firstKey, secondKey, (*msgModel.Message)(msg))
	if err != nil {
		log.Println(err)
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

// 重新推送用户不在线时收到的数据
func rePushMsg(uid uint, conn *websocket.Conn, ctx context.Context) {
	mongoConn := global.MongoClient.Database("pending_message").Collection("msgs")

	filter := bson.D{
		{
			Key: "$or",
			Value: bson.A{
				bson.D{
					{Key: "uid1", Value: uid},
				},
				bson.D{
					{Key: "uid2", Value: uid},
				},
			},
		},
	}

	option := options.Find().SetSort(bson.D{{Key: "pubTime", Value: 1}})

	cursor, err := mongoConn.Find(ctx, filter, option)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var results []Msg
	if err = cursor.All(ctx, &results); err != nil {
		log.Println(err)
	}

	for _, i := range results {
		messageQueue <- Message{Conn: conn, Msg: i}
	}
}
