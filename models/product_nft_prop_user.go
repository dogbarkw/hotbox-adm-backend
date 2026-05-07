package models

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"hotbox-adm-backend/cli"

	"gopkg.in/guregu/null.v4"
)

var NftPropUserDal = &NftPropUserModel{}

type NftPropUserModel struct {
	PropUserUuid string    `json:"prop_user_uuid"`
	UserId       uint64    `json:"user_id"`
	PropId       uint64    `json:"prop_id"`
	StartTime    uint64    `json:"start_time"`
	EndTime      uint64    `json:"end_time"`
	IsDelete     uint32    `json:"is_delete"`
	Source       string    `json:"source"`
	BenefitType  string    `json:"benefit_type"`
	IsUsable     uint32    `gorm:"column:is_usable;default:1" json:"is_usable"`
	IsTransfer   uint32    `gorm:"column:is_transfer;default:0" json:"is_transfer"`
	IsOpen       uint32    `gorm:"column:is_open;default:0" json:"is_open"`
	IsCanOpen    uint32    `gorm:"column:is_can_open;default:1" json:"is_can_open"`
	Remark       string    `gorm:"column:remark" json:"remark"`
	UseTime      null.Time `gorm:"column:use_time;default:null" json:"use_time"`
	CreateTime   time.Time `gorm:"column:create_time;AUTOCREATETIME" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;AUTOUPDATETIME" json:"update_time"`
}

func (NftPropUserModel) TableName() string {
	return "ai_match_product_nft_prop_user"
}

// 获取用户拥有指定道具个数
func (p *NftPropUserModel) UserTotalPropNum(ctx context.Context, propId uint64) (num int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(p).
		Select("COUNT(*) AS num")
	err = query.
		Where("prop_id = ?", propId).
		Where("user_id not in (?)", cli.SpecialUserIds). // 排除特殊用户
		Where("is_lock = 0").
		Where("is_usable = 1").
		Where("is_delete = 0").
		Where("end_time > ?", time.Now().UnixMilli()).
		First(&num).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return num, nil
}
