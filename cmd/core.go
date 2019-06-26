package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

var stat sync.Map
var tok string
var exp time.Duration

func main() {
	db.InitDB()
	// 初始化引擎
	engine := gin.Default()
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
	// 绑定端口，然后启动应用
	err := engine.Run(":8080")
	if err != nil {
		fmt.Printf("ListenAndServe err:%s", err.Error())
	}
}

func getNoteListByDay(c *gin.Context) {
	req := newNoteReqByGET(c)
	res := newBaseRes()
	if req.CreateTime == 0 {
		req.CreateTime = getCurrDay(time.Now())
	} else {
		req.CreateTime = getCurrDay(time.Unix(req.CreateTime, 0))
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

func deleteNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := getNowTimeByMilli()
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

func updateNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	ts := getNowTimeByMilli()
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

func addNote(c *gin.Context) {
	req := newNoteReqByPOST(c)
	res := newBaseRes()
	//ts := getNowTimeByMilli()
	nID := generateUUID()
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
func generateUUID() string {
	u := uuid.Must(uuid.NewV4()).String()
	return strings.Replace(u, "-", "", -1)
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

func getCurrDay(t time.Time) int64 {
	now := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return now.UnixNano() / 1e6
}

func updateContinuedNum(userID string, ts int64) {
	stat := db.QueryUserActionStatValByTypeAndUint(
		userID,
		x.UserActStatTypeForContinueRecord,
		x.UserActStatUintForContinueRecord,
	)
	if stat.WID == "" {
		db.CreateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     int64(0),
			ActUnit:    x.UserActStatUintForContinueRecord,
			CreateTime: ts,
		})
	} else {
		if stat.UpdateTime+x.OneDay > ts {
			return
		}
		db.UpdateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     stat.ActVal + int64(1),
			ActUnit:    x.UserActStatUintForContinueRecord,
			UpdateTime: ts,
		})
	}
}

func removeRepByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]struct{}{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = struct{}{}
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
