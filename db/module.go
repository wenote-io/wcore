package db

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
