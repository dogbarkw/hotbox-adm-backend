package until

import (
	"hotbox-adm-backend/pkg/errno"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
type Cont struct {
	c      *gin.Context
	reData *Response
}

func (r Cont) Responses(code int, msg string, data interface{}) {
	r.c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func NewResponse(c *gin.Context) *Cont {
	res := new(Cont)
	res.c = c
	res.reData = new(Response)
	return res
}

func (r Cont) ResponseSuccessWithList(data interface{}, total int) {
	code := errno.Success
	msg := errno.GetMsg(errno.Success)
	r.c.JSON(200, gin.H{
		"code":      code,
		"msg":       msg,
		"data":      data,
		"res_count": total,
	})
}

func (r Cont) ResponseSuccess(data interface{}) {
	code := errno.Success
	msg := errno.GetMsg(errno.Success)
	r.Responses(code, msg, data)
}

func (r Cont) ResponseFail(msg string) {
	code := errno.Error
	r.Responses(code, msg, struct{}{})
}

func (r Cont) ResponseFails(msg string, code int) {
	r.Responses(code, msg, struct{}{})
}
