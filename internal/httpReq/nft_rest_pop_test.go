package httpReq

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"hotbox-adm-backend/dto"
)

func TestNftRestPop(t *testing.T) {
	gotResult, err := NftRestPop("5075bcd0d1d048c4ba19b5dc2176c080", dto.NftRestPopSendReq{
		ProductId:     381,
		ProductSizeId: 379,
		Limit:         1,
	})
	assert.Nil(t, err)
	t.Log(gotResult)
}
