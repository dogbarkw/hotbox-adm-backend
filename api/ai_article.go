package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models/hd_adb_models"
	"hotbox-adm-backend/models/hd_task_models"
	"hotbox-adm-backend/pkg/constant"

	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	_ "github.com/k0kubun/pp/v3"
)

// @BasePath /hotbox/operation

// @Summary ChatGPT ai公告
// @Description ChatGPT ai公告
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GenAiArticleMsgReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/gen_ai_article [post]
func GenAiArticle(c *gin.Context) {
	req := form.GenAiArticleMsgReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	prompt := constant.GenAiArticlePrompt
	input, err := completeGptInputMsg(c, req.VerifyContent)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	_, gptResp, err := httpReq.ChatgtpApiDal.SendMsg(c, prompt, input, true)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	if len(gptResp.Choices) == 0 {
		response.ResponseFail("GPT无返回，请稍后重试")
		return
	}

	// fmt.Println(fmt.Sprintf("%s\n%s", prompt, input))
	msg := struct {
		GptInput string `json:"gpt_input"`
		GptResp  string `json:"gpt_resp"`
	}{
		GptInput: fmt.Sprintf("%s\n%s", prompt, input),
		GptResp:  gptResp.Choices[0].Message.Content,
	}
	response.ResponseSuccess(msg)
}

// 拼接user input部分
func completeGptInputMsg(ctx *gin.Context, content string) (string, error) {
	yesterday := time.Now().AddDate(0, 0, -1)

	// 健康名单
	healthPre := fmt.Sprintf(`健康名单：
%s健康名单快照
 （从中选择藏品做活动）
`, yesterday.Format("2006-01-02"))

	// 异常名单部分
	unHealthBody := `异常名单：
%s
请按示例给出审核结果与S1-S3输出，不要有计算错误，不要忘记验算流通市值`

	productHealths, err := hd_task_models.DailyProductHealthDal.GetProductHealthByYmd(ctx, yesterday.Format("20060102"))
	if err != nil {
		return "", err
	}

	nftTmp := `
藏品名称: %v 
实时当日平均成交价: %v
`

	healthContent, unHealthContent := "", ""
	avgCostMap := make(map[string]float64)

	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location()).Format("2006-01-02 15:04:05")
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 59, yesterday.Location()).Format("2006-01-02 15:04:05")
	group, _ := errgroup.WithContext(ctx)
	group.SetLimit(50)
	var mu sync.Mutex
	for _, v := range productHealths {
		health := v
		group.Go(func() error {
			avgCost, err := hd_adb_models.AiMatchProductOrder{Ctx: ctx}.GetProductAvgCostByDate(int(health.ProductId), int(health.NftProductSizeId), startOfDay, endOfDay)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrapf(err, "GetProductAvgCostByDateerr,prodId = %v,nftProductSizeId=%v", health.ProductId, health.NftProductSizeId)
			}
			mu.Lock()
			avgCostMap[fmt.Sprintf("%v-%v", health.ProductId, health.NftProductSizeId)] = avgCost.AvgPayAmount
			mu.Unlock()
			return nil
		})
	}
	if err = group.Wait(); err != nil {
		return "", err
	}

	for index, productHealth := range productHealths {
		var (
			infoMap        = make(map[string]any)
			nftName        string
			productAvgCost float64
			// artist, userNumStandard                 string
			// secordBuyRate, avgPrice, productAvgCost float64
			// nftCreateTime, nftArticleTime           int64
		)

		if err := json.Unmarshal([]byte(productHealth.Info), &infoMap); err != nil {
			return "", errors.Wrapf(err, "unmarshal productHealthInfo err,id = %v", productHealth.Id)
		}

		for k, v := range infoMap {
			switch k {
			case "藏品名称":
				nftName = cast.ToString(v)
				// case "艺术家":
				//	artist = cast.ToString(v)
				// case "用户数达标进度":
				//	userNumStandard = cast.ToString(v)
				// case "二级购买比例":
				//	secordBuyRate = decimal.NewFromFloat(cast.ToFloat64(v)).InexactFloat64()
				// case "平均成交价":
				//	avgPrice = decimal.NewFromFloat(cast.ToFloat64(v)).InexactFloat64()
				// case "藏品生成时间":
				//	nftCreateTime = cast.ToInt64(v)
				// case "最近活动公告时间":
				//	nftArticleTime = cast.ToInt64(v)
			}
		}

		if index == 0 {
			healthContent = healthContent + healthPre
		}

		cost, ok := avgCostMap[fmt.Sprintf("%v-%v", productHealth.ProductId, productHealth.NftProductSizeId)]
		if ok {
			productAvgCost = cost
		}
		//avgCost, err := models.AiMatchProductOrder{Ctx: ctx}.GetProductAvgCostByDate(int(productHealth.ProductId), int(productHealth.NftProductSizeId), startOfDay, endOfDay)
		//if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//	return "", errors.Wrapf(err, "GetProductAvgCostByDateerr,prodId = %v,nftProductSizeId=%v", productHealth.ProductId, productHealth.NftProductSizeId)
		//}
		//productAvgCost = avgCost.AvgPayAmount

		nftInfoStr := fmt.Sprintf(nftTmp, nftName, productAvgCost)
		if productHealth.Type == 1 {
			healthContent = healthContent + nftInfoStr
		} else {
			unHealthContent = unHealthContent + nftInfoStr
		}
	}

	return fmt.Sprintf(`%s
%s
%s`, content, healthContent, fmt.Sprintf(unHealthBody, unHealthContent)), nil
}
