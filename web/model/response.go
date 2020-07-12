package model

// ResponseData 是响应的json数据
type ResponseData struct {
	State    bool        `json:"state"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data"`
	Error    string      `json:"error"`
	HTTPCode int         `json:"code"`
}
