package route

import (
	"os"

	"hotbox-adm-backend/api"
	"hotbox-adm-backend/until"

	"github.com/gin-gonic/gin"
)

func RegisterRouters(r *gin.Engine) {
	r.Use(until.CORS())
	poker := r.Group("/hotbox", until.LoginChk())
	if os.Getenv("ENV") == "dev" {
		// 本地环境不验token
		poker = r.Group("/hotbox")
	}
	operation := poker.Group("/operation")
	{
		operation.POST("/nft_second/list", api.NftSecondPriceList)
		// 刷新普通用户剩余份数
		operation.POST("/nft/flush/surplus", api.NftSecondPriceFlushSurplus)
		// 获取 普通用户剩余份数 和 用户百分比
		operation.GET("/nft/flush/surplus", api.GetNftSecondPriceFlushSurplus)
		// 刷新用户百分比
		operation.POST("/nft/flush/user_percentage", api.NftSecondPriceFlushUserPercentage)
		// 按数量回收
		operation.POST("/nft/recycling/by_count", api.NftRecyclingByCount)
		// 按数量摧毁
		operation.POST("/nft/destroy/by_count", api.NftDestroyByCount)
		// 按数量空投
		operation.POST("/nft/airdrop/by_count", api.NftAirdropByCount)
		// 按数量摧毁（空投）记录
		operation.POST("/nft/recycle/record/list", api.GetRecycleRecordList)
		operation.POST("/nft/recycle/record/:id", api.GetRecycleRecordById)

		// NFT未出售列表
		operation.POST("/nft/rest/list", api.NftRestList)
		operation.POST("/nft/product/simple_list", api.GetSimpleList)
		// 给制定用户新增空投道具
		operation.POST("/airdrop/treasure", api.AirdropTreasure)
		// 置换
		displace := operation.Group("/nft/displace")
		{
			displace.POST("/list", api.GetDisplaceList)
			displace.PUT("/update/replace_count/:replaceId", api.UpdateNftReplaceDisplaceCount)
			displace.PUT("/update/reserve_num/:replaceId", api.UpdateNftReplaceReserveNum)
		}
		// 合成
		combination := operation.Group("/nft/combination")
		{
			combination.POST("/list", api.GetNftCombinationList)
		}
		// 活动材料预留
		activityMaterialReserve := operation.Group("/nft")
		{
			// 合成活动
			{
				// 获取合成活动预留材料任务信息
				activityMaterialReserve.POST("/combination/material_reserve/query", api.GetNftCombinationMaterialReserveTask)
				// 新建/更新合成活动预留材料任务信息
				activityMaterialReserve.POST("/combination/material_reserve/create", api.UpdateNftCombinationMaterialReserveTask)
				// 删除合成活动预留材料任务信息
				activityMaterialReserve.POST("/combination/material_reserve/delete", api.DeleteNftCombinationMaterialReserveTask)
				// 计算合成活动预留材料任务信息合成次数
				activityMaterialReserve.POST("/combination/material_reserve/calculate", api.CalculateNftCombinationMaterialReserveTask)
			}
			// 置换、分解
			{
				// 获取合成活动预留材料任务信息
				activityMaterialReserve.POST("/displace/material_reserve/query", api.GetNftDisplaceMaterialReserveTask)
				// 新建/更新合成活动预留材料任务信息
				activityMaterialReserve.POST("/displace/material_reserve/create", api.UpdateNftDisplaceMaterialReserveTask)
				// 删除合成活动预留材料任务信息
				activityMaterialReserve.POST("/displace/material_reserve/delete", api.DeleteNftDisplaceMaterialReserveTask)
				// 计算合成活动预留材料任务信息合成次数
				activityMaterialReserve.POST("/displace/material_reserve/calculate", api.CalculateNftDisplaceMaterialReserveTask)
			}

		}
	}

	operation = r.Group("/hotbox/v2/operation", until.GTokenChk())
	{
		// gmv列表
		operation.POST("/gmv/list", api.GmvList)
	}
}
