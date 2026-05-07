package models

import (
	"context"
	"time"

	"gorm.io/gorm/clause"

	"hotbox-adm-backend/form"

	"hotbox-adm-backend/cli"

	"gorm.io/gorm"
)

var ActivityScoreDal = &ActivityScore{}

// 活动打分
type ActivityScore struct {
	Id                      int            `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	ActivityId              int64          `gorm:"column:activity_id;type:int(11);NOT NULL" json:"activity_id"`
	ActivityType            int            `gorm:"column:activity_type;type:tinyint(3);NOT NULL" json:"activity_type"`
	ActivityStartTs         int64          `gorm:"column:activity_start_ts;type:bigint(20);default:0;comment:活动开始时间;NOT NULL" json:"activity_start_ts"`
	ActivityTitle           string         `gorm:"column:activity_title;type:varchar(255)" json:"activity_title"`
	TotalCost               float64        `gorm:"column:total_cost;type:decimal(10,2);default:0.00;comment:总成本;NOT NULL" json:"total_cost"`
	ActivityDuration        int64          `gorm:"column:activity_duration;type:int(11);default:0;comment:活动持续时长;NOT NULL" json:"activity_duration"`
	WarnTime                int            `gorm:"column:warn_time;type:int(11);default:0;comment:报警次数;NOT NULL" json:"warn_time"`
	NotchNumber             int64          `gorm:"column:notch_number;type:int(11);default:0;comment:缺口数;NOT NULL" json:"notch_number"`
	ExpectedProductAmount   int64          `gorm:"column:expected_product_amount;type:int(11);default:0;comment:预期流通份数;NOT NULL" json:"expected_product_amount"`
	RecommendScoreBeforeEnd float64        `gorm:"column:recommend_score_before_end;type:decimal(10,2);default:0.00;comment:活动结束前的推荐分;NOT NULL" json:"recommend_score_before_end"`
	PraiseScore             float64        `gorm:"column:praise_score;type:decimal(10,2);default:0.00;comment:活动结束口碑分;NOT NULL" json:"praise_score"`
	Status                  int            `gorm:"column:status;type:tinyint(2);comment:状态:0未结束；1：已结束;NOT NULL" json:"status"`
	CreatedAt               time.Time      `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time      `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名称
func (*ActivityScore) TableName() string {
	return "activity_score"
}

func (a *ActivityScore) GetByParams(ctx context.Context, where map[string]any, order []string) (list []ActivityScore, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}

func (a *ActivityScore) GetEndActivityScorePage(ctx context.Context, req form.GetEndedActivityScoreListReq) (list []ActivityScore, total int64, err error) {
	order := []string{"activity_start_ts DESC"}
	if req.OrderType == 2 {
		order = []string{"warn_time desc"}
	}
	where := map[string]any{
		"activity_type": req.ActivityType,
		"status":        1,
	}
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}
	if err := query.Count(&total).Error; err != nil {
		return list, 0, err
	}

	if req.PageSize != 0 && req.PageNumber != 0 {
		query.Offset((int(req.PageNumber) - 1) * req.PageSize).Limit(req.PageSize)
	}

	for _, v := range order {
		query.Order(v)
	}
	err = query.Scan(&list).Error
	return
}

func (a *ActivityScore) One(ctx context.Context, id int) (resp ActivityScore, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("id", id).First(&resp).Error
	return
}

func (a *ActivityScore) Update(ctx context.Context, id int, payload map[string]any) (err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("id", id).Updates(payload).Error
	return
}

func (*ActivityScore) Create(ctx context.Context, dm *ActivityScore) (err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Save(&dm).Error
	return err
}

func (*ActivityScore) FirstOrCreate(ctx context.Context, dm ActivityScore) (affectRow int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Where("activity_id", dm.ActivityId).Where("activity_type", dm.ActivityType).Clauses(clause.Insert{Modifier: "IGNORE"}).FirstOrCreate(&dm)
	err = query.Error
	if err != nil {
		return 0, err
	}
	return query.RowsAffected, nil
}

func (*ActivityScore) Delete(ctx context.Context, id uint, force bool) error {
	query := cli.HotDogGormDB.WithContext(ctx).Model(&ActivityScore{})
	if force {
		query = query.Unscoped()
	}
	return query.Delete(&ActivityScore{}, "id", id).Error
}

func (a *ActivityScore) ActivityInitAndUpdate(ctx context.Context, dm ActivityScore) (affectRow int64, err error) {
	affectRow, err = a.FirstOrCreate(ctx, dm)
	if err != nil {
		return
	}
	if affectRow == 1 {
		return
	}
	data, err := a.GetByUniqKey(ctx, int(dm.ActivityId), dm.ActivityType)
	if err != nil {
		return 0, err
	}
	data.ActivityTitle = dm.ActivityTitle
	data.ActivityStartTs = dm.ActivityStartTs
	if data.NotchNumber == 0 && data.TotalCost == 0 && data.ExpectedProductAmount == 0 { // 保证只快照一次
		data.TotalCost = dm.TotalCost
		data.NotchNumber = dm.NotchNumber
		data.ExpectedProductAmount = dm.ExpectedProductAmount
	}
	query := cli.HotDogGormDB.WithContext(ctx).Save(&data)
	err = query.Error
	if err != nil {
		return 0, err
	}
	return query.RowsAffected, nil
}

func (a *ActivityScore) Save(ctx context.Context, dm ActivityScore) error {
	return cli.HotDogGormDB.WithContext(ctx).Save(&dm).Error
}

func (a *ActivityScore) GetByUniqKey(ctx context.Context, activityId int, activityType int) (data ActivityScore, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Where("activity_id = ? AND activity_type = ?", activityId, activityType).First(&data).Error
	return
}
