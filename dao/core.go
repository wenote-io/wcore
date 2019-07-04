package dao

import (
	"fmt"
	"log"
	"wcore/db"
	"wcore/x"
)

// CheckUser 校验用户是否存在
func CheckUser(userID string) bool {
	user := QueryUserByWID(userID)
	if user.WID == "" {
		return false
	}
	return true
}

// CheckUserIsLogin 校验用户是否登陆
func CheckUserIsLogin(OpenID string) (string, bool) {
	//_, ok = loginMap.Load(userID)
	// 查库看看 user注册过没
	user := QueryUserByExtID(OpenID)
	// 若没用则添加用户
	if user.WID == "" {
		return "", false
	}
	return user.WID, true
}

// UpdateContinuedNum 更新连续记录值
func UpdateContinuedNum(userID string, ts int64, isADD bool) {
	stat := QueryUserActionStatValByTypeAndUint(
		userID,
		x.UserActStatTypeForContinueRecord,
		x.UserActStatUintForContinueRecord,
	)
	if stat.WID == "" {
		// 不是添加操作
		if !isADD {
			return
		}
		// 若记录不存在,则初始化一个连续记录天数
		CreateUserActionStat(&db.UserActionStat{
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
		UpdateUserActionStat(&db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActVal:     stat.ActVal + int64(1),
			ActUnit:    x.UserActStatUintForContinueRecord,
			UpdateTime: ts,
		})
	case dt > x.OneDay:
		s := db.UserActionStat{
			WID:        userID,
			ActType:    1,
			ActUnit:    x.UserActStatUintForContinueRecord,
			UpdateTime: ts,
		}
		if isADD {
			s.ActVal = int64(1)
		}
		UpdateUserActionStat(&s)
	}
}

// CreateUser add user
func CreateUser(user db.User) error {
	tx := db.DBCtl.GetDB().MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			fmt.Printf("CreateUser err:%s\n", err)
		}
	}()
	res := tx.MustExec(db.CreateUserSQL, user.WID, user.ExtID, user.WType, user.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		fmt.Printf("CreateUser err:%s row:%d\n", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// QueryUserByWID id
func QueryUserByWID(wID string) (user *db.User) {
	user = &db.User{}
	db.DBCtl.GetDB().Get(user, db.QueryUserByWIDSQL, wID)
	return user
}

// QueryUserByExtID id
func QueryUserByExtID(oID string) (user *db.User) {
	user = &db.User{}
	db.DBCtl.GetDB().Get(user, db.QueryUserByExtIDSQL, oID)
	return user
}

// CreateNote add user
func CreateNote(note *db.WNote) error {
	tx := db.DBCtl.GetDB().MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("CreateNote err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(db.CreateNoteSQL, note.WID, note.NID, note.WMood, note.WDesc, note.WAction, note.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateNote err:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// DelNote update user
func DelNote(note *db.WNote) error {
	tx := db.DBCtl.GetDB().MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("DelNote err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(db.DELNoteSQL, note.DeleteTime, note.NID)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("DelNote:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// UpdateNote update user
func UpdateNote(note *db.WNote) error {
	tx := db.DBCtl.GetDB().MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("UpdateNote err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(db.UpdateNoteSQL, note.WMood, note.WDesc, note.WAction, note.UpdateTime, note.WID, note.NID)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("UpdateNote:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// QueryNotesByWID get notes
func QueryNotesByWID(wID string, limit, offset int) (notes []*db.WNote) {
	notes = make([]*db.WNote, 0)
	db.DBCtl.GetDB().Select(&notes, db.QueryNotesByIDSQL, wID, limit, offset)
	return notes
}

// QueryTotalNoteNumByUserID get total num
func QueryTotalNoteNumByUserID(userID string) (num int) {
	db.DBCtl.GetDB().Get(&num, db.QueryTotalNoteNumByUserIDSQL, userID)
	return num
}

// QueryNotesByOneDay get total num
func QueryNotesByOneDay(userID string, start, end int64, limit, offset int) (notes []*db.WNote) {
	notes = make([]*db.WNote, 0)
	db.DBCtl.GetDB().Select(&notes, db.QueryNotesByOneDaySQL, userID, start, end, limit, offset)
	return notes
}

// QueryNoteTimeByUserIDAndTimeRange time list
func QueryNoteTimeByUserIDAndTimeRange(wID string, start, end int64) (times []int64) {
	times = make([]int64, 0)
	db.DBCtl.GetDB().Select(&times, db.QueryNoteTimeByUserIDAndTimeRangeSQL, wID, start, end)
	return times
}

// QueryUserActionStatValByTypeAndUint  action stat
func QueryUserActionStatValByTypeAndUint(wID string, t int, u string) (stat *db.UserActionStat) {
	stat = &db.UserActionStat{}
	db.DBCtl.GetDB().Get(stat, db.QueryUserActionStatValByTypeAndUintSQL, wID, t, u)
	return stat
}

// CreateUserActionStat  create action stat
func CreateUserActionStat(stat *db.UserActionStat) {
	tx := db.DBCtl.GetDB().MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("CreateUserActionStat err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(db.CreateUserActionStatSQL, stat.WID, stat.ActType, stat.ActVal, stat.ActUnit, stat.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateUserActionStat:%s row:%d", err.Error(), row)
		return
	}
	tx.Commit()
}

// UpdateUserActionStat  update action stat
func UpdateUserActionStat(stat *db.UserActionStat) {
	tx := db.DBCtl.GetDB().MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("UpdateUserActionStat err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(db.UpdateUserActionStatSQL, stat.ActType, stat.ActVal, stat.ActUnit, stat.UpdateTime, stat.WID)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateUserActionStat:%s row:%d", err.Error(), row)
		return
	}
	tx.Commit()
}
