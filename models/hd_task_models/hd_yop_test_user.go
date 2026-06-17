package hd_task_models

import (
	"context"
	"database/sql"
	"time"

	"hotbox-adm-backend/cli"
)

var HdYopTestUserDal = &HdYopTestUser{}

// 特殊账号表
type HdYopTestUser struct {
	Id          int64        `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	UserId      int64        `gorm:"column:user_id;type:bigint(20);default:0;comment:平台用户ID;NOT NULL" json:"user_id"`
	Mobile      string       `gorm:"column:mobile;type:varchar(32);comment:手机号;NOT NULL" json:"mobile"`
	RealName    string       `gorm:"column:real_name;type:varchar(32);comment:帐号实名;NOT NULL" json:"real_name"`
	UserType    int          `gorm:"column:user_type;type:tinyint(4);default:0;comment:账号类型 1实名账号 2测试账号;NOT NULL" json:"user_type"`
	Rate        int          `gorm:"column:rate;type:int(11);default:0;comment:分成比例;NOT NULL" json:"rate"`
	FreezeRate  int          `gorm:"column:freeze_rate;type:int(11);default:0;comment:到账冻结比例;NOT NULL" json:"freeze_rate"`
	TotalIncome float64      `gorm:"column:total_income;type:decimal(10,2);default:0.00;comment:累计进账;NOT NULL" json:"total_income"`
	Remark      string       `gorm:"column:remark;type:varchar(255);comment:备注;NOT NULL" json:"remark"`
	CountTime   sql.NullTime `gorm:"column:count_time;type:datetime;comment:上次进账统计时间" json:"count_time"`
	CreatedAt   time.Time    `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdYopTestUser) TableName() string {
	return "hotbox_yop_test_user"
}

func (m *HdYopTestUser) FirstOrCreate(ctx context.Context, payload *HdYopTestUser) (int64, error) {
	tx := cli.HotDogTaskGormDB.WithContext(ctx).Where("user_id = ?", payload.UserId).FirstOrCreate(&payload)
	return tx.RowsAffected, tx.Error
}

func (m *HdYopTestUser) UpdateByParams(ctx context.Context, where map[string]any, params map[string]any) error {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&m)
	for k, v := range where {
		query.Where(k, v)
	}
	return query.Updates(params).Error
}

func (m *HdYopTestUser) GetList(ctx context.Context, where map[string][]any, order []string, page, size int) (list []HdYopTestUser, count int64, err error) {
	offset := (page - 1) * size
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&HdYopTestUser{})
	for k, v := range where {
		query.Where(k, v...)
	}
	for _, v := range order {
		query.Order(v)
	}
	if order != nil {
		query.Order(order)
	}
	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	if count == 0 {
		return nil, 0, nil
	}
	if size > 0 {
		query = query.Limit(size)
	}
	if offset > 0 {
		query.Offset(offset)
	}
	err = query.Scan(&list).Error
	return
}

func (m *HdYopTestUser) GetHdYopTestUsers(ctx context.Context, where map[string][]any, order []string, limit int) (list []HdYopTestUser, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&HdYopTestUser{})
	for k, v := range where {
		query.Where(k, v...)
	}
	for _, v := range order {
		query.Order(v)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err = query.Scan(&list).Error
	return
}

func (m *HdYopTestUser) Delete(ctx context.Context, ids []int64) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).Unscoped().Delete(&m, "id", ids).Error
}

func (m *HdYopTestUser) SumTotalIncome(ctx context.Context) (total float64, err error) {
	query := cli.HotDogTaskGormDB.WithContext(ctx).Model(&HdYopTestUser{})
	query.Select("IFNULL(SUM(total_income),0) AS total")
	err = query.First(&total).Error
	return
}

func (m *HdYopTestUser) One(ctx context.Context, id int64) (data HdYopTestUser, err error) {
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).Where("id", id).First(&data).Error
	return
}
