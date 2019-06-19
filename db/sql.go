package db

const (
	// CreateUserSQL sql
	CreateUserSQL = "INSERT INTO w_user (w_id, ext_id, w_type,create_time) VALUES (?, ?, ?, ?)"
	// QueryUserByWIDSQL sql
	QueryUserByWIDSQL = "SELECT * FROM w_user WHERE w_id=? AND delete_time=0"
	// QueryUserByExtIDSQL sql
	QueryUserByExtIDSQL = "SELECT * FROM w_user WHERE ext_id=? AND delete_time=0"
	// CreateNoteSQL sql
	CreateNoteSQL = "INSERT INTO w_note (w_id,n_id,w_mood, w_desc, w_action_type,create_time) VALUES (?, ?, ?, ?, ?, ?)"
	// DELNoteSQL sql
	DELNoteSQL = "UPDATE  w_note SET delete_time=? WHERE w_id=? AND n_id=?"
	// UpdateNoteSQL sql
	UpdateNoteSQL = "UPDATE  w_note SET w_mood=?, w_desc=?, w_action_type=?,update_time=? WHERE delete_time=0 AND w_id=? AND n_id=?"
	// QueryNotesByIDSQL sql
	QueryNotesByIDSQL = "SELECT * FROM w_note WHERE w_id=? AND delete_time=0 ORDER BY create_time DESC LIMIT ? OFFSET ?"
	// QueryNoteTimeByUserIDAndTimeRangeSQL sql
	QueryNoteTimeByUserIDAndTimeRangeSQL = "SELECT create_time FROM w_note WHERE  delete_time=0 AND w_id=? AND create_time BETWEEN ? AND ?"
	// QueryUserActionStatValByTypeAndUintSQL  sql
	QueryUserActionStatValByTypeAndUintSQL = "SELECT * FROM user_action_stat WHERE  w_id=? AND act_type=? AND act_unit=?"
	// CreateUserActionStatSQL  sql
	CreateUserActionStatSQL = "INSERT INTO user_action_stat (w_id, act_type, act_val, act_unit, create_time) VALUES (?, ?, ?, ?, ?)"
	// UpdateUserActionStatSQL  sql
	UpdateUserActionStatSQL = "UPDATE  user_action_stat SET act_type=?, act_val=?, act_unit=?, update_time=? WHERE delete_time=0 AND w_id=?"
)
