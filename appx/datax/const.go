package datax

const Authorization = "Authorization"

const (
	NotFound = -1
)
// 在每个模型中的意义不同
const (
	IndexStatusActive   = iota + 1 // 正常(用户)
	IndexStatusDisabled = 2        // 禁用(用户)
	IndexStatusDraft    = 3        // 草稿(用户)
	IndexStatusLock     = 20       // 锁定(管理员)
	IndexStatusDelete   = 21       // 删除(管理员)
	IndexStatusWaring   = 22       // 警告(管理员)
	IndexStatusCheck    = 23       // 待审核(管理员)
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusDraft    = "draft"
	StatusLock     = "lock"
	StatusWarning  = "warning"
	StatusDelete   = "delete"
	StatusCheck    = "check"
)

const (
	Yes           = 1
	No            = 2
	TextTrue      = "true"
	TextFalse     = "false"
	TextYes       = "yes"
	TextNo        = "no"
	OK            = "OK"
	SUCCESS       = "SUCCESS"
)


const (
	MIMEApplicationJSON     = "application/json"
	MIMEApplicationForm     = "application/x-www-form-urlencoded"
	MIMEApplicationFormData = "application/form-data"
	MIMEMultipartForm       = "multipart/form-data"
	MIMEMultipartMixed      = "multipart/mixed"
)


const (
	FormatJson   = "json"   // 直接输出
	Format304    = "304"    // 304 跳转
	FormatCookie = "cookie" // 设置 cookie
	SuccessTag   = "ok"
	Exit         = "Nil"
)
