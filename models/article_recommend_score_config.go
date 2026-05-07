package models

import (
	"context"
	"time"

	"hotbox-adm-backend/cli"
)

var ArtistRecommendScoreConfigDal = &ArtistRecommendScoreConfig{}

// 艺术家推荐分设置
type ArtistRecommendScoreConfig struct {
	Id          int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	DeletedAt   time.Time `gorm:"column:deleted_at;type:datetime" json:"deleted_at"`
	ArtistId    int64     `gorm:"column:artist_id;type:bigint(20);default:0;comment:艺术家ID;NOT NULL" json:"artist_id"`
	ArtistName  string    `gorm:"column:artist_name;type:varchar(128);comment:艺术家名称;NOT NULL" json:"artist_name"`
	NftNum      int       `gorm:"column:nft_num;type:int(11);default:0;comment:藏品数;NOT NULL" json:"nft_num"`
	ActivityNum int       `gorm:"column:activity_num;type:int(11);default:0;comment:活动数;NOT NULL" json:"activity_num"`
	Score       string    `gorm:"column:score;type:varchar(128);default:100;comment:口碑分,初始值100;NOT NULL" json:"score"`
}

func (a *ArtistRecommendScoreConfig) TableName() string {
	return "artist_recommend_score_config"
}

func (a *ArtistRecommendScoreConfig) GetByParams(ctx context.Context, where map[string]any) (list []ArtistRecommendScoreConfig, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query = query.Where(k, v)
	}

	err = query.Scan(&list).Error
	return
}

func (a *ArtistRecommendScoreConfig) GetArtistRecommendScoreConfigList(ctx context.Context, where map[string]any, order []string, pageNum, pageSize int) (list []*ArtistRecommendScoreConfig, count int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range where {
		query.Where(k, v)
	}
	err = query.Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		return
	}
	for _, v := range order {
		query.Order(v)
	}
	err = query.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	return
}

func (a *ArtistRecommendScoreConfig) One(ctx context.Context, id int) (resp ArtistRecommendScoreConfig, err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Model(a).Where("id", id).First(&resp).Error
	return
}

func (a *ArtistRecommendScoreConfig) GetOneByParams(ctx context.Context, params map[string]any) (resp ArtistRecommendScoreConfig, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Model(a)
	for k, v := range params {
		query.Where(k, v)
	}
	err = query.First(&resp).Error
	return
}

func (a *ArtistRecommendScoreConfig) UpdateParamsByArtistId(ctx context.Context, artistId int64, params map[string]any) error {
	return cli.HotDogGormDB.WithContext(ctx).Model(a).Where("artist_id", artistId).Updates(params).Error
}

func (a *ArtistRecommendScoreConfig) UpdateParamsById(ctx context.Context, id int64, params map[string]any) error {
	return cli.HotDogGormDB.WithContext(ctx).Model(a).Where("id", id).Updates(params).Error
}

func (a *ArtistRecommendScoreConfig) Save(ctx context.Context, dm *ArtistRecommendScoreConfig) (err error) {
	err = cli.HotDogGormDB.WithContext(ctx).Save(&dm).Error
	return err
}

func (*ArtistRecommendScoreConfig) FirstOrCreate(ctx context.Context, dm *ArtistRecommendScoreConfig) (affectRow int64, err error) {
	query := cli.HotDogGormDB.WithContext(ctx).Where("artist_id", dm.ArtistId).FirstOrCreate(&dm)
	err = query.Error
	if err != nil {
		return 0, err
	}
	return query.RowsAffected, nil
}

func (a *ArtistRecommendScoreConfig) SaveNftNum(ctx context.Context, dm *ArtistRecommendScoreConfig) (err error) {
	affectRow, err := a.FirstOrCreate(ctx, dm)
	if err != nil {
		return
	}
	if affectRow == 1 {
		return
	}
	query := cli.HotDogGormDB.WithContext(ctx).Save(&dm)
	err = query.Error
	if err != nil {
		return err
	}
	return nil
}
