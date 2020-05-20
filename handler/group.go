package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/raylax/imx/core"
	"github.com/raylax/imx/router"
)
import pb "github.com/raylax/imx/proto"

var groupHandlers = make(map[string]GroupHandler)

type GroupHandler interface {
	Key() string
	// 用户进入群组
	HandleUserJoin(group core.Group, user core.User)
	// 用户离开群组
	HandleUserLeave(group core.Group, user core.User)
}

func AddGroupHandler(handler GroupHandler) {
	groupHandlers[handler.Key()] = handler
}

func RemoveGroupHandler(key string) {
	delete(groupHandlers, key)
}

func ResetGroupHandler() {
	groupHandlers = make(map[string]GroupHandler)
}

func GetGroupHandlers() map[string]GroupHandler {
	return groupHandlers
}

// 默认handler
type DefaultGroupHandler struct {
	MessageRouter router.MessageRouter
}

func (handler DefaultGroupHandler) Key() string {
	return "Default"
}

func (handler DefaultGroupHandler) HandleUserJoin(group core.Group, user core.User) {
	message := &pb.WsMessageRequest{
		MessageId: uuid.New().String(),
		TargetId:  group.Id,
		Type:      pb.MessageType_GROUP,
		Data:      fmt.Sprintf("[%s]加入群聊", user.Id),
		SourceId:  "0",
	}
	_ = handler.MessageRouter.RouteMessage(message)
}

func (handler DefaultGroupHandler) HandleUserLeave(group core.Group, user core.User) {
	message := &pb.WsMessageRequest{
		MessageId: uuid.New().String(),
		TargetId:  group.Id,
		Type:      pb.MessageType_GROUP,
		Data:      fmt.Sprintf("[%s]离开群聊", user.Id),
		SourceId:  "0",
	}
	_ = handler.MessageRouter.RouteMessage(message)
}
