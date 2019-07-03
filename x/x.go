package x

import (
	"strings"
	"wcore_old/x"

	uuid "github.com/satori/go.uuid"
)

// BaseRes res module
type BaseRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewBaseRes 通用响应模版
func NewBaseRes() *BaseRes {
	return &BaseRes{
		Code: x.SuccessCode,
		Msg:  x.SuccessMsg,
		Data: nil,
	}
}

// RemoveRepByMap int 去除重复
func RemoveRepByMap(slc []int) []int {
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

// GenerateUUID 生成唯一id
func GenerateUUID() string {
	u := uuid.Must(uuid.NewV4()).String()
	return strings.Replace(u, "-", "", -1)
}
