package model

type ErrMod struct {
	Status int   ` json:"status"` //指http状态码
	Err    error `json:"err"`     //错误信息
}
