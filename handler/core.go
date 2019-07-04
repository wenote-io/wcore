package api

import (
	"fmt"
	"time"
	"wcore/module/dao"
	"wcore/module/view"
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

// 按天数获取notelist
func getNoteListByDay(c *gin.Context) {
	req := view.NewNoteReqByGET(c)
	res := x.NewBaseRes()
	if req.CreateTime == 0 {
		req.CreateTime = x.GetCurrDay()
	}
	fmt.Printf("getNoteListByDay s:%d,e:%d\n", req.CreateTime, req.CreateTime+x.OneDay)
	if dao.CheckUser(req.UserID) {
		notes := dao.QueryNotesByOneDay(req.UserID, req.CreateTime, req.CreateTime+x.OneDay, req.Limit, req.Offset)
		res.Data = view.NewNotesRes(notes)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

// getNoteTotalNum 获取记录总数
func getNoteTotalNum(c *gin.Context) {
	req := view.NewNoteReqByGET(c)
	res := x.NewBaseRes()
	if dao.CheckUser(req.UserID) {
		res.Data = dao.QueryTotalNoteNumByUserID(req.UserID)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  删除记录
func deleteNote(c *gin.Context) {
	req := view.NewNoteReqByPOST(c)
	res := x.NewBaseRes()
	ts := x.GetNowTimeByMilli()
	if dao.CheckUser(req.UserID) {
		if dao.DelNote(&dao.WNote{
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
	req := view.NewNoteReqByPOST(c)
	res := x.NewBaseRes()
	ts := x.GetNowTimeByMilli()
	if dao.CheckUser(req.UserID) {
		if dao.UpdateNote(&dao.WNote{
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
	req := view.NewNoteReqByPOST(c)
	res := x.NewBaseRes()
	//ts := getNowTimeByMilli()
	nID := x.GenerateUUID()
	if dao.CheckUser(req.UserID) {
		if req.CreateTime == 0 {
			req.CreateTime = x.GetNowTimeByMilli()
		}
		if dao.CreateNote(&dao.WNote{
			WID:        req.UserID,
			NID:        nID,
			WMood:      req.Mood,
			WDesc:      req.Desc,
			CreateTime: req.CreateTime,
		}) != nil {
			res.Code = x.CreateNoteFailErrCode
			res.Msg = x.CreateNoteFailMsg
		}
		go dao.UpdateContinuedNum(req.UserID, req.CreateTime, true)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取本周的note 创建时间的列表
func getNoteTimeListForWeek(c *gin.Context) {
	req := view.NewNoteReqByGET(c)
	res := x.NewBaseRes()
	onMonday := x.GetWeekDay()
	weekEnd := onMonday + x.OneWeek
	if dao.CheckUser(req.UserID) {
		times := dao.QueryNoteTimeByUserIDAndTimeRange(
			req.UserID,
			onMonday,
			weekEnd,
		)
		weeks := make([]int, 0)
		for _, v := range times {
			t := time.Unix(v/1e3, 0)                //将其转换为秒,转换为日期
			weeks = append(weeks, int(t.Weekday())) // 获取本周星期数
		}
		res.Data = x.RemoveRepByMap(weeks)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取本月的note 创建时间的列表
func getNoteTimeListForMonth(c *gin.Context) {
	req := view.NewNoteReqByGET(c)
	res := x.NewBaseRes()
	if dao.CheckUser(req.UserID) {
		times := dao.QueryNoteTimeByUserIDAndTimeRange(
			req.UserID,
			req.CreateTime,
			req.CreateTime+x.OneMonth,
		)
		days := make([]int, 0)
		for _, v := range times {
			t := time.Unix(v/1e3, 0)     //将其转换为秒,转换为日期
			days = append(days, t.Day()) // 获取时间当月的天数
		}
		res.Data = x.RemoveRepByMap(days)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取记录列表分页
func getNoteList(c *gin.Context) {
	req := view.NewNoteReqByGET(c)
	res := x.NewBaseRes()
	if dao.CheckUser(req.UserID) {
		notes := dao.QueryNotesByWID(req.UserID, req.Limit, req.Offset)
		res.Data = view.NewNotesRes(notes)
		go dao.UpdateContinuedNum(req.UserID, x.GetCurrDay(), false)
	} else {
		res.Msg = x.UserNotFoundErrMsg
		res.Code = x.UserNotFoundErrCode
	}
	c.JSON(200, res)
}

//  获取连续记录数
func getContinuedNum(c *gin.Context) {
	res := x.NewBaseRes()
	req := view.NewNoteReqByGET(c)
	if dao.CheckUser(req.UserID) {
		// 查询用户连续登录天数
		stat := dao.QueryUserActionStatValByTypeAndUint(
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
	res := x.NewBaseRes()
	code := c.Query("code")
	if code == "" {
		res.Msg = x.CodeInvalidMsg
		res.Code = x.CodeInvalidErrCode
		c.JSON(200, res)
		return
	}
	ts := x.GetNowTimeByMilli()
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
	if wID, ok = dao.CheckUserIsLogin(loginRes.OpenID); ok {
		res.Data = wID
	} else {
		wID = x.GenerateUUID()
		if dao.CreateUser(dao.User{
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
