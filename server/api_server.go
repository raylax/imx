package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raylax/imx/core"
	"github.com/raylax/imx/handler"
	"github.com/raylax/imx/registry"
	"github.com/raylax/imx/router"
	"net/http"
)

type apiServer struct {
	addr          string
	registry      registry.Registry
	router        *gin.Engine
	messageRouter router.MessageRouter
}

func (a *apiServer) Serve() error {
	gin.SetMode(gin.ReleaseMode)
	a.router = gin.Default()
	group := a.router.Group("/group/:gid")
	group.GET("/join", a.handleGroupJoinUser)
	group.GET("/leave", a.handleGroupLeaveUser)
	handler.AddGroupHandler(&handler.DefaultGroupHandler{
		MessageRouter: a.messageRouter,
	})
	return a.router.Run(a.addr)
}

func (a *apiServer) Shutdown() {
	// noop
}

func (a *apiServer) handleGroupJoinUser(ctx *gin.Context) {
	gid := ctx.Param("gid")
	uid := ctx.Param("uid")
	group := core.Group{Id: gid}
	user := core.User{Id: uid}
	err := a.registry.RegGroup(group, user)
	result(ctx, err)
	go func() {
		for _, h := range handler.GetGroupHandlers() {
			h.HandleUserJoin(group, user)
		}
	}()
}

func (a *apiServer) handleGroupLeaveUser(ctx *gin.Context) {
	gid := ctx.Param("gid")
	uid := ctx.Param("uid")
	group := core.Group{Id: gid}
	user := core.User{Id: uid}
	a.registry.UnRegGroup(group, user)
	result(ctx, nil)
	go func() {
		for _, h := range handler.GetGroupHandlers() {
			h.HandleUserLeave(group, user)
		}
	}()
}

type resultCode int

const (
	codeOk  resultCode = 0
	codeErr resultCode = 1
)

type apiResult struct {
	Code    resultCode  `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func result(ctx *gin.Context, err error) {
	if err != nil {
		resultErr(ctx, err)
	} else {
		resultOk(ctx)
	}
}

func resultOk(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, apiResult{
		Code: codeOk,
	})
}

func resultErr(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, apiResult{
		Code:    codeErr,
		Message: err.Error(),
	})
}
