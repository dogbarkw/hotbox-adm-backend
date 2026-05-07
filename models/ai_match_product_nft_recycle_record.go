package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/form"

	"github.com/sirupsen/logrus"
)

type AiMatchProductNftRecycleRecord struct {
	ID                 uint   `gorm:"primarykey" json:"id"`
	ProductTitle       string `gorm:"product_title" json:"product_title"`
	ProductId          int    `gorm:"product_id" json:"product_id"`
	NftProductSizeId   int    `gorm:"nft_product_size_id" json:"nft_product_size_id"`
	RecycleTargetCount int    `gorm:"recycle_target_count" json:"recycle_target_count"`
	RecycleCount       int    `gorm:"recycle_count" json:"recycle_count"`
	Status             int    `gorm:"status; comment: 0:待处理,1:在处理,-1:已完成" json:"status"`
	UserId             int    `gorm:"user_id" json:"user_id"`
	AdmUserName        string `gorm:"adm_user_name" json:"adm_user_name"`
	Msg                string `gorm:"msg" json:"msg"`
	Type               int    `gorm:"type; comment: 0:回收,1:空投"`
	OperatorId         int    `gorm:"operator_id; comment:后台操作者id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*AiMatchProductNftRecycleRecord) TableName() string {
	return "ai_match_product_nft_recycle_record"
}

func (AiMatchProductNftRecycleRecord) Create(ctx context.Context, m *AiMatchProductNftRecycleRecord) error {
	return cli.HotDogGormDB.WithContext(ctx).Create(&m).Error
}

func (AiMatchProductNftRecycleRecord) Update(ctx context.Context, id int, payload map[string]any) (int64, error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(AiMatchProductNftRecycleRecord{}).Where("id", id).Updates(payload)
	affectRow := query.RowsAffected
	if err := query.Error; err != nil {
		return 0, err
	}
	return affectRow, nil
}

func (AiMatchProductNftRecycleRecord) GetRecycleRecordList(ctx context.Context, req form.RecycleRecordListReq, order []string) (list []AiMatchProductNftRecycleRecord, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftRecycleRecord{})
	if req.Type != nil {
		query = query.Where("type", req.Type)
	}
	err = query.Count(&total).Error
	if err != nil {
		logrus.Errorf("GetBlackListList query count error: %v", err)
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	if req.PageSize != 0 && req.PageNumber != 0 {
		query.Offset((int(req.PageNumber) - 1) * int(req.PageSize)).Limit(int(req.PageSize))
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}
	return
}

func (AiMatchProductNftRecycleRecord) GetByParams(ctx context.Context, where map[string]any, limit int) (list []AiMatchProductNftRecycleRecord, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftRecycleRecord{})
	for k, v := range where {
		query = query.Where(k, v)
	}
	if limit != 0 {
		query = query.Limit(limit)
	}
	err = query.Scan(&list).Error
	return
}

func (AiMatchProductNftRecycleRecord) One(ctx context.Context, id int) (data AiMatchProductNftRecycleRecord, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(&AiMatchProductNftRecycleRecord{}).Where("id", id).First(&data).Error

	return
}
