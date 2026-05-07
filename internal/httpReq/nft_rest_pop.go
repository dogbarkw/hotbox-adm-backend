package httpReq

import (
	"errors"
	"os"

	"hotbox-adm-backend/dto"
)

func NftRestPop(token string, payload dto.NftRestPopSendReq) (result dto.NftRestPopSendResp, err error) {
	errResult := dto.CommonResp{}
	client := NewClient()
	resp, err := client.R().
		SetBody(&payload).
		SetHeaders(map[string]string{
			"token": token,
		}).
		SetErrorResult(errResult).
		SetSuccessResult(&result).
		Post(os.Getenv("AI_CHAO_LIU_APP_URL") + "/hotbox/v1/operation/product/nft/rest/pop")
	if err != nil {
		return result, err
	}
	if !resp.IsSuccessState() {
		return result, errors.New("bad response status: " + resp.Status)
	}
	if result.Code != 0 {
		return result, errors.New(result.Msg)
	}
	return result, nil
}
