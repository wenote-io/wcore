package db

import (
	"log"
	// mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// User moudle
type User struct {
	ID         int    `db:"id"`
	WID        string `db:"w_id"`
	ExtID      string `db:"ext_id"`
	WType      int    `db:"w_type"`
	CreateTime int64  `db:"create_time"`
	DeleteTime int64  `db:"delete_time"`
	UpdateTime int64  `db:"update_time"`
}

// WNote moudle
type WNote struct {
	ID         int    `db:"id"`
	WID        string `db:"w_id"`
	WMood      int    `db:"w_mood"`
	WDesc      string `db:"w_desc"`
	WAction    int    `db:"w_action"`
	CreateTime int64  `db:"create_time"`
	DeleteTime int64  `db:"delete_time"`
	UpdateTime int64  `db:"update_time"`
}

// UserActionStat stat info
type UserActionStat struct {
	ID         int    `db:"id"`
	WID        string `db:"w_id"`
	ActType    int    `db:"act_type"`
	ActVal     int64  `db:"act_val"`
	ActUnit    string `db:"act_unit"`
	CreateTime int64  `db:"create_time"`
	DeleteTime int64  `db:"delete_time"`
	UpdateTime int64  `db:"update_time"`
}

// DB Connect
var db *sqlx.DB

// InitDB init tidb
func InitDB() {
	var err error
	db, err = sqlx.Connect("mysql", "root:@(172.21.0.6:4000)/wcore")
	if err != nil {
		log.Fatalln(err)
	}
}

// CreateUser add user
func CreateUser(user User) error {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	res := tx.MustExec(CreateUserSQL, user.WID, user.ExtID, user.WType, user.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateNote err:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// QueryUserByWID id
func QueryUserByWID(wID string) (user *User) {
	db.Select(user, QueryUserByWIDSQL, wID)
	return user
}

// QueryUserByExtID id
func QueryUserByExtID(oID string) (user *User) {
	db.Select(user, QueryUserByExtIDSQL, oID)
	return user
}

// CreateNote add user
func CreateNote(note *WNote) error {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	res := tx.MustExec(CreateNoteSQL, note.WID, note.WMood, note.WDesc, note.WAction, note.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateNote err:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// DelNote update user
func DelNote(note *WNote) error {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	res := tx.MustExec(DELNoteSQL, note.DeleteTime, note.WID)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("DelNote:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// UpdateNote update user
func UpdateNote(note *WNote) error {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	res := tx.MustExec(UpdateNoteSQL, note.WID, note.WMood, note.WDesc, note.WAction, note.UpdateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("UpdateNote:%s row:%d", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// QueryNotesByWID get notes
func QueryNotesByWID(wID string, limit, offset int) (notes []*WNote) {
	db.Select(&notes, QueryNotesByIDSQL, wID, limit, offset)
	return notes
}

// QueryNoteTimeByUserIDAndTimeRange time list
func QueryNoteTimeByUserIDAndTimeRange(wID string, start, end int64) (times []int64) {
	db.Select(&times, QueryNoteTimeByUserIDAndTimeRangeSQL, wID, start, end)
	return times
}

// QueryUserActionStatValByTypeAndUint  action stat
func QueryUserActionStatValByTypeAndUint(wID string, t int, u string) (stat *UserActionStat) {
	db.Select(&stat, QueryUserActionStatValByTypeAndUintSQL, wID, t, u)
	return stat
}

// CreateUserActionStat  create action stat
func CreateUserActionStat(stat *UserActionStat) {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	tx.MustExec(CreateUserActionStatSQL, stat.WID, stat.ActType, stat.ActVal, stat.ActUnit, stat.CreateTime)
	tx.Commit()
}

// UpdateUserActionStat  update action stat
func UpdateUserActionStat(stat *UserActionStat) {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	tx.MustExec(UpdateUserActionStatSQL, stat.ActType, stat.ActVal, stat.ActUnit, stat.CreateTime, stat.WID)
	tx.Commit()
}
