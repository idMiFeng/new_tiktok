package global

import "errors"

var (
	Success                    = errors.New("success")
	ErrUserNotFound            = errors.New("用户不存在")
	ErrUsernameAlreadyExist    = errors.New("用户名已存在")
	ErrInvalidUsernameOrPass   = errors.New("用户名或密码无效")
	ErrFailedToCreateUser      = errors.New("创建用户失败")
	ErrFailedToGetUserInfo     = errors.New("获取用户信息失败")
	ErrFailedToGetUserCount    = errors.New("获取用户统计信息失败")
	ErrFailedToCreateUserCount = errors.New("创建用户统计信息失败")
)

const (
	SuccessCode                    = 0
	ErrCodeUserNotFound            = 4000
	ErrCodeUsernameAlreadyExist    = 4001
	ErrCodeInvalidUsernameOrPass   = 4002
	ErrCodeFailedToCreateUser      = 4003
	ErrCodeFailedToGetUserInfo     = 4004
	ErrCodeFailedToGetUserCount    = 4005
	ErrCodeFailedToCreateUserCount = 4006
)

var errorMsg = map[int]string{
	SuccessCode:                    Success.Error(),
	ErrCodeUserNotFound:            ErrUserNotFound.Error(),
	ErrCodeUsernameAlreadyExist:    ErrUsernameAlreadyExist.Error(),
	ErrCodeInvalidUsernameOrPass:   ErrInvalidUsernameOrPass.Error(),
	ErrCodeFailedToCreateUser:      ErrFailedToCreateUser.Error(),
	ErrCodeFailedToGetUserInfo:     ErrFailedToGetUserInfo.Error(),
	ErrCodeFailedToGetUserCount:    ErrFailedToGetUserCount.Error(),
	ErrCodeFailedToCreateUserCount: ErrFailedToCreateUserCount.Error(),
}

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	return errorMsg[code]
}
