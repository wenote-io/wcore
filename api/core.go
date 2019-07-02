package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"wcore/db"
	"wcore/x"

	"github.com/gin-gonic/gin"
	"github.com/medivhzhan/weapp"
)

// InitHandle 初始化handler
func InitHandle(engine *gin.Engine) {
	engine.GET("/we_note/user/login", login)
	engine.GET("/we_note/note/continued_num", getContinuedNum)
	engine.GET("/we_note/note/list", getNoteList)
	engine.GET("/we_note/note/day/list", getNoteListByDay)
	engine.GET("/we_note/note/month/times", getNoteTimeListForMonth)
	engine.GET("/we_note/note/week/times", getNoteTimeListForWeek)
	engine.GET("/we_note/note/total", getNoteTotalNum)
	engine.POST("/we_note/note", addNote)
	engine.PUT("/we_note/note", updateNote)
	engine.DELETE("/we_note/note", deleteNote)
}

// NoteReq core module
type NoteReq struct {
	UserID     string `json:"user_id"`
	NoteID     string `json:"note_id"`
	CreateTime int64  `json:"create_time"`
	Mood       int    `json:"mood"` // 1:狂喜 2: 开心 3:还行 4:不爽 5:超烂
	Desc       string `json:"desc"`
	Limit      int    `josn:"limit"`
	Offset     int    `json:"offset"`
}

// NoteRes core module
type NoteRes struct {
	UserID     string `json:"user_id"`
	NoteID     string `json:"note_id"`
	CreateTime int64  `json:"create_time"`
	Mood       int    `json:"mood"` // 1:狂喜 2: 开心 3:还行 4:不爽 5:超烂
	Desc       string `json:"desc"`
}

// BaseRes res module
type BaseRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newNotesRes(notes []*db.WNote) []*NoteRes {
	res := make([]*NoteRes, 0)
	for _, v := range notes {
		res = append(res, newNoteRes(v))
	}
	return res
}

func newNoteRes(note *db.WNote) *NoteRes {
	return &NoteRes{
		UserID:     note.WID,
		NoteID:     note.NID,
		Mood:       note.WMood,
		Desc:       note.WDesc,
		CreateTime: note.CreateTime,
	}
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
		UserID: c.Query("user_id"),
		Desc:   c.Query("desc"),
	}
	req.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "0"))
	req.Offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	req.Mood, _ = strconv.Atoi(c.DefaultQuery("mood", "0"))
	req.CreateTime, _ = strconv.ParseInt(c.DefaultQuery("create_time", "0"), 10, 64)
	data, _ := json.Marshal(req)
	fmt.Printf("get param:[%s]\n", string(data))
	return req
}

func newNoteReqByPOST(c *gin.Context) *NoteReq {
	var note NoteReq
	_ = c.BindJSON(&note)
	data, _ := json.Marshal(note)
	fmt.Printf("post param:[%s]\n", string(data))
	return &note
}

// 按天数获取notelist
func getNoteListByDay(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if req.CreateTime == 0 {
		req.CreateTime = x.GetCurrDay()
	}
	fmt.Printf("getNoteListByDay s:%d,e:%d\n", req.CreateTime, req.CreateTime+x.OneDay)
	if checkUser(req.UserID) {
		notes := db.QueryNotesByOneDay(req.UserID, req.CreateTime, req.CreateTime+x.OneDay, req.Limit, req.Offset)
		res.Data = newNotesRes(notes)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

// getNoteTotalNum 获取记录总数
func getNoteTotalNum(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if checkUser(req.UserID) {
		res.Data = db.QueryTotalNoteNumByUserID(req.UserID)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  删除记录
func deleteNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := x.GetNowTimeByMilli()
	if checkUser(req.UserID) {
		if db.DelNote(&db.WNote{
			//WID:        req.UserID,
			NID:        req.NoteID,
			DeleteTime: ts,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFailMsg
		}
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  更新记录
func updateNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := x.GetNowTimeByMilli()
	if checkUser(req.UserID) {
		if db.UpdateNote(&db.WNote{
			WID:        req.UserID,
			NID:        req.NoteID,
			WMood:      req.Mood,
			WDesc:      req.Desc,
			UpdateTime: ts,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFailMsg
		}
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  添加记录
func addNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	//ts := getNowTimeByMilli()
	nID := x.GenerateUUID()
	if checkUser(req.UserID) {
		if req.CreateTime == 0 {
			req.CreateTime = getNowTimeByMilli()
		}
		if db.CreateNote(&db.WNote{
			WID:        req.UserID,
			NID:        nID,
			WMood:      req.Mood,
			WDesc:      req.Desc,
			CreateTime: req.CreateTime,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFailMsg
		}
		go updateContinuedNum(req.UserID, req.CreateTime)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取本周的note 创建时间的列表
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
		weeks := make([]int, 0)
		for _, v := range times {
			t := time.Unix(v/1e3, 0)                //将其转换为秒,转换为日期
			weeks = append(weeks, int(t.Weekday())) // 获取本周星期数
		}
		res.Data = removeRepByMap(weeks)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取本月的note 创建时间的列表
func getNoteTimeListForMonth(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if checkUser(req.UserID) {
		times := db.QueryNoteTimeByUserIDAndTimeRange(
			req.UserID,
			req.CreateTime,
			req.CreateTime+x.OneMonth,
		)
		days := make([]int, 0)
		for _, v := range times {
			t := time.Unix(v/1e3, 0)     //将其转换为秒,转换为日期
			days = append(days, t.Day()) // 获取时间当月的天数
		}
		res.Data = removeRepByMap(days)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取记录列表分页
func getNoteList(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if checkUser(req.UserID) {
		notes := db.QueryNotesByWID(req.UserID, req.Limit, req.Offset)
		res.Data = newNotesRes(notes)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取连续记录数
func getContinuedNum(c *gin.Context) {
	res := newBaseRes()
	req := newNoteReqByGET(c)
	if checkUser(req.UserID) {
		// 查询用户连续登录天数
		stat := db.QueryUserActionStatValByTypeAndUint(
			req.UserID,
			x.UserActStatTypeForContinueRecord,
			x.UserActStatUintForContinueRecord,
		)
		res.Data = stat.ActVal
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  用户登陆
func login(c *gin.Context) {
	res := newBaseRes()
	code := c.Query("code")
	if code == "" {
		res.Msg = x.CodeInvalidMsg
		res.Code = x.CodeInvalidErrCode
		c.JSON(200, res)
		return
	}
	ts := getNowTimeByMilli()
	//loginRes := weapp.LoginResponse{OpenID: "123"}
	loginRes, err := weapp.Login(x.AppID, x.Ssecret, code)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
		res.Msg = x.CodeInvalidMsg
		res.Code = x.CodeInvalidErrCode
		c.JSON(200, res)
		return
	}
	var wID string
	var ok bool
	//loginMap.Store(loginRes.OpenID, 1)
	if wID, ok = checkUserIsLogin(loginRes.OpenID); ok {
		res.Data = wID
	} else {
		wID = generateUUID()
		if db.CreateUser(db.User{
			WID:        wID,
			ExtID:      loginRes.OpenID,
			WType:      1, // 微信用户
			CreateTime: ts,
		}) != nil {
			res.Msg = x.CreateUserFailMsg
			res.Code = x.CreateUserFailErrCode
		} else {
			res.Data = wID
		}
	}
	c.JSON(200, res)
}
