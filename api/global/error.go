package global

var (
	ErrCodeJWTInvalidToken     = 4000 // 无效的Token
	ErrCodeJWTExpiredToken     = 4001 // Token已过期
	ErrCodeJWTMissingToken     = 4002 // 缺少Token
	ErrCodeJWTInvalidSignature = 4003 // 无效的签名
	ErrBadRequest              = 4004
)
