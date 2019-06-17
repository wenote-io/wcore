package x

// User Action Stat Type
const (
	_ = iota
	UserActStatTypeForContinueRecord
)

// User Action Stat Uint
const (
	UserActStatUintForContinueRecord = "sec"
)

// Code
const (
	SuccessCode = 2000
	// 用户不存在
	UserNotFoundErrCode   = 1000
	CreateNoteFailErrCode = 1001
	CreateUserFailErrCode = 1002
)

// msg
const (
	SuccessMsg         = "ok"
	UserNotFoundErrMsg = "UserNotFound"
	CreateNoteFail     = "CreateNoteFail"
	CreateUserFail     = "CreateUserFail"
)

// time
const (
	OneMonth = 30 * 24 * 60 * 60
)

// 小程序 授权
const (
	AppID   = "wx2f90a70b5866bb89"
	Ssecret = "38e8b556b6259a4c5620b247f21688ee"
)
