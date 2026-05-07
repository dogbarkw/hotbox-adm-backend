package httpReq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/openai/openai-go/shared"
	"github.com/spf13/cast"

	"github.com/cloudwego/kitex/tool/internal_pkg/log"
	"github.com/pkg/errors"

	"hotbox-adm-backend/dto"
)

var ChatgtpApiDal = &ChatgtpApi{}

type ChatgtpApi struct{}

// SendMsg 发送消息无上下文
func (i *ChatgtpApi) SendMsg(ctx context.Context, prompt string, userMsg string, enableGpt bool) (msgList []dto.GptMsg, resp dto.ChatgptSendResp, err error) {
	msgReq := make([]dto.GptMsg, 0)
	if prompt != "" {
		msgReq = append(msgReq, dto.GptMsg{
			Role:    "system",
			Content: prompt,
		})
	}
	if userMsg != "" {
		msgReq = append(msgReq, dto.GptMsg{
			Role:    "user",
			Content: userMsg,
		})
	}
	msgList = msgReq
	if enableGpt { // 是否调用gpt
		resp, err = i.SendGptMsg(ctx, msgReq)
	}
	return
}

// SendGptMsg 发送消息,带上下文
func (i *ChatgtpApi) SendGptMsg(ctx context.Context, userMsg []dto.GptMsg) (resp dto.ChatgptSendResp, err error) {
	if len(userMsg) == 0 {
		return dto.ChatgptSendResp{}, errors.New("message is empty")
	}

	c := http.Client{}
	if strings.ToLower(os.Getenv("ENV")) == "qa" {
		return
	}
	c.Timeout = 5 * time.Minute // 设置超时时间

	if strings.ToLower(os.Getenv("ENV")) == "dev" {
		urli := url.URL{}
		urlproxy, _ := urli.Parse("http://127.0.0.1:7890")
		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		}
	}

	req := dto.ChatgptSendReq{
		Model:       os.Getenv("GPT_MODEL"),
		Messages:    userMsg,
		Temperature: cast.ToFloat64(os.Getenv("GPT_TEMPERATURE")),
	}

	b, _ := json.Marshal(req)
	res, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(b))
	if err != nil {
		err = errors.Wrapf(err, "chatgptClients.SendMsg,NewRequestErr,req:%+v", string(b))
		return
	}

	gptKeyArray := strings.Split(os.Getenv("GPT_KEY"), ",")
	gptKey := gptKeyArray[rand.Intn(len(gptKeyArray))]
	res.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gptKey))
	res.Header.Set("Content-Type", "application/json")
	r, err := c.Do(res)
	if err != nil {
		err = errors.Wrapf(err, "chatgptClients.SendMsg,DoRequestErr,req:%+v", string(b))
		return
	}

	httpResp, _ := io.ReadAll(r.Body)
	log.Infof("chatgptClients.SendMsg, httpResp:%s", string(httpResp))
	if r.StatusCode != 200 {
		err = fmt.Errorf("gptResponseErr, errorCode=%d, response=%s", r.StatusCode, string(httpResp))
		FeiShuGptErrRootBot(fmt.Sprintf("SendGptMsg, chatGptCompletionApiErr, errorCode=%d, response=%s, secretKey:%s", r.StatusCode, string(httpResp), gptKey))
		return
	}

	err = json.Unmarshal(httpResp, &resp)

	logrus.Infof("SendGptMsg, GptUsage:%+v", resp.Usage)
	return
}

// SendGptMsgNew 发送消息,带上下文,支持流模式
func (i *ChatgtpApi) SendGptMsgStream(ctx context.Context, userMsg []dto.GptMsg) (chatCompletionStream *ssestream.Stream[openai.ChatCompletionChunk], err error) {
	if len(userMsg) == 0 {
		return nil, errors.New("message is empty")
	}

	if strings.ToLower(os.Getenv("ENV")) == "qa" {
		return
	}

	logrus.Infof("SendGptMsgStream, gptInput:%+v", userMsg)

	gptKeyArray := strings.Split(os.Getenv("GPT_KEY"), ",")
	gptKey := gptKeyArray[rand.Intn(len(gptKeyArray))]
	client := openai.NewClient(
		option.WithAPIKey(gptKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
	)

	msgsParams := make([]openai.ChatCompletionMessageParamUnion, 0, len(userMsg))
	for _, msg := range userMsg {
		msgsParams = append(msgsParams, openai.ChatCompletionUserMessageParam{
			Role:    openai.F(openai.ChatCompletionUserMessageParamRole(msg.Role)),
			Content: openai.F[openai.ChatCompletionUserMessageParamContentUnion](shared.UnionString(msg.Content)),
		})
	}

	chatCompletionStream = client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages:    openai.F(msgsParams),
		Model:       openai.F(openai.ChatModel(os.Getenv("GPT_MODEL"))),
		Temperature: openai.F[float64](cast.ToFloat64(os.Getenv("GPT_TEMPERATURE"))),
	})
	if chatCompletionStream.Err() != nil {
		FeiShuGptErrRootBot(fmt.Sprintf("SendGptMsgStream, chatGptCompletionStream.Err:%+v, secretKey:%s", chatCompletionStream.Err(), gptKey))
		return nil, chatCompletionStream.Err()
	}
	return
}
