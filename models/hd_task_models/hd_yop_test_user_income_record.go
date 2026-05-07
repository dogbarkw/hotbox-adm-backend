package hd_task_models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var HdYopTestUserIncomeRecordDal = &HdYopTestUserIncomeRecord{}

// 特殊账号进账记录表
type HdYopTestUserIncomeRecord struct {
	Id            uint64     `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT;comment:主键 ID" json:"id"`
	YopTestUserId int64      `gorm:"column:yop_test_user_id;type:bigint(20);default:0;comment:特殊账号ID;NOT NULL" json:"yop_test_user_id"`
	Fee           float64    `gorm:"column:fee;type:decimal(10,2);default:0.00;comment:金额;NOT NULL" json:"fee"`
	Rate          int        `gorm:"column:rate;type:int(11);default:0;comment:当前分成比例;NOT NULL" json:"rate"`
	Income        float64    `gorm:"column:income;type:decimal(10,2);default:0.00;comment:进账;NOT NULL" json:"income"`
	IncomeTime    *time.Time `gorm:"column:income_time;type:datetime;comment:进账时间" json:"income_time"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (m *HdYopTestUserIncomeRecord) TableName() string {
	return "hotbox_yop_test_user_income_record"
}

func (m *HdYopTestUserIncomeRecord) BatchCreate(ctx context.Context, payload []*HdYopTestUserIncomeRecord) error {
	return cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).CreateInBatches(&payload, 30).Error
}

type UserDailyIncome struct {
	YopTestUserId int64   `json:"yop_test_user_id"`
	Total         float64 `json:"total"`
}

func (m *HdYopTestUserIncomeRecord) SumIncomeByTimeRange(ctx context.Context, testUserIds []int64, startTime, endTime time.Time) (userIncomeMap map[int64]float64, err error) {
	userIncomeMap = make(map[int64]float64, len(testUserIds))
	result := make([]UserDailyIncome, 0, len(testUserIds))
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).
		Where("income_time >= ?", startTime).
		Where("income_time <= ?", endTime).
		Where("yop_test_user_id in ?", testUserIds).
		Select([]string{"yop_test_user_id", "IFNULL(SUM(income),0) as total"}).
		Group("yop_test_user_id").
		Scan(&result).Error
	if err != nil {
		return userIncomeMap, err
	}
	for _, v := range result {
		userIncomeMap[v.YopTestUserId] = v.Total
	}
	return
}

func (m *HdYopTestUserIncomeRecord) SumTodayIncome(ctx context.Context, start, end time.Time) (total float64, err error) {
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).
		Where("income_time >= ?", start).
		Where("income_time <= ?", end).
		Select("IFNULL(SUM(income),0) as total").
		First(&total).Error
	return
}

type YopTestUserIncomeStats struct {
	Month       string `json:"month"`        // 月份
	TotalIncome string `json:"total_income"` // 总收入
}
type YopTestUserIncomeDetail struct {
	CreatedAt time.Time `json:"created_at"` // 月份
	Mobile    string    `json:"mobile"`     // 手机号
	RealName  string    `json:"real_name"`  // 真实姓名
	Remark    string    `json:"remark"`     // 备注
	Rate      int       `json:"rate"`       // 分成比例
	Fee       string    `json:"fee"`        // 金额
}

func (m *HdYopTestUserIncomeRecord) StatByMouth(ctx context.Context) (result []YopTestUserIncomeStats, err error) {
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).
		Select("DATE_FORMAT(income_time, '%Y-%m') AS month", "SUM(income) AS total_income").
		Group("DATE_FORMAT(income_time, '%Y-%m')").
		Order("month desc").
		Scan(&result).Error
	return
}

func (m *HdYopTestUserIncomeRecord) GetIncomeRecord(ctx context.Context, timeStart, timeEnd time.Time) (result []YopTestUserIncomeDetail, err error) {
	err = cli.HotDogTaskGormDB.WithContext(ctx).Model(&m).Raw(`SELECT
	hotbox_yop_test_user_income_record.income_time as created_at,
	hotbox_yop_test_user.mobile,
	hotbox_yop_test_user.real_name,
	hotbox_yop_test_user.remark,
	hotbox_yop_test_user_income_record.rate, 
	hotbox_yop_test_user_income_record.income as fee
FROM
	hotbox_yop_test_user_income_record
	LEFT JOIN hotbox_yop_test_user ON hotbox_yop_test_user.id=hotbox_yop_test_user_income_record.yop_test_user_id
WHERE
	income_time >= ? 
	AND income_time < ? 
ORDER BY
	income_time`, timeStart, timeEnd).Scan(&result).Error
	return
}
