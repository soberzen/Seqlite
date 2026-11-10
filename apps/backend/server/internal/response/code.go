// response/code.go
package response

const (
	// 成功
	CodeSuccess = 0

	// 4xxxx 段：客户端错误
	CodeInvalidParam = 40000
	CodeUnauthorized = 40100
	CodeTokenExpired = 40101
	CodeInvalidToken = 40102
	CodeInvalidCreds = 40103
	CodeForbidden    = 40300
	CodeNotFound     = 40400
	CodeUserNotFound = 40401
	CodeConflict     = 40900
	CodeEmailTaken   = 40901

	// 5xxxx 段：服务端错误
	CodeInternal = 50000
)
