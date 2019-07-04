package view

import (
	"encoding/json"
	"fmt"
	"strconv"
	"wcore/module/dao"

	"github.com/gin-gonic/gin"
)

// NoteReq 请求结构体
type NoteReq struct {
	UserID     string `json:"user_id"`
	NoteID     string `json:"note_id"`
	CreateTime int64  `json:"create_time"`
	Mood       int    `json:"mood"` // 1:狂喜 2: 开心 3:还行 4:不爽 5:超烂
	Desc       string `json:"desc"`
	Limit      int    `josn:"limit"`
	Offset     int    `json:"offset"`
}

// NoteRes 返回结构体
type NoteRes struct {
	UserID     string `json:"user_id"`
	NoteID     string `json:"note_id"`
	CreateTime int64  `json:"create_time"`
	Mood       int    `json:"mood"` // 1:狂喜 2: 开心 3:还行 4:不爽 5:超烂
	Desc       string `json:"desc"`
}

// NewNoteReqByGET x
func NewNoteReqByGET(c *gin.Context) *NoteReq {
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

// NewNoteReqByPOST xx
func NewNoteReqByPOST(c *gin.Context) *NoteReq {
	var note NoteReq
	_ = c.BindJSON(&note)
	data, _ := json.Marshal(note)
	fmt.Printf("post param:[%s]\n", string(data))
	return &note
}

// NewNotesRes 创建结果集
func NewNotesRes(notes []*dao.WNote) []*NoteRes {
	res := make([]*NoteRes, 0)
	for _, v := range notes {
		res = append(res, NewNoteRes(v))
	}
	return res
}

// NewNoteRes 创建结果
func NewNoteRes(note *dao.WNote) *NoteRes {
	return &NoteRes{
		UserID:     note.WID,
		NoteID:     note.NID,
		Mood:       note.WMood,
		Desc:       note.WDesc,
		CreateTime: note.CreateTime,
	}
}
