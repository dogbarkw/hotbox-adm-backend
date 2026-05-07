package until

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ApiValidator struct {
	err     error
	ReqData interface{}
}

// 获取错误信息
func (f *ApiValidator) GetError() (errMsg string) {
	if f.err == nil {
		return
	}
	if errs, ok := f.err.(validator.ValidationErrors); ok {
		var errList []string
		for _, e := range errs {
			msg := f.TranslateError(e)
			errList = append(errList, msg)
		}
		errMsg = strings.Join(errList, " ")
		return
	}
	if errMsg == "" {
		errMsg = f.err.Error()
	}
	return
}

func (f *ApiValidator) TranslateError(e validator.FieldError) (msg string) {
	field := e.Namespace()
	if strings.ContainsRune(e.Namespace(), '.') {
		field = strings.Split(e.Namespace(), ".")[1]
	}
	tag, _ := GetStructTagContent(f.ReqData, field, "json")
	desc, _ := GetStructTagContent(f.ReqData, field, "desc")
	if desc == "" {
		tag = field
	} else {
		tag = desc
	}

	switch e.ActualTag() {
	case "required":
		msg = fmt.Sprintf("%s字段必填", tag)
	case "min":
		msg = fmt.Sprintf("%s字段最小只能为%s", tag, e.Param())
	case "max":
		msg = fmt.Sprintf("%s字段最大只能为%s", tag, e.Param())
	case "len":
		msg = fmt.Sprintf("%s字段长度必须是%s个字符", tag, e.Param())
	case "oneof":
		msg = fmt.Sprintf("%s字段必须是[%s]中的一个", tag, e.Param())
	default:
		msg = f.err.Error()
	}
	return
}

// 验证入参
func (f *ApiValidator) Validator(c *gin.Context) (ret bool) {
	f.err = c.ShouldBindBodyWith(f.ReqData, binding.JSON)
	body := ""
	if bodyBuff, ok := c.Get(gin.BodyBytesKey); ok {
		body = string(bodyBuff.([]byte))
		body = strings.ReplaceAll(body, "\n", "\\n")
		body = strings.ReplaceAll(body, "\"", "\\\"")
		body = fmt.Sprintf("\"%s\"", body)
	}
	errMsg := ""
	if f.err != nil {
		errMsg = fmt.Sprintf("Error: %s(%s)", f.GetError(), f.err.Error())
		fmt.Println(errMsg)
	}
	ret = nil == f.err
	return
}

func (f *ApiValidator) ValidatorQuery(c *gin.Context) (ret bool) {
	f.err = c.ShouldBindQuery(f.ReqData)
	errMsg := ""
	if f.err != nil {
		errMsg = fmt.Sprintf("Error: %s(%s)", f.GetError(), f.err.Error())
		fmt.Println(errMsg)
	}
	ret = nil == f.err
	return
}
