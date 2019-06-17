package db

const (
	// CreateUserSQL sql
	CreateUserSQL = "INSERT INTO w_user (w_id, ext_id, w_type,create_time) VALUES ($1, $2, $3, $4)"
	// CreateNoteSQL sql
	CreateNoteSQL = "INSERT INTO w_note (w_mood, w_desc, w_action,create_time) VALUES ($1, $2, $3, $4)"
	// QueryNotesByIDSQL sql
	QueryNotesByIDSQL = "SELECT * FROM w_note WHERE w_id=$1 ORDER BY create_time DESC LIMIT $2 OFFSET $3"
	// QueryUserByExtIDSQL sql
	QueryUserByExtIDSQL = "SELECT * FROM w_note WHERE ext_id=$1"
	// QueryUserByWIDSQL sql
	QueryUserByWIDSQL = "SELECT * FROM w_note WHERE w_id=$1"
	// QueryNoteTimeByUserIDAndTimeRangeSQL sql
	QueryNoteTimeByUserIDAndTimeRangeSQL = "SELECT create_time FROM w_note WHERE w_id=$1 AND create_time BETWEEN $2 AND $3"
	// QueryUserActionStatValByTypeAndUintSQL  sql
	QueryUserActionStatValByTypeAndUintSQL = "SELECT act_val FROM user_action_stat WHERE w_id=$1 AND act_type=$2 AND act_unit=$3"
	// CreateUserActionStatSQL  sql
	CreateUserActionStatSQL = "INSERT INTO user_action_stat (w_id, act_type, act_val, act_unit, create_time) VALUES ($1, $2, $3, $4, $5)"
	// UpdateUserActionStatSQL  sql
	UpdateUserActionStatSQL = "UPDATE  user_action_stat SET act_type=$1, act_val=$2, act_unit=$3, create_time=$4 WHERE w_id=$5"
)
