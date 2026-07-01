package dto

import (
	"time"

	"hotbox-adm-backend/models/hd_task_models"
)

type GmvData struct {
	Dt       time.Time `json:"dt"`
	Gmv      float64   `json:"gmv"`
	UserCnt  uint      `json:"user_cnt"`
	RGmv     float64   `json:"r_gmv"`
	RUserCnt uint      `json:"r_user_cnt"`
}

type GmvResponse struct {
	List     []hd_task_models.CmDailyGmv `json:"list"`
	Total    int64                       `json:"total"`
	TotalGmv float64                     `json:"total_gmv"`
	Profit   float64                     `json:"profit"`
}

type NftCategoryGmvData struct {
	Dt           time.Time `json:"dt"`
	Gmv          float64   `json:"gmv"`
	UserCnt      uint      `json:"user_cnt"`
	CategoryPath string    `json:"category_path"`
	Category     string    `json:"category"`
	Rk           uint      `json:"rk"`
}
