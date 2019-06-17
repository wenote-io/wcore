package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"vn-light-core/db"
	"vn-light-core/x"

	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
	uuid "github.com/satori/go.uuid"
)

// NoteReq core module
type NoteReq struct {
	UserID     string `json:"user_id"`
	CreateTime int64  `json:"create_time"`
	Mood       int    `json:"mood"` // 1:狂喜 2: 开心 3:还行 4:不爽 5:超烂
	Desc       string `json:"desc"`
	Limit      int    `josn:"limit"`
	Offset     int    `json:"offset"`
}

// BaseRes res module
type BaseRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newBaseRes() *BaseRes {
	return &BaseRes{
		Code: x.SuccessCode,
		Msg:  x.SuccessMsg,
		Data: nil,
	}
}

func newNoteReqByGET(c *gin.Context) *NoteReq {
	req := &NoteReq{
		UserID:     c.Query("user_id"),
		CreateTime: c.GetInt64("create_time"),
		Mood:       c.GetInt("mood"),
		Desc:       c.Query("desc"),
	}
	req.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "0"))
	req.Offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	return req
}

func newNoteReqByPOST(c *gin.Context) *NoteReq {
	var note NoteReq
	_ = c.BindJSON(&note)
	return &note
}

var loginMap sync.Map
var tok string
var exp time.Duration

func main() {
	db.InitDB()
	// 初始化引擎
	engine := gin.Default()
	engine.GET("/we_note/user/login", login)
	engine.GET("/we_note/note/continued_num", getContinuedNum)
	engine.GET("/we_note/note/list", getNoteList)
	engine.GET("/we_note/note/times", getNoteTimeList)
	engine.GET("/we_note/note", addNote)
	// 绑定端口，然后启动应用
	err := engine.Run(":8080")
	if err != nil {
		fmt.Printf("ListenAndServe err:%s", err.Error())
	}
}

func addNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	var err error
	if checkUser(req.UserID) {
		err = db.CreateNote(&db.WNote{
			WMood:      req.Mood,
			WDesc:      req.Desc,
			CreateTime: time.Now().Unix(),
		})
	} else if err != nil {
		res.Code = x.CreateNoteFailErrCode
		res.Msg = x.CreateNoteFail
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

func getNoteTimeList(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if checkUser(req.UserID) {
		times := db.QueryNoteTimeByUserIDAndTimeRange(req.UserID, req.CreateTime, req.CreateTime+x.OneMonth)
		res.Data = times
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}
func getNoteList(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if checkUser(req.UserID) {
		notes := db.QueryNotesByWID(req.UserID, req.Limit, req.Offset)
		res.Data = notes
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}
func getContinuedNum(c *gin.Context) {
	res := newBaseRes()
	req := newNoteReqByGET(c)
	if checkUser(req.UserID) {
		// 查询用户连续登录天数
		res.Data = db.QueryUserActionStatValByTypeAndUint(req.UserID, x.UserActStatTypeForContinueRecord, x.UserActStatUintForContinueRecord)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

func login(c *gin.Context) {
	res := newBaseRes()
	code := c.Query("code")
	loginRes, err := weapp.Login(x.AppID, x.Ssecret, code)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	//loginMap.Store(loginRes.OpenID, 1)
	if wID, ok := checkUserIsLogin(loginRes.OpenID); ok {
		res.Data = wID
	} else {
		err = db.CreateUser(db.User{
			WID:        generateUUID(),
			ExtID:      loginRes.OpenID,
			WType:      1, // 微信用户
			CreateTime: time.Now().Unix(),
		})
	}
	if err != nil {
		res.Msg = x.CreateUserFail
		res.Code = x.CreateUserFailErrCode
	}
	c.JSON(200, res)
}
func generateUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func checkUser(userID string) bool {
	user := db.QueryUserByWID(userID)
	if user.WID == "" {
		return false
	}
	return true
}
func checkUserIsLogin(OpenID string) (string, bool) {
	//_, ok = loginMap.Load(userID)
	// 查库看看 user注册过没
	user := db.QueryUserByExtID(OpenID)
	// 若没用则添加用户
	if user.WID == "" {
		return "", false
	}
	return user.WID, true
}
