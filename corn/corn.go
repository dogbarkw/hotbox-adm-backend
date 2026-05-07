package corn

import (
	"hotbox-adm-backend/corn/activity_score_nft_second_price"
	"hotbox-adm-backend/corn/airdrop_record"
	"hotbox-adm-backend/corn/artist_nft_activity_corn"
	"hotbox-adm-backend/corn/daily_gmv"
	"hotbox-adm-backend/corn/dg_yop_test_user"
	"hotbox-adm-backend/corn/nft_reserved"
	"hotbox-adm-backend/corn/recycle_record"
)

func NewActivityScoreInitCornJob() *activity_score_nft_second_price.ActivityScoreInitCorn {
	return &activity_score_nft_second_price.ActivityScoreInitCorn{}
}

func NewAirdropRecordCornJob() *airdrop_record.AirdropRecordCornJob {
	return &airdrop_record.AirdropRecordCornJob{}
}

func NewArtistNftActivityCountCornJob() *artist_nft_activity_corn.ArtistNftActivityCountCornJob {
	return &artist_nft_activity_corn.ArtistNftActivityCountCornJob{}
}

func NewNftReservedCornJob() *nft_reserved.NftReservedCornJob {
	return &nft_reserved.NftReservedCornJob{}
}

func NewRecycleRecordCornJob() *recycle_record.RecycleRecordCornJob {
	return &recycle_record.RecycleRecordCornJob{}
}

func NewDailyGmvCornJob() *daily_gmv.DailyGmvCornJob {
	return &daily_gmv.DailyGmvCornJob{}
}

func NewPreDailyGmvCornJob() *daily_gmv.DailyBeforeGmvCornJob {
	return &daily_gmv.DailyBeforeGmvCornJob{}
}

func NewYopTestUserJob() *dg_yop_test_user.DgYopTestUserIncomeCornJob {
	return &dg_yop_test_user.DgYopTestUserIncomeCornJob{}
}
