package types

// Resp stands for object to response to client
// 2020/02/19 22:25:52
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
