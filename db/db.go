package db

import (
	"fmt"
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
	NID        string `db:"n_id"`
	WID        string `db:"w_id"`
	WMood      int    `db:"w_mood"`
	WDesc      string `db:"w_desc"`
	WAction    int    `db:"w_action_type"`
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
			fmt.Printf("CreateUser err:%s\n", err)
		}
	}()
	res := tx.MustExec(CreateUserSQL, user.WID, user.ExtID, user.WType, user.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		fmt.Printf("CreateUser err:%s row:%d\n", err.Error(), row)
		return err
	}
	tx.Commit()
	return nil
}

// QueryUserByWID id
func QueryUserByWID(wID string) (user *User) {
	user = &User{}
	db.Get(user, QueryUserByWIDSQL, wID)
	return user
}

// QueryUserByExtID id
func QueryUserByExtID(oID string) (user *User) {
	user = &User{}
	db.Get(user, QueryUserByExtIDSQL, oID)
	return user
}

// CreateNote add user
func CreateNote(note *WNote) error {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("CreateNote err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(CreateNoteSQL, note.WID, note.NID, note.WMood, note.WDesc, note.WAction, note.CreateTime)
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
			fmt.Printf("DelNote err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(DELNoteSQL, note.DeleteTime, note.NID)
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
			fmt.Printf("UpdateNote err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(UpdateNoteSQL, note.WMood, note.WDesc, note.WAction, note.UpdateTime, note.WID, note.NID)
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
	notes = make([]*WNote, 0)
	db.Select(&notes, QueryNotesByIDSQL, wID, limit, offset)
	return notes
}

// QueryTotalNoteNumByUserID get total num
func QueryTotalNoteNumByUserID(userID string) (num int) {
	db.Get(&num, QueryTotalNoteNumByUserIDSQL, userID)
	return num
}

// QueryNotesByOneDay get total num
func QueryNotesByOneDay(userID string, start, end int64, limit, offset int) (notes []*WNote) {
	notes = make([]*WNote, 0)
	db.Select(&notes, QueryNotesByOneDaySQL, userID, start, end, limit, offset)
	return notes
}

// QueryNoteTimeByUserIDAndTimeRange time list
func QueryNoteTimeByUserIDAndTimeRange(wID string, start, end int64) (times []int64) {
	times = make([]int64, 0)
	db.Select(&times, QueryNoteTimeByUserIDAndTimeRangeSQL, wID, start, end)
	return times
}

// QueryUserActionStatValByTypeAndUint  action stat
func QueryUserActionStatValByTypeAndUint(wID string, t int, u string) (stat *UserActionStat) {
	stat = &UserActionStat{}
	db.Get(stat, QueryUserActionStatValByTypeAndUintSQL, wID, t, u)
	return stat
}

// CreateUserActionStat  create action stat
func CreateUserActionStat(stat *UserActionStat) {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("CreateUserActionStat err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(CreateUserActionStatSQL, stat.WID, stat.ActType, stat.ActVal, stat.ActUnit, stat.CreateTime)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateUserActionStat:%s row:%d", err.Error(), row)
		return
	}
	tx.Commit()
}

// UpdateUserActionStat  update action stat
func UpdateUserActionStat(stat *UserActionStat) {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("UpdateUserActionStat err:%s\n", err)
			tx.Rollback()
		}
	}()
	res := tx.MustExec(UpdateUserActionStatSQL, stat.ActType, stat.ActVal, stat.ActUnit, stat.UpdateTime, stat.WID)
	if row, err := res.RowsAffected(); err != nil || row == 0 {
		tx.Rollback()
		log.Printf("CreateUserActionStat:%s row:%d", err.Error(), row)
		return
	}
	tx.Commit()
}
