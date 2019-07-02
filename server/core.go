package server

import (
	"wcore/db"
	"wcore/x"
)

// CheckUser 校验用户是否存在
func CheckUser(userID string) bool {
	user := db.QueryUserByWID(userID)
	if user.WID == "" {
		return false
	}
	return true
}

// CheckUserIsLogin 校验用户是否登陆
func CheckUserIsLogin(OpenID string) (string, bool) {
	//_, ok = loginMap.Load(userID)
	// 查库看看 user注册过没
	user := db.QueryUserByExtID(OpenID)
	// 若没用则添加用户
	if user.WID == "" {
		return "", false
	}
	return user.WID, true
}

// UpdateContinuedNum 更新连续记录值
func UpdateContinuedNum(userID string, ts int64) {
	stat := db.QueryUserActionStatValByTypeAndUint(
		userID,
		x.UserActStatTypeForContinueRecord,
		x.UserActStatUintForContinueRecord,
	)
	if stat.WID == "" {
		// 若记录不存在,则初始化一个连续记录天数
		db.CreateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     int64(1),
			ActUnit:    x.UserActStatUintForContinueRecord,
			CreateTime: ts,
		})
		return
	}
	dt := x.GetCurrDay() - x.GetByZeroMorningTs(stat.UpdateTime)
	switch {
	case dt == x.OneDay: // 若记录数在上次记录后的一天之内则增加记录值
		db.UpdateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     stat.ActVal + int64(1),
			ActUnit:    x.UserActStatUintForContinueRecord,
			UpdateTime: ts,
		})
	case dt > x.OneDay:
		db.UpdateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     int64(1),
			ActUnit:    x.UserActStatUintForContinueRecord,
			UpdateTime: ts,
		})
	default:
	}
}
