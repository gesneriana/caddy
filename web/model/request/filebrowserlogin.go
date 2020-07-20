package request

// FileBrowserLogin 登录请求参数
type FileBrowserLogin struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	Recaptcha string `json:"recaptcha"`
}
