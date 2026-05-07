package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var ActivityUpgradeDal = &ActivityUpgrade{}

// ActivityUpgrade 升级活动
type ActivityUpgrade struct {
	ID                      int64     `json:"id" gorm:"id"`                                           // 主键ID
	Title                   string    `json:"title" gorm:"title"`                                     // 标题
	SubTitle                string    `json:"sub_title" gorm:"sub_title"`                             // 副标题
	Type                    int8      `json:"type" gorm:"type"`                                       // 合成类型 1=藏品
	Image                   string    `json:"image" gorm:"image"`                                     // 活动封面图
	ProductId               int64     `json:"product_id" gorm:"product_id"`                           // 藏品ID
	ProductSizeId           int64     `json:"product_size_id" gorm:"product_size_id"`                 // 藏品规格ID
	BeginTime               time.Time `json:"begin_time" gorm:"begin_time"`                           // 开始时间
	EndTime                 time.Time `json:"end_time" gorm:"end_time"`                               // 结束时间
	SaleTime                time.Time `json:"sale_time" gorm:"sale_time"`                             // 二级发售时间
	Point                   int64     `json:"point" gorm:"point"`                                     // 升级奖励积分
	Weight                  int64     `json:"weight" gorm:"weight"`                                   // 权重
	LowestVersion           int64     `json:"lowest_version" gorm:"lowest_version"`                   // 最低版本
	FailTreasureBoxId       int64     `json:"fail_treasure_box_id" gorm:"fail_treasure_box_id"`       // 失败奖励百宝箱ID
	DestroyTreasureBoxId    int64     `json:"destroy_treasure_box_id" gorm:"destroy_treasure_box_id"` // 摧毁奖励百宝箱ID
	Status                  int8      `json:"status" gorm:"status"`                                   // 状态 1=下架 2=正常 3=删除
	SuccessRate             int16     `json:"success_rate" gorm:"success_rate"`                       // 升级成功概率
	FailRate                int16     `json:"fail_rate" gorm:"fail_rate"`                             // 升级失败概率
	DestroyRate             int16     `json:"destroy_rate" gorm:"destroy_rate"`                       // 升级摧毁概率
	Hot                     int64     `json:"hot" gorm:"hot"`
	CreatedAt               time.Time `json:"created_at" gorm:"created_at"`                                   // 创建时间
	UpdatedAt               time.Time `json:"updated_at" gorm:"updated_at"`                                   // 更新时间
	DestroyFailNum          int16     `json:"destroy_fail_num" gorm:"destroy_fail_num"`                       // 升级催毁失败次数
	UpgradeFailNum          int16     `json:"upgrade_fail_num" gorm:"upgrade_fail_num"`                       // 升级失败次数
	UnsuccessNum            int16     `json:"unsuccess_num" gorm:"unsuccess_num"`                             // 升级不成功次数
	IsDisplayCount          int16     `json:"is_display_count" gorm:"is_display_count"`                       // 是否展示活动数量
	IsDisplayTime           int16     `json:"is_display_time" gorm:"is_display_time"`                         // 是否展示活动时间
	PropId                  int64     `json:"prop_id" gorm:"prop_id"`                                         // 道具id
	SuccessMinNum           int64     `json:"success_min_num" gorm:"success_min_num"`                         // 升级 N次必中最小值
	SuccessMaxNum           int64     `json:"success_max_num" gorm:"success_max_num"`                         // 升级 N次必中最大值
	FailDestroyRate         int64     `json:"fail_destroy_rate" gorm:"fail_destroy_rate"`                     // 升级失败摧毁概率
	APoolMaxNum             int64     `json:"a_pool_max_num" gorm:"a_pool_max_num"`                           // a池可升级成功次数
	APoolSuccessMinNum      int64     `json:"a_pool_success_min_num" gorm:"a_pool_success_min_num"`           // a池用户 N次必中最小值
	APoolSuccessMaxNum      int64     `json:"a_pool_success_max_num" gorm:"a_pool_success_max_num"`           // a 池用户 N次必中最大值
	APoolSuccessNum         int64     `json:"a_pool_success_num" gorm:"a_pool_success_num"`                   // a池已升级成功次数
	APoolMutationTotalNum   int64     `json:"a_pool_mutation_total_num" gorm:"a_pool_mutation_total_num"`     // a池突变总次数
	APoolMutationMinNum     int64     `json:"a_pool_mutation_min_num" gorm:"a_pool_mutation_min_num"`         // a池用户 N次突变最小值
	APoolMutationMaxNum     int64     `json:"a_pool_mutation_max_num" gorm:"a_pool_mutation_max_num"`         // a池用户 N次突变最大值
	APoolMutationSuccessNum int64     `json:"a_pool_mutation_success_num" gorm:"a_pool_mutation_success_num"` // a池已突变成功次数
	ReserveNum              int64     `json:"reserve_num" gorm:"reserve_num"`                                 // 预留量
	Channel                 int64     `json:"channel" gorm:"channel"`                                         // 通道顺序 【需求：运营活动列表展示自动排序】
	TotalNum                int64     `json:"total_num" gorm:"total_num"`                                     // 活动总参与次数
	FinishedNum             int64     `json:"finished_num" gorm:"finished_num"`                               // 已参与次数
	MutationTotalNum        int64     `json:"mutation_total_num" gorm:"mutation_total_num"`                   // 突变总数量
	MutationFinishedNum     int64     `json:"mutation_finished_num" gorm:"mutation_finished_num"`             // 突变已完成数量
	MutationReserveNum      int64     `json:"mutation_reserve_num" gorm:"mutation_reserve_num"`               // 突变预留数量
}

// TableName 表名称
func (*ActivityUpgrade) TableName() string {
	return "activity_upgrade"
}

func (a *ActivityUpgrade) GetByParams(ctx context.Context, where map[string]any) (list []ActivityUpgrade, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (a *ActivityUpgrade) One(ctx context.Context, id int) (resp ActivityUpgrade, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("id", id).First(&resp).Error
	return
}

func (a *ActivityUpgrade) GetActivityUpgradeByParams(ctx context.Context, where map[string]any, order []string, limit, offset int) (list []ActivityUpgrade, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&ActivityUpgrade{})
	for k, v := range where {
		query.Where(k, v)
	}
	if limit > 0 {
		query.Limit(limit)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}
