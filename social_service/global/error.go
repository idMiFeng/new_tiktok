package global

import "errors"

var (
	Success = errors.New("成功")

	VideoUnExist    = errors.New("视频不存在")
	VideoUploadErr  = errors.New("视频上传错误")
	UserNoFavorite  = errors.New("用户没有喜欢的视频")
	UserNoVideo     = errors.New("用户没有视频")
	CacheError      = errors.New("缓存错误")
	MQSendErr       = errors.New("消息队列发送失败")
	MYSQLQueryErr   = errors.New("数据库查询错误")
	RequestParamErr = errors.New("请求参数错误")
)

const (
	SuccessCode = 0

	VideoUnExistCode    = 5000
	VideoUploadErrCode  = 5001
	UserNoFavoriteCode  = 5002
	UserNoVideoCode     = 5003
	CacheErrorCode      = 5004
	MQSendErrCode       = 5005
	MYSQLQueryErrCode   = 5006
	RequestParamErrCode = 5007
)

var errorMsg = map[int]string{
	SuccessCode: Success.Error(),

	VideoUnExistCode:    VideoUnExist.Error(),
	VideoUploadErrCode:  VideoUploadErr.Error(),
	UserNoFavoriteCode:  UserNoFavorite.Error(),
	UserNoVideoCode:     UserNoVideo.Error(),
	CacheErrorCode:      CacheError.Error(),
	MQSendErrCode:       MQSendErr.Error(),
	MYSQLQueryErrCode:   MYSQLQueryErr.Error(),
	RequestParamErrCode: RequestParamErr.Error(),
}

// GetErrorMessage 根据错误码获取错误信息
func GetErrorMessage(code int) string {
	return errorMsg[code]
}
