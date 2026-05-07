package models

import (
	"context"
	"strings"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BolingPropRecordsModel struct {
	Id         uint64    `gorm:"column:id" json:"id"`
	Source     uint16    `gorm:"column:source" json:"source"` //   1:绝地求生 2:热狗运动会 10:贝壳扭蛋机
	UserId     uint64    `gorm:"column:user_id" json:"user_id"`
	Num        uint64    `gorm:"column:num" json:"num"`
	RemainNum  uint64    `gorm:"column:remain_num" json:"remain_num"`
	Remark     string    `gorm:"column:remark" json:"remark"`
	ChangeType uint8     `gorm:"column:change_type" json:"change_type"`
	CreateDate string    `gorm:"column:create_date" json:"create_date"`
	Status     uint8     `gorm:"column:status;default:1" json:"status"`
	PropId     uint64    `gorm:"column:prop_id" json:"prop_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (BolingPropRecordsModel) TableName() string {
	return "ai_match_product_boling_prop_num_record"
}

type SendPropReq struct {
	PropUserUuid    string
	UserId          uint64
	PropId          uint64
	Number          uint32
	Tm              time.Time
	Source          string
	OutPropUserUuid string
	PropSource      uint16
	Remark          string
}

// 发道具
func SendProp(ctx context.Context, req SendPropReq) error {
	tx := cli.HotDogGormDB.Begin()
	defer func() {
		if r := recover(); r != nil { // 捕获 panic
			tx.Rollback()
		}
	}()
	prop := AiMatchProductNftPropModel{}
	if err := tx.Model(&AiMatchProductNftPropModel{}).Where("prop_id = ?", req.PropId).First(&prop).Error; err != nil {
		tx.Rollback()
		return err
	}
	if strings.EqualFold(prop.BenefitType, "boling") {
		propNum := AiMatchProductBolingPropNumModel{}
		if err := tx.Model(&propNum).Where("prop_id = ? and user_id = ? and is_delete = ?", req.PropId, req.UserId, 0).First(&propNum).Error; err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			return errors.WithStack(err)
		}

		if propNum.BolingUuid != "" { // 存在 number++
			updateTx := tx.Model(&AiMatchProductBolingPropNumModel{}).
				Where("boling_uuid = ?", propNum.BolingUuid).
				Update("number", gorm.Expr("number + ?", req.Number))
			if updateTx.Error != nil {
				tx.Rollback()
				return errors.WithStack(updateTx.Error)
			}
			if updateTx.RowsAffected != 1 {
				tx.Rollback()
				return errors.New("send prop failed")
			}
		} else { // 不存在 新增
			propNum = AiMatchProductBolingPropNumModel{
				BolingUuid: uuid.New().String(),
				UserId:     int64(req.UserId),
				PropId:     int64(req.PropId),
				Number:     int64(req.Number),
				CreateTime: req.Tm,
				UpdateTime: req.Tm,
				IsDelete:   0,
			}
			if err := tx.Create(&propNum).Error; err != nil {
				tx.Rollback()
				return errors.WithStack(err)
			}
		}
		// 新增记录
		record := &BolingPropRecordsModel{
			Source:     req.PropSource,
			UserId:     req.UserId,
			Num:        uint64(req.Number),
			RemainNum:  uint64(propNum.Number),
			Remark:     req.Remark,
			ChangeType: 1,
			CreateDate: req.Tm.Format(util.StandardFormat),
			PropId:     req.PropId,
		}
		if err := tx.Create(record).Error; err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}

		return tx.Commit().Error
	} else if strings.EqualFold(prop.BenefitType, "random") || strings.EqualFold(prop.BenefitType, "lottery") {
		otherProp := AiMatchProductOtherPropNumModel{}
		if err := tx.Model(&otherProp).Where("user_id = ? AND prop_id = ? AND is_delete = ?", req.UserId, req.PropId, 0).First(&otherProp).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return errors.WithStack(err)
		}
		if otherProp.OtherUuid != "" {
			updateTx := tx.Model(&otherProp).Where("other_uuid = ?", otherProp.OtherUuid).Updates(map[string]interface{}{
				"number": gorm.Expr("number + ?", req.Number),
			})
			if updateTx.Error != nil {
				tx.Rollback()
				return errors.WithStack(updateTx.Error)
			}
			if updateTx.RowsAffected != 1 {
				tx.Rollback()
				return errors.New("update other prop number failed!")
			}
		} else {
			otherProp = AiMatchProductOtherPropNumModel{
				OtherUuid:   strings.ReplaceAll(uuid.New().String(), "-", ""),
				BenefitType: prop.BenefitType,
				UserId:      req.UserId,
				PropId:      int64(req.PropId),
				Number:      req.Number,
				CreateTime:  req.Tm,
				UpdateTime:  req.Tm,
			}
			if err := tx.Create(&otherProp).Error; err != nil {
				tx.Rollback()
				return errors.WithStack(err)
			}
		}
		req.OutPropUserUuid = otherProp.OtherUuid
		return tx.Commit().Error
	}

	var (
		endTime   int64
		startTime = req.Tm.UnixMilli()
	)

	if strings.EqualFold(prop.BenefitType, "first_buy") {
		productSeckill := ProductBuySeckillModel{}
		if err := tx.Model(&productSeckill).Where("id = ? AND is_delete = ?", prop.ActiveId, 0).First(&productSeckill).Error; err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}
		startTime = productSeckill.StartTime - int64(prop.PriorityHour*60*60*1000)
		endTime = productSeckill.StartTime
	} else if strings.EqualFold(prop.BenefitType, "first_combine") {
		combine := AiMatchProductNftCombinationModel{}
		if err := tx.Model(&combine).Where("id = ? AND on_sale_status = ? AND is_delete = ?", prop.ActiveId, 1, 0).First(&combine).Error; err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}
		// startTime = combine.StartTime - int64(prop.PriorityHour*60*60*1000)
		// endTime = combine.StartTime
		startTime = int64(prop.StartTime)
		endTime = int64(prop.EndTime)
	} else {
		if prop.AvailableDay > 0 {
			endTime = startTime + 1000*60*60*24*int64(prop.AvailableDay)
		} else {
			startTime = int64(prop.StartTime)
			endTime = int64(prop.EndTime)
		}
	}

	if req.Source == "" {
		req.Source = "open_bind_box"
	}

	propUserList := make([]NftPropUserModel, 0, req.Number)
	for i := 0; i < int(req.Number); i++ {
		propUserList = append(propUserList, NftPropUserModel{
			PropUserUuid: strings.ReplaceAll(uuid.New().String(), "-", ""),
			UserId:       req.UserId,
			PropId:       req.PropId,
			StartTime:    uint64(startTime),
			EndTime:      uint64(endTime),
			Source:       req.Source,
			BenefitType:  prop.BenefitType,
		})
		if len(propUserList) == 500 {
			if err := tx.Create(&propUserList).Error; err != nil {
				tx.Rollback()
				return errors.WithStack(err)
			}
			propUserList = make([]NftPropUserModel, 0, req.Number)
		}
	}
	if len(propUserList) > 0 {
		if err := tx.Create(&propUserList).Error; err != nil {
			tx.Rollback()
			return errors.WithStack(err)
		}
	}

	return tx.Commit().Error
}
