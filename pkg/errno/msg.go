package errno

var MsgFlags = map[int]string{
	Success:         "ok",
	Error:           "休息一下，请稍后再试～",
	InvalidParams:   "参数无效",
	ErrorToken:      "用户凭证无效",
	ManyPeople:      "当前活动人数过多，请稍后再试～",
	TaskTimeFailed:  "配置的开始时间需至少比当前时间晚5分钟",
	RequestTooOfter: "操作太频繁，请稍后重试~",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}
