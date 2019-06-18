package db

const (
	// CreateUserSQL sql
	CreateUserSQL = "INSERT INTO w_user (w_id, ext_id, w_type,create_time) VALUES ($1, $2, $3, $4)"
	// QueryUserByWIDSQL sql
	QueryUserByWIDSQL = "SELECT * FROM w_user WHERE w_id=$1 AND delete_time=0"
	// QueryUserByExtIDSQL sql
	QueryUserByExtIDSQL = "SELECT * FROM w_user WHERE ext_id=$1 AND delete_time=0"
	// CreateNoteSQL sql
	CreateNoteSQL = "INSERT INTO w_note (w_id,w_mood, w_desc, w_action,create_time) VALUES ($1, $2, $3, $4,$5)"
	// DELNoteSQL sql
	DELNoteSQL = "UPDATE  w_note SET delete_time=$1 WHERE w_id=$2"
	// UpdateNoteSQL sql
	UpdateNoteSQL = "UPDATE  w_note SET w_id=$1,w_mood=$2, w_desc=$3, w_action=$4,update_time=$5 WHERE delete_time=0 AND w_id=$6"
	// QueryNotesByIDSQL sql
	QueryNotesByIDSQL = "SELECT * FROM w_note WHERE w_id=$1 AND delete_time=0 ORDER BY create_time DESC LIMIT $2 OFFSET $3"
	// QueryNoteTimeByUserIDAndTimeRangeSQL sql
	QueryNoteTimeByUserIDAndTimeRangeSQL = "SELECT create_time FROM w_note WHERE  delete_time=0 AND w_id=$1 AND create_time BETWEEN $2 AND $3"
	// QueryUserActionStatValByTypeAndUintSQL  sql
	QueryUserActionStatValByTypeAndUintSQL = "SELECT * FROM user_action_stat WHERE  w_id=$1 AND act_type=$2 AND act_unit=$3"
	// CreateUserActionStatSQL  sql
	CreateUserActionStatSQL = "INSERT INTO user_action_stat (w_id, act_type, act_val, act_unit, create_time) VALUES ($1, $2, $3, $4, $5)"
	// UpdateUserActionStatSQL  sql
	UpdateUserActionStatSQL = "UPDATE  user_action_stat SET act_type=$1, act_val=$2, act_unit=$3, create_time=$4 WHERE delete_time=0 AND w_id=$5"
)
