package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/models/hd_adb_models"

	"github.com/shopspring/decimal"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/internal/httpReq"

	"hotbox-adm-backend/form"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
	_ "github.com/k0kubun/pp/v3"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// @BasePath /hotbox/operation

// @Summary 获取活动策划方案(带gpt结果)
// @Description 获取活动策划方案(带gpt结果)
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetGptEventPlanningSchemaReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/event_planning/schema [post]
func GetGptEventPlanningSchema(c *gin.Context) {
	req := form.GetGptEventPlanningSchemaReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	validateErr := form.ValidGetGptEventPlanningSchemaReq(req)
	if validateErr != nil {
		response.Responses(errno.Error, validateErr.Error(), nil)
		return
	}
	prompt, userMsg, err := getGptEventPlanningSchemaMsgLogic(c, req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	gptResponse := ""
	msgs := []dto.GptMsg{}
	// 请求gpt
	if os.Getenv("ENABLE_GPT") != "false" {
		msgList, gptRes, err := httpReq.ChatgtpApiDal.SendMsg(c, prompt, userMsg, true)
		if err != nil {
			response.ResponseFail(err.Error())
			return
		}
		if len(gptRes.Choices) == 0 {
			response.ResponseFail("GPT返回数据有误")
			return
		}
		msgs = msgList
		gptResponse = gptRes.Choices[0].Message.Content
	}
	resp := gin.H{
		"code":         0,
		"prompt":       userMsg,
		"gpt_response": gptResponse,
		"msg_list":     msgs,
		"msg":          "success",
	}
	// TODO:对接gpt
	if os.Getenv("EVENT_PLANNING_SWITCH") == "false" {
		delete(resp, "prompt")
	}
	c.JSON(200, resp)
}

// getGptEventPlanningSchemaMsgLogic 获取提示词和用户输入消息
func getGptEventPlanningSchemaMsgLogic(c *gin.Context, req form.GetGptEventPlanningSchemaReq) (string, string, error) {
	var (
		userMsg string
		prompt  string
	)

	newReq, err := fillWithCollectionPayload(c, req)
	if err != nil {
		return "", "", err
	}

	switch req.TemplateType {
	case 1:
		userMsg, err = assembleGeneralActivityCallWord(newReq)
		if err != nil {
			return "", "", err
		}
		prompt = constant.GeneralActivityCallWordPrefix
	case 2:
		userMsg, err = assembleHighCostActivityCallWord(newReq)
		if err != nil {
			return "", "", err
		}
		prompt = constant.GenHighCostActivityCallWordPrefix
	case 3:
		userMsg, err = assembleLineMergingActivityCallWord(newReq)
		if err != nil {
			return "", "", err
		}
		prompt = constant.GenLineMergingCallWordPrefix
	}
	return prompt, userMsg, err
}

// @Summary 获取活动策划方案(不带gpt结果)
// @Description 获取活动策划方案(不带gpt结果)
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetGptEventPlanningSchemaReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/event_planning/schema_msg [post]
func GetGptEventPlanningSchemaMsg(c *gin.Context) {
	req := form.GetGptEventPlanningSchemaReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	validateErr := form.ValidGetGptEventPlanningSchemaReq(req)
	if validateErr != nil {
		response.Responses(errno.Error, validateErr.Error(), nil)
		return
	}
	prompt, userMsg, err := getGptEventPlanningSchemaMsgLogic(c, req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	msgList, _, _ := httpReq.ChatgtpApiDal.SendMsg(c, prompt, userMsg, false)
	resp := gin.H{
		"code":     0,
		"prompt":   userMsg,
		"msg_list": msgList,
		"msg":      "success",
	}
	// TODO:对接gpt
	if os.Getenv("EVENT_PLANNING_SWITCH") == "false" {
		delete(resp, "prompt")
	}
	c.JSON(200, resp)
}

// @Summary 获取活动策划方案(stream)
// @Description 获取活动策划方案(stream)
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetGptEventPlanningSchemaReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/event_planning/schema_stream [post]
func GetGptEventPlanningSchemaStream(c *gin.Context) {
	req := form.GetGptEventPlanningSchemaReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	validateErr := form.ValidGetGptEventPlanningSchemaReq(req)
	if validateErr != nil {
		response.Responses(errno.Error, validateErr.Error(), nil)
		return
	}
	prompt, userMsg, err := getGptEventPlanningSchemaMsgLogic(c, req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	setSSEHeader(c)
	msgList, _, _ := httpReq.ChatgtpApiDal.SendMsg(c, prompt, userMsg, false)
	completionStream, err := httpReq.ChatgtpApiDal.SendGptMsgStream(c, msgList)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	role := ""
	for completionStream.Next() {
		if len(completionStream.Current().Choices) == 0 {
			continue
		}
		if completionStream.Current().Choices[0].Delta.Role != "" {
			role = string(completionStream.Current().Choices[0].Delta.Role)
		}
		if completionStream.Current().Choices[0].Delta.Content == "" {
			continue
		}

		ev := map[string]any{
			"message": dto.GptMsg{
				Role:    role,
				Content: completionStream.Current().Choices[0].Delta.Content,
			},
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(ev)
		fmt.Fprintf(c.Writer, "data: %v\n", buf.String())

		c.Writer.Flush()
	}
}

// @Summary 获取活动策划方案
// @Description 获取活动策划方案
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.GetGptEventPlanningCollectionInfoReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/event_planning/collection_info [post]
func GetGptEventPlanningCollectionInfo(c *gin.Context) {
	req := form.GetGptEventPlanningCollectionInfoReq{}
	response := until.NewResponse(c)

	validate := until.ApiValidator{ReqData: &req}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	var realUserSurplus int64  // 真实用户剩余份数(公池)
	var restSurplusCount int64 // 私池(剩余份数-公池)
	var nftCount int64         // 剩余份数
	var err error
	g, _ := errgroup.WithContext(c)
	g.Go(func() error {
		return models.GetProductOrderCount(c, cli.SpecialUserIds, &realUserSurplus, req.ProductId, req.NftProductSizeId)
	})
	g.Go(func() error {
		nftMarketWarehouseTotalCount, err := models.NewBusinessNftMarketWarehouseTotalCount(c).GetByProductIdAndSizeId(int64(req.ProductId), int64(req.NftProductSizeId))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		nftCount = nftMarketWarehouseTotalCount.NftCount
		if nftCount == 0 {
			dd, err := models.NewSaleProductNftSize().GetOneByParams(c, map[string]any{
				"product_id":          req.ProductId,
				"nft_product_size_id": req.NftProductSizeId,
			})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			nftCount = dd.TotalCount
		}
		return nil
	})
	err = g.Wait()
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	restSurplusCount = nftCount - realUserSurplus
	resp := gin.H{
		"code": 0,
		"ac":   nftCount,
		"pbac": realUserSurplus,
		"prac": restSurplusCount,
		"msg":  "success",
	}
	c.JSON(200, resp)
}

func parseTpl(name, tpl string, desc any) (string, error) {
	funcMap := template.FuncMap{
		"convertDecimal": func(s any) string {
			t := cast.ToFloat64(s)
			return fmt.Sprintf("%.2f", t)
		},
	}
	var buffer bytes.Buffer
	t := template.New(name)
	t = t.Funcs(funcMap)
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}

	err = t.Execute(&buffer, desc)
	if err != nil {
		return "", err
	}
	defer buffer.Reset()
	return buffer.String(), nil
}

// 组装一般活动提示词
func assembleGeneralActivityCallWord(newReq form.GetGptEventPlanningSchemaReq) (string, error) {
	materialTpl := `可使用的材料列表如下:
{{- range .CollectionPayload}}
{{.ProductTitle}}(lp={{.Lp}},ct={{.Ct}},ac={{.Ac}},pbac={{.Pbac}},prac={{.Prac}})
{{- end }}
活动条件:
`
	materialText, err := parseTpl("materialTpl", materialTpl, newReq)
	if err != nil {
		return "", err
	}

	mainMaterial := lo.Filter[form.GptCollectionArr](newReq.CollectionPayload, func(item form.GptCollectionArr, index int) bool {
		return item.IsMain
	})
	s1Struct := struct {
		ProductTitles string
		Quantitys     string
	}{
		ProductTitles: strings.Join(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return item.ProductTitle
		}), ","),
		Quantitys: strings.Join(lo.Compact(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return cast.ToString(item.Quantity)
		})), ","),
	}

	s1 := `S1
{{if .ProductTitles }}本次活动指定{{.ProductTitles}}为主材料，必须被使用。{{if .Quantitys }}并且系数必须为{{.Quantitys}}，也就是每套材料必须包括{{.Quantitys}}张{{ .ProductTitles }}{{end}}
可以搭配其他材料，也可以只用{{.ProductTitles}}
通过活动，产出一种新产物，名称必须跟主材料{{.ProductTitles}}的名称风格类似。{{else}}本次活动没有指定主材料{{end}}
`
	s1Text, err := parseTpl("s1Template", s1, s1Struct)
	if err != nil {
		return "", err
	}
	s2Struct := struct {
		ActivityType           *string
		IncreaseProfitMultiple string
		ProductCoefficient     *uint
		CostCeiling            *uint // 成本上限
		CostAdvice             *uint
	}{
		ProductCoefficient:     newReq.ProductCoefficient,
		ActivityType:           newReq.ActivityType,
		IncreaseProfitMultiple: fmt.Sprintf("%.2f", newReq.IncreaseProfitMultiple[0]) + "-" + fmt.Sprintf("%.2f", newReq.IncreaseProfitMultiple[1]),
		CostCeiling:            newReq.CostCeiling,
		CostAdvice:             newReq.CostAdvice,
	}

	s2Tpl := `S2
{{ if .ActivityType }}
本次活动必须使用{{.ActivityType}}类型的活动，不可以使用其他活动类型{{ else }}可以任意使用合成、置换、分解这些类型，但不可以用升级类型{{ end }}
{{- if .ProductCoefficient }}
新产物系数为{{.ProductCoefficient}}{{- end }}
{{- if .CostCeiling }}
请使你的配方tct不超过{{.CostCeiling}}{{- end }}
{{- if .CostAdvice }}
建议是你的配方tct在{{.CostAdvice}}左右{{- end }}
请控制你的lct在1.15~1.55之间
本次活动的增润倍数要求在{{.IncreaseProfitMultiple}}之间随机取一个值。
如果任一条件不满足，需要调整你的材料选择和配方，直至满足条件才可以继续往下运行。
如果连续调整5次还是没有找到合适配方，中断运行，不必再进行本次策划案编写，直接输出“没有找到可以使用的配方。”后终止回答即可。
`
	s2Text, err := parseTpl("s2Tpl", s2Tpl, s2Struct)
	if err != nil {
		return "", err
	}
	s3Struct := struct {
		ProductTitles        string
		ConsumedQuantity     string
		MinimumGuaranteeFund uint
	}{
		ProductTitles: strings.Join(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return item.ProductTitle
		}), ","),
		ConsumedQuantity: strings.Join(lo.Compact(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return cast.ToString(item.ConsumedQuantity)
		})), ","),
		MinimumGuaranteeFund: *newReq.MinimumGuaranteeFund * 10000,
	}
	s3Tpl := `S3
{{if  and (.ProductTitles) (.ConsumedQuantity)}}本次消耗{{.ProductTitles}}{{.ConsumedQuantity}}份,请决定活动参与次数，并且决定公池配额、私池配额{{else}}请决定活动参与次数，并且决定公池配额、私池配额{{end}}，同时注意产物公池容量mt上限(mtmax)为{{.MinimumGuaranteeFund}}。不要忘记验算pbcs≤mtmax 。
本次活动的私池配额prcs请尽可能地大，但公池配额pbcs不允许为0。
`
	s3Text, err := parseTpl("s3Tpl", s3Tpl, s3Struct)
	if err != nil {
		return "", err
	}
	s4Text := `S4
本次活动无门槛。
`
	s5Tpl := `S5
{{if .IsRaid}}可以定时发生、也可以突袭发生。{{else}}不可以突袭发生。{{end}}
`
	s5Text, err := parseTpl("s5Tpl", s5Tpl, newReq)
	if err != nil {
		return "", err
	}
	s6Text := `S6
输出活动策划案。

`
	endingText := `请设计一个活动案，输出你的完整思路（S1-S5）和活动策划案。
简洁地回答，确保符合每一个要求，不要有计算错误。不要漏掉过程中的验算。
不要忘记验算pbcs≤mtmax=mt上限/(产物限价*产物系数)`

	return materialText + s1Text + s2Text + s3Text + s4Text + s5Text + s6Text + endingText, nil
}

// 组装高成本活动
func assembleHighCostActivityCallWord(newReq form.GetGptEventPlanningSchemaReq) (string, error) {
	materialTpl := `可使用的材料列表如下:
{{- range .CollectionPayload}}
{{.ProductTitle}}(lp={{.Lp}},ct={{.Ct}},ac={{.Ac}},pbac={{.Pbac}},prac={{.Prac}})
{{- end}}
活动条件:
`
	materialText, err := parseTpl("materialTpl", materialTpl, newReq)
	if err != nil {
		return "", err
	}

	mainMaterial := lo.Filter[form.GptCollectionArr](newReq.CollectionPayload, func(item form.GptCollectionArr, index int) bool {
		return item.IsMain
	})
	s1Struct := struct {
		ProductTitles string
		Quantitys     string
	}{
		ProductTitles: strings.Join(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return item.ProductTitle
		}), ","),
		Quantitys: strings.Join(lo.Compact(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return cast.ToString(item.Quantity)
		})), ","),
	}

	s1 := `S1
{{if .ProductTitles }}本次活动指定{{.ProductTitles}}为主材料，必须被使用。{{if .Quantitys }}并且系数必须为{{.Quantitys}}，也就是每套材料必须包括{{.Quantitys}}张{{ .ProductTitles }}{{end}}
可以搭配其他材料，也可以只用{{.ProductTitles}}
通过活动，产出一种新产物，名称必须跟主材料{{.ProductTitles}}的名称风格类似。{{else}}本次活动没有指定主材料{{end}}
`
	s1Text, err := parseTpl("s1Template", s1, s1Struct)
	if err != nil {
		return "", err
	}
	s2Struct := struct {
		ActivityType           *string
		IncreaseProfitMultiple string
		ProductCoefficient     *uint
		CostCeiling            *uint // 成本上限
		CostAdvice             *uint
	}{
		ProductCoefficient:     newReq.ProductCoefficient,
		ActivityType:           newReq.ActivityType,
		IncreaseProfitMultiple: fmt.Sprintf("%.2f", newReq.IncreaseProfitMultiple[0]) + "-" + fmt.Sprintf("%.2f", newReq.IncreaseProfitMultiple[1]),
		CostCeiling:            newReq.CostCeiling,
		CostAdvice:             newReq.CostAdvice,
	}

	s2Tpl := `S2
{{ if .ActivityType }}
本次活动必须使用{{.ActivityType}}类型的活动，不可以使用其他活动类型{{ else }}可以任意使用合成、置换、分解这些类型，但不可以用升级类型{{ end }}
{{- if .ProductCoefficient }}
新产物系数为{{.ProductCoefficient}}{{- end }}
{{- if .CostCeiling }}
请使你的配方tct不超过{{.CostCeiling}}{{- end }}
{{- if .CostAdvice }}
建议是你的配方tct在{{.CostAdvice}}左右{{- end }}
请控制你的lct在1.15~1.8之间
本次活动的增润倍数要求在{{.IncreaseProfitMultiple}}之间随机取一个值。
如果任一条件不满足，需要调整你的材料选择和配方，直至满足条件才可以继续往下运行。
如果连续调整5次还是没有找到合适配方，中断运行，不必再进行本次策划案编写，直接输出“没有找到可以使用的配方。”后终止回答即可。
`
	s2Text, err := parseTpl("s2Tpl", s2Tpl, s2Struct)
	if err != nil {
		return "", err
	}
	s3Struct := struct {
		ProductTitles        string
		ConsumedQuantity     string
		MinimumGuaranteeFund uint
	}{
		ProductTitles: strings.Join(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return item.ProductTitle
		}), ","),
		ConsumedQuantity: strings.Join(lo.Compact(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			return cast.ToString(item.ConsumedQuantity)
		})), ","),
		MinimumGuaranteeFund: *newReq.MinimumGuaranteeFund * 10000,
	}
	s3Tpl := `S3
{{if  and (.ProductTitles) (.ConsumedQuantity)}}本次消耗{{.ProductTitles}}{{.ConsumedQuantity}}份,请决定活动参与次数，并且决定公池配额、私池配额{{else}}请决定活动参与次数，并且决定公池配额、私池配额{{end}}，同时注意产物公池容量mt上限(mtmax)为{{.MinimumGuaranteeFund}}。不要忘记验算pbcs≤mtmax 。
本次活动的私池配额prcs请尽可能地大，但公池配额pbcs不允许为0。
`
	s3Text, err := parseTpl("s3Tpl", s3Tpl, s3Struct)
	if err != nil {
		return "", err
	}
	s4Text := `S4
本次活动无门槛。
`
	s5Tpl := `S5
{{if .IsRaid}}可以定时发生、也可以突袭发生。{{else}}不可以突袭发生。{{end}}
`
	s5Text, err := parseTpl("s5Tpl", s5Tpl, newReq)
	if err != nil {
		return "", err
	}
	s6Text := `S6
输出活动策划案。

`
	endingText := `请设计一个活动案，输出你的完整思路（S1-S5）和活动策划案。
简洁地回答，确保符合每一个要求，不要有计算错误。不要漏掉过程中的验算。
不要忘记验算pbcs≤mtmax=mt上限/(产物限价*产物系数)`

	return materialText + s1Text + s2Text + s3Text + s4Text + s5Text + s6Text + endingText, nil
}

// 组装并流提示词
func assembleLineMergingActivityCallWord(newReq form.GetGptEventPlanningSchemaReq) (string, error) {
	mainMaterial := lo.Filter[form.GptCollectionArr](newReq.CollectionPayload, func(item form.GptCollectionArr, index int) bool {
		return item.IsMain
	})
	noMainMaterial := lo.Filter[form.GptCollectionArr](newReq.CollectionPayload, func(item form.GptCollectionArr, index int) bool {
		return !item.IsMain
	})
	materialStruct := struct {
		MainMaterial   []form.GptCollectionArr
		NoMainMaterial []form.GptCollectionArr
	}{
		MainMaterial:   mainMaterial,
		NoMainMaterial: noMainMaterial,
	}
	materialTpl := `{{ if len .MainMaterial }}主材料列表:{{- range .MainMaterial}}
{{.ProductTitle}}(lp={{.Lp}},ct={{.Ct}},ac={{.Ac}},pbac={{.Pbac}},prac={{.Prac}})
{{ end}}{{ end }}
{{- if len .NoMainMaterial }}副材料列表:{{- range .NoMainMaterial}}
{{.ProductTitle}}(lp={{.Lp}},ct={{.Ct}},ac={{.Ac}},pbac={{.Pbac}},prac={{.Prac}})
{{- end}}{{ end }}
`
	materialText, err := parseTpl("materialTpl", materialTpl, materialStruct)
	if err != nil {
		return "", err
	}
	s1Struct := struct {
		MainMaterial     []form.GptCollectionArr
		MainMaterialText string
	}{
		MainMaterial: mainMaterial,
		MainMaterialText: strings.Join(lo.Map[form.GptCollectionArr, string](mainMaterial, func(item form.GptCollectionArr, index int) string {
			if item.ConsumedQuantity != nil {
				if item.Quantity != nil {
					return fmt.Sprintf("%s(消耗%d份,主材料系数为:%d)", item.ProductTitle, *item.ConsumedQuantity, *item.Quantity)
				} else {
					return fmt.Sprintf("%s(消耗%d份)", item.ProductTitle, *item.ConsumedQuantity)
				}
			}
			return item.ProductTitle
		}), "、"),
	}
	s1 := `{{ if len .MainMaterial}}要求根据这些材料，一次性做{{ len .MainMaterial }}个活动。
{{ len .MainMaterial }}个活动分别以：{{ .MainMaterialText }}为主材料，每个活动使用其中一种。{{ end }}`
	s1Text, err := parseTpl("s1Template", s1, s1Struct)
	if err != nil {
		return "", err
	}
	s2Struct := struct {
		MainMaterial []form.GptCollectionArr
		ActivityType *string
	}{
		MainMaterial: mainMaterial,
		ActivityType: newReq.ActivityType,
	}
	s2 := `{{- if .ActivityType }}
{{ if len .MainMaterial}}{{ len .MainMaterial}}个活动都用{{ .ActivityType }}{{ end }}{{ end }}
每个活动可以搭配副材料（从副材料列表中自由选择），或者仅使用主材料。
`
	s2Text, err := parseTpl("s2Template", s2, s2Struct)
	if err != nil {
		return "", err
	}
	s3Struct := struct {
		MainMaterial                 []form.GptCollectionArr
		ProductRecommendedPriceLimit *uint
		NewProductTitle              *string
		ProductLimitRangeText        string
	}{
		MainMaterial:                 mainMaterial,
		ProductRecommendedPriceLimit: newReq.ProductRecommendedPriceLimit,
		NewProductTitle:              newReq.NewProductTitle,
	}
	if newReq.ProductRecommendedPriceLimit != nil {
		s3Struct.ProductLimitRangeText = fmt.Sprintf("%d到%d", newReq.ProductLimitRange[0], newReq.ProductLimitRange[1])
	}
	s3 := `要求{{ if len .MainMaterial }}{{len .MainMaterial}}个{{ end }}活动生成的新产物为同一种（名称、限价相同）,新产物限价为{{ .ProductRecommendedPriceLimit }}(不一定非得是{{.ProductRecommendedPriceLimit}}, {{ .ProductLimitRangeText }}间均可)
新产物名称为{{.NewProductTitle}}
`
	s3Text, err := parseTpl("s3Template", s3, s3Struct)
	if err != nil {
		return "", err
	}

	s4Struct := struct {
		IncreaseProfitMultiple string
	}{
		IncreaseProfitMultiple: fmt.Sprintf("%.2f", newReq.IncreaseProfitMultiple[0]) + "-" + fmt.Sprintf("%.2f", newReq.IncreaseProfitMultiple[1]),
	}
	s4 := `每个活动lct处在1-1.9之间
每个活动增润倍数在{{ .IncreaseProfitMultiple }}之间。
`
	s4Text, err := parseTpl("s4Template", s4, s4Struct)
	if err != nil {
		return "", err
	}
	s5Text := `不可以突袭发生。
无门槛。
`
	s6Struct := struct {
		MainMaterial       []form.GptCollectionArr
		TotalGuaranteeFund *uint
	}{
		MainMaterial:       mainMaterial,
		TotalGuaranteeFund: newReq.TotalGuaranteeFund,
	}
	s6 := `
注意，新产物的总产物公池容量tmt=∑（每个活动的产物公池份数）*产物限价。
tmt上限为{{ .TotalGuaranteeFund }}（如果超过{{ .TotalGuaranteeFund }}，请适当降低其中一个活动的公池配额，也可以降低多个活动的公池配额，来使tmt达到要求。）
并在最后，给出{{ if len .MainMaterial}}{{ len .MainMaterial}}个{{ end}}活动策划案，新产物名称、限价、总份数、公池份数、私池份数、tmt（验算tmt是否符合上限要求），并根据使用情况更新主副材料列表。
`
	s6Text, err := parseTpl("s6Template", s6, s6Struct)
	if err != nil {
		return "", err
	}
	s7Struct := struct {
		MainMaterial       []form.GptCollectionArr
		TotalGuaranteeFund *uint
		DivisionResult     uint
		ResultText         string
	}{
		MainMaterial:       mainMaterial,
		TotalGuaranteeFund: newReq.TotalGuaranteeFund,
	}
	if newReq.TotalGuaranteeFund != nil && len(mainMaterial) > 0 {
		s7Struct.DivisionResult = *newReq.TotalGuaranteeFund / uint(len(mainMaterial))
		s7Struct.ResultText = fmt.Sprintf("每个活动的mt上限为tmt上限/活动个数=%d/%d=%d", *newReq.TotalGuaranteeFund, len(mainMaterial), s7Struct.DivisionResult)
	}
	s7 := `{{ if .ResultText }}{{.ResultText}}{{ end }}不要忘记验算pbcs≤pbmax，建议pbcs≤mtmax*80%`
	s7Text, err := parseTpl("s7Template", s7, s7Struct)
	if err != nil {
		return "", err
	}
	s8Text := `以上任一条件不满足，都需要调整你的材料系数、材料选择和配方，直到满足条件才可以继续往下运行
如果连续调整5次还没有找到合适配方，中断运行，不必再进行本次策划案编写，直接输出“没有找到可用配方”后终止运行，进行下一个配方的任务。
`
	s9Struct := struct {
		MainMaterial []form.GptCollectionArr
	}{
		MainMaterial: mainMaterial,
	}
	s9 := `
好了，请开始做题吧！记得简洁回答，不要有计算错误。
请分步完成任务
1、先按D1-D2定{{ if len .MainMaterial }} {{ len .MainMaterial }}个{{ end }}活动的配方。
等我回复你配方可以后
2、再按运行D3（S3-S6定单个活动的cs、pbcs、prcs）、D4汇总输出最后结果。
每个活动的S1-S6运行完毕后，等我回复你可以，你再运行下一个活动的S1-S6`
	s9Text, err := parseTpl("s9Template", s9, s9Struct)
	if err != nil {
		return "", err
	}
	return materialText + s1Text + s2Text + s3Text + s4Text + s5Text + s6Text + s7Text + s8Text + s9Text, nil
}

// 填充藏品信息
func fillWithCollectionPayload(c *gin.Context, req form.GetGptEventPlanningSchemaReq) (newReq form.GetGptEventPlanningSchemaReq, err error) {
	for _, v := range req.CollectionPayload {
		var realUserSurplus int64   // 真实用户剩余份数(公池)
		var restSurplusCount int64  // 私池(剩余份数-公池)
		var nftCount int64          // 剩余份数
		var sellPriceMaxLimit int64 // 最高限价
		var ct uint                 // 成本
		var err error
		g, _ := errgroup.WithContext(c)
		// 获取 公池(Pbac)
		if v.Pbac != nil {
			realUserSurplus = *v.Pbac
		} else {
			g.Go(func() error {
				return models.GetProductOrderCount(c, cli.SpecialUserIds, &realUserSurplus, v.ProductId, v.NftProductSizeId)
			})
		}
		// 获取 限价(Lp)
		if v.Lp != nil {
			sellPriceMaxLimit = *v.Lp
		} else {
			g.Go(func() error {
				orderData, err := models.NewAiMatchProductNftSecondPrice(c).GetByProductIdAndNftProductSizeId(v.ProductId, v.NftProductSizeId)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				sellPriceMaxLimit = cast.ToInt64(orderData.SellPriceMaxLimit)
				return nil
			})
		}

		// 获取剩余份数（ac）
		if v.Ac != nil {
			nftCount = *v.Ac
		} else {
			g.Go(func() error {
				nftMarketWarehouseTotalCount, err := models.NewBusinessNftMarketWarehouseTotalCount(c).GetByProductIdAndSizeId(int64(v.ProductId), int64(v.NftProductSizeId))
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				nftCount = nftMarketWarehouseTotalCount.NftCount
				if nftCount == 0 {
					dd, err := models.NewSaleProductNftSize().GetOneByParams(c, map[string]any{
						"product_id":          v.ProductId,
						"nft_product_size_id": v.NftProductSizeId,
					})
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
					nftCount = dd.TotalCount
				}
				return nil
			})
		}

		// 获取 私池(prac)
		if v.Prac != nil {
			restSurplusCount = *v.Prac
		} else {
			restSurplusCount = nftCount - realUserSurplus
		}

		// 获取成本
		if v.Ct != nil {
			ct = *v.Ct
		} else {
			g.Go(func() error {
				avgData, err := hd_adb_models.AiMatchProductOrder{Ctx: c}.GetProductAvgCost(v.ProductId, v.NftProductSizeId)
				if err != nil {
					return err
				}
				// 按二级购买平均价来填这个值
				if !decimal.NewFromFloat(avgData.AvgPayAmount).Equal(decimal.NewFromInt(0)) {
					ct = uint(avgData.AvgPayAmount)
					return nil
				}
				// 二级没有发生过购买，按当前最低寄售价来填这个值
				orderData, err := models.NewAiMatchProductNftSecondPrice(c).GetByProductIdAndNftProductSizeId(v.ProductId, v.NftProductSizeId)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if orderData.SaleMinPrice > 0 {
					ct = uint(orderData.SaleMinPrice)
					return nil
				}
				// 如果也没有，再按限价来填
				ct = uint(orderData.SellPriceMaxLimit)
				return nil
			})
		}

		err = g.Wait()
		if err != nil {
			return newReq, err
		}

		pp := form.GptCollectionArr{
			ProductTitle:     v.ProductTitle,
			ProductId:        v.ProductId,
			NftProductSizeId: v.NftProductSizeId,
			IsMain:           v.IsMain,
			Quantity:         v.Quantity,
			ConsumedQuantity: v.ConsumedQuantity,
			Lp:               &sellPriceMaxLimit, // 限价
			Ct:               &ct,                // 成本
			Ac:               &nftCount,          // 剩余份数
			Pbac:             &realUserSurplus,   // 公池图份数
			Prac:             &restSurplusCount,  // 私池图份数
		}
		newReq.CollectionPayload = append(newReq.CollectionPayload, pp)
	}
	if req.IncreaseProfitMultiple == nil {
		var defaultVal *[2]float64
		switch req.TemplateType {
		case 1:
			defaultVal = &[2]float64{1.15, 1.35}
		case 2:
			defaultVal = &[2]float64{1.20, 1.55}
		case 3:
			defaultVal = &[2]float64{1.05, 1.22}
		}
		newReq.IncreaseProfitMultiple = defaultVal
	} else {
		newReq.IncreaseProfitMultiple = req.IncreaseProfitMultiple
	}
	if req.MinimumGuaranteeFund == nil {
		defaultVal := uint(2)
		newReq.MinimumGuaranteeFund = &defaultVal
	} else {
		newReq.MinimumGuaranteeFund = req.MinimumGuaranteeFund
	}
	if req.ProductLimitRange != nil {
		var defaultVal *[2]uint
		switch req.TemplateType {
		// case 1:
		// 	defaultVal = &[2]float64{1.15, 1.35}
		// case 2:
		// 	defaultVal = &[2]float64{1.20, 1.55}
		case 3:
			defaultVal = &[2]uint{148, 200}
		}
		if req.ProductLimitRange != nil {
			newReq.ProductLimitRange = req.ProductLimitRange
		} else {
			newReq.ProductLimitRange = defaultVal
		}

	}

	newReq.ProductCoefficient = req.ProductCoefficient
	newReq.CostCeiling = req.CostCeiling
	newReq.CostAdvice = req.CostAdvice
	newReq.ActivityType = req.ActivityType
	newReq.IsRaid = req.IsRaid
	newReq.TemplateType = req.TemplateType
	newReq.ProductRecommendedPriceLimit = req.ProductRecommendedPriceLimit
	newReq.TotalGuaranteeFund = req.TotalGuaranteeFund
	newReq.NewProductTitle = req.NewProductTitle
	return
}

// @Summary ChatGPT 会话
// @Description ChatGPT 会话
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.SendGptMsgReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/send_gpt_msg [post]
func SendGptMsg(c *gin.Context) {
	req := form.SendGptMsgReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	resp, err := httpReq.ChatgtpApiDal.SendGptMsg(c, req.Msg)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	response.ResponseSuccess(resp)
}

// @Summary ChatGPT 会话 (stream)
// @Description ChatGPT 会话 (stream)
// @Tags Gpt
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object body form.SendGptMsgReq true "查询参数"
// @Success 200 {object} any
// @Router /gpt/send_gpt_msg_stream [post]
func SendGptMsgStream(c *gin.Context) {
	req := form.SendGptMsgReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	setSSEHeader(c)
	completionStream, err := httpReq.ChatgtpApiDal.SendGptMsgStream(c, req.Msg)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}

	role := ""
	for completionStream.Next() {
		if len(completionStream.Current().Choices) == 0 {
			continue
		}
		if completionStream.Current().Choices[0].Delta.Role != "" {
			role = string(completionStream.Current().Choices[0].Delta.Role)
		}
		if completionStream.Current().Choices[0].Delta.Content == "" {
			continue
		}
		//sse.Encode(c.Writer, sse.Event{
		//	Data: map[string]any{
		//		"message": dto.GptMsg{
		//			Role:    role,
		//			Content: completionStream.Current().Choices[0].Delta.Content,
		//		},
		//	},
		//})
		ev := map[string]any{
			"message": dto.GptMsg{
				Role:    role,
				Content: completionStream.Current().Choices[0].Delta.Content,
			},
		}
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(ev)
		fmt.Fprintf(c.Writer, "data: %v\n", buf.String())

		c.Writer.Flush()
	}
}

func setSSEHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
}
