package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"wcore/db"
	"wcore/x"

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
	engine.GET("/we_note/note/month/times", getNoteTimeListForMonth)
	engine.GET("/we_note/note/week/times", getNoteTimeListForWeek)
	engine.POST("/we_note/note", addNote)
	engine.PUT("/we_note/note", updateNote)
	engine.DELETE("/we_note/note", deleteNote)
	// 绑定端口，然后启动应用
	err := engine.Run(":8080")
	if err != nil {
		fmt.Printf("ListenAndServe err:%s", err.Error())
	}
}
func deleteNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := getNowTimeByMilli()
	if checkUser(req.UserID) {
		if db.DelNote(&db.WNote{
			WID:        req.UserID,
			DeleteTime: ts,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFail
		}
		go updateContinuedNum(req.UserID, ts, true)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

func updateNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := getNowTimeByMilli()
	if checkUser(req.UserID) {
		if db.UpdateNote(&db.WNote{
			WID:        req.UserID,
			WMood:      req.Mood,
			WDesc:      req.Desc,
			UpdateTime: ts,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFail
		}
		go updateContinuedNum(req.UserID, ts, false)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

func addNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := getNowTimeByMilli()
	if checkUser(req.UserID) {
		if db.CreateNote(&db.WNote{
			WID:        req.UserID,
			WMood:      req.Mood,
			WDesc:      req.Desc,
			CreateTime: ts,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFail
		}
		go updateContinuedNum(req.UserID, ts, false)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

func getNoteTimeListForWeek(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	onMonday := getWeekDay()
	weekEnd := onMonday + x.OneWeek
	if checkUser(req.UserID) {
		times := db.QueryNoteTimeByUserIDAndTimeRange(
			req.UserID,
			onMonday,
			weekEnd,
		)
		res.Data = times
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

func getNoteTimeListForMonth(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if checkUser(req.UserID) {
		times := db.QueryNoteTimeByUserIDAndTimeRange(
			req.UserID,
			req.CreateTime,
			req.CreateTime+x.OneMonth,
		)
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
		res.Data = db.QueryUserActionStatValByTypeAndUint(
			req.UserID,
			x.UserActStatTypeForContinueRecord,
			x.UserActStatUintForContinueRecord,
		)
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
	ts := getNowTimeByMilli()
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	}
	//loginMap.Store(loginRes.OpenID, 1)
	if wID, ok := checkUserIsLogin(loginRes.OpenID); ok {
		res.Data = wID
	} else {
		wID := generateUUID()
		if db.CreateUser(db.User{
			WID:        wID,
			ExtID:      loginRes.OpenID,
			WType:      1, // 微信用户
			CreateTime: ts,
		}) != nil {
			res.Msg = x.CreateUserFail
			res.Code = x.CreateUserFailErrCode
		} else {
			go updateContinuedNum(wID, ts, false)
		}
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

func getNowTimeByMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

func getWeekDay() int64 {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStart.UnixNano() / 1e6
}

func updateContinuedNum(userID string, ts int64, isLess bool) {
	stat := db.QueryUserActionStatValByTypeAndUint(
		userID,
		x.UserActStatTypeForContinueRecord,
		x.UserActStatUintForContinueRecord,
	)
	switch {
	case isLess == false && stat.WID == "":
		db.CreateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     int64(0),
			ActUnit:    x.UserActStatUintForContinueRecord,
			CreateTime: ts,
		})
	case isLess == false && stat.WID != "":
		db.UpdateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     stat.ActVal + int64(1),
			ActUnit:    x.UserActStatUintForContinueRecord,
			CreateTime: ts,
		})
	case isLess == true && stat.WID != "":
		db.UpdateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     int64(0),
			ActUnit:    x.UserActStatUintForContinueRecord,
			CreateTime: ts,
		})
	default:
		return
	}
}
