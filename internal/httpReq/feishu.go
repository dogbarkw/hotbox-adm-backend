package httpReq

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"hotbox-adm-backend/pkg/constant"

	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
)

const (
	feishuPrefix = "kk:newbackend:feishu"
	qaFeishuUrl  = "https://open.feishu.cn/open-apis/bot/v2/hook/2104c392-6d1c-4df2-973d-ea6aca4daece"
)

func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func warpFeiShuBodyPayload(msgType, msg string) struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
} {
	message := msg
	if msgType != "" {
		message = fmt.Sprintf("[%s]\n%s", constant.CMD_TITLE+msgType, msg)
	}
	return struct {
		MsgType string `json:"msg_type"`
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: message,
		},
	}
}

func FeiShuRootBot(format string, msg ...any) error {
	url := qaFeishuUrl
	if strings.ToLower(os.Getenv("ENV")) == "production" {
		url = "https://open.feishu.cn/open-apis/bot/v2/hook/b3e28e3b-8449-4d76-a455-824a1314ec64"
	}
	envMap := map[string]string{
		"qa":         "测试环境",
		"production": "正式环境",
	}

	body := struct {
		MsgType string `json:"msg_type"`
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: fmt.Sprintf(fmt.Sprintf("【%s】 %s", constant.CMD_TITLE+envMap[strings.ToLower(os.Getenv("ENV"))], format), msg...),
		},
	}
	client := req.C().SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	_, err := client.R().
		SetBody(&body).
		Post(url)
	return err
}

func FeiShuDebugRootBot(format string, msg ...any) error {
	url := qaFeishuUrl
	if strings.ToLower(os.Getenv("ENV")) == "production" {
		url = "https://open.feishu.cn/open-apis/bot/v2/hook/c484e622-5675-40d8-9c98-1c1d9e7ed581"
	}
	envMap := map[string]string{
		"qa":         "测试环境",
		"production": "正式环境",
	}

	body := struct {
		MsgType string `json:"msg_type"`
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: fmt.Sprintf(fmt.Sprintf("【%s】 %s", constant.CMD_TITLE+envMap[strings.ToLower(os.Getenv("ENV"))], format), msg...),
		},
	}
	client := req.C().SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	_, err := client.R().
		SetBody(&body).
		Post(url)
	return err
}

func FeiShuGptErrRootBot(format string, msg ...any) error {
	url := qaFeishuUrl
	if strings.ToLower(os.Getenv("ENV")) == "production" {
		url = "https://open.feishu.cn/open-apis/bot/v2/hook/6e27b41f-1992-41b9-bae5-4f815b7f4bfe"
	}
	envMap := map[string]string{
		"dev":        "开发环境",
		"qa":         "测试环境",
		"production": "正式环境",
	}

	body := struct {
		MsgType string `json:"msg_type"`
		Content struct {
			Text string `json:"text"`
		} `json:"content"`
	}{
		MsgType: "text",
		Content: struct {
			Text string `json:"text"`
		}{
			Text: fmt.Sprintf(fmt.Sprintf("【%s】 %s", constant.CMD_TITLE+envMap[strings.ToLower(os.Getenv("ENV"))], format), msg...),
		},
	}
	client := req.C().SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	_, err := client.R().
		SetBody(&body).
		Post(url)
	return err
}

// 根据传入的机器人url发送回执
func FeiShuWithUrlRootBot(url, msgType, msg string) error {
	env := os.Getenv("ENV")
	if strings.ToUpper(env) == "DEV" {
		return nil
	}
	body := warpFeiShuBodyPayload(msgType, msg)
	client := NewFeiShuClient()

	if strings.ToUpper(env) != "PRODUCTION" {
		url = qaFeishuUrl
	}
	resp, err := client.R().
		SetBody(&body).
		Post(url)
	if resp.IsErrorState() {
		logrus.Warnf("飞书请求失败:%s\n body: %s", resp.String(), string(resp.Request.Body))
	}
	if err != nil {
		logrus.Errorf("飞书请求失败:%s\n body: %s", err, string(resp.Request.Body))
	}
	return err
}

func FeiShuWithUrlMentionRootBot(url, msgType, msg string, mentionUIds []string) error {
	env := os.Getenv("ENV")
	if strings.ToUpper(env) == "DEV" {
		return nil
	}
	replyMsg := ""
	if strings.ToUpper(env) == "PRODUCTION" {
		for k, v := range mentionUIds {
			replyMsg += fmt.Sprintf("<at user_id=\"%s\">%s</at>", v, "name")
			if k == len(mentionUIds)-1 {
				replyMsg += "\n"
			}
		}
	}
	replyMsg += msg
	body := warpFeiShuBodyPayload(msgType, replyMsg)
	client := NewClient()

	if strings.ToUpper(env) != "PRODUCTION" {
		url = qaFeishuUrl
	}
	resp, err := client.R().
		SetBody(&body).
		Post(url)
	if resp.IsErrorState() {
		logrus.Warnf("飞书请求失败:%s\n body: %s", resp.String(), string(resp.Request.Body))
	}
	if err != nil {
		logrus.Errorf("飞书请求失败:%s\n body: %s", resp.String(), string(resp.Request.Body))
	}
	return err
}
