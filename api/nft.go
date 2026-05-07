package api

import (
	"encoding/json"
	"fmt"
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/dto"
	"hotbox-adm-backend/form"
	"hotbox-adm-backend/internal/httpReq"
	"hotbox-adm-backend/models"
	"hotbox-adm-backend/until"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"hotbox-adm-backend/pkg/constant"
	"hotbox-adm-backend/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const (
	PFP_VIDEO    = "PFP_video"
	PFP_PHOTO    = "PFP_photo"
	SYG_photo    = "SYG_photo"
	NFT_BOX_TYPE = "box"
)

var (
	PFP_LIST              = []string{PFP_VIDEO, PFP_PHOTO, SYG_photo}
	RecycleNoNftSizeError = errors.New("NFTAdmRecycling, no nft size")
)

func NftSecondPriceList(c *gin.Context) {
	req := form.SecondPriceListReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	res, total, err := models.Nft{Ctx: c}.NftProductPriceList(req)
	if err != nil {
		response.ResponseFail("mysql query fail")
		return
	}

	newResList := make([]dto.NftSecondPriceListRespItem, 0, len(res))
	nftSizeList := make([]uint64, 0, len(res))
	for _, v := range res {
		nftSizeList = append(nftSizeList, v.NftProductSizeId)
	}
	sizeInfoList, err := models.NewSaleProductNftSize().NftProductListByCondition(c, models.SaleProductNftSizeListByConditionReq{
		Ids: lo.Uniq(nftSizeList),
	})
	if err != nil {
		response.ResponseFail(fmt.Sprintf("查询NftProductListByCondition出错 err: %v", err))
		return
	}
	sizeInfoMap := make(map[int64]*models.SaleProductNftSizeModel)
	for _, v := range sizeInfoList {
		sizeInfoMap[v.NftProductSizeId] = v
	}
	for i := range res {
		v := res[i]
		if v == nil {
			continue
		}
		// 盲盒
		if v.NftType == NFT_BOX_TYPE && v.NftProductSizeId > 0 {
			// 获取盲盒内藏品名字
			if _, ok := sizeInfoMap[int64(v.NftProductSizeId)]; ok {
				v.ProductTitle = sizeInfoMap[int64(v.NftProductSizeId)].ProductTitle
			}
		}
		newRes := dto.NftSecondPriceListRespItem{
			ProductNftSecondPriceJoinSaleCalendar: *v,
		}
		// 查询藏品预留数量
		taskList := make([]dto.ProductNftSecondPriceJoinSaleCalendarReserveTask, 0)
		if v.CountResetTime != "0" && v.CountResetValue != 0 {
			taskList = append(taskList, dto.ProductNftSecondPriceJoinSaleCalendarReserveTask{
				ExecTime:   cast.ToInt64(v.CountResetTime),
				ReserveNum: int64(v.CountResetValue),
			})
		}
		list, tempErr := models.ActivityMaterialReserveDetailDal.GetJoinMainWithParams(map[string]any{
			"product_id":      v.ProductId,
			"product_size_id": v.NftProductSizeId,
			"exec_status":     0,
			"status":          0,
		})
		if tempErr != nil {
			newRes.TaskList = taskList
			newResList = append(newResList, newRes)
			_ = httpReq.FeiShuDebugRootBot(fmt.Sprintf("[NftSecondPriceList]，GetJoinMainWithParams fail, product_id:%d, product_size_id:%d, err%v", v.ProductId, v.NftProductSizeId, tempErr))
			continue
		}
		for _, d := range list {
			taskList = append(taskList, dto.ProductNftSecondPriceJoinSaleCalendarReserveTask{
				ExecTime:   d.ExecTime,
				ReserveNum: d.ReserveNum,
			})
		}
		newRes.TaskList = taskList
		newResList = append(newResList, newRes)
	}

	c.JSON(200, gin.H{
		"data":        newResList,
		"total_count": total,
		"code":        0,
		"msg":         errno.GetMsg(errno.Success),
		"is_list":     1,
	})
}

func NftSecondPriceFlushSurplus(c *gin.Context) {
	req := form.SecondPriceFlushSurplusReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	userId, _ := c.Get("user_id")
	flushKey := fmt.Sprintf("NFTFlushSurplusAdm:%d:%d:%d-%d", userId, req.FlushType, req.ProductId, req.NftProductSizeId)

	st := cli.HotDogRedis.SetNX(c, flushKey, "lock", 2*time.Second).Val()
	if !st {
		response.Responses(10112, "", nil)
		return
	}
	code, total, err := models.Nft{Ctx: c}.NftSecondPriceFlushSurplus(req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.Responses(code, errno.GetMsg(errno.Success), gin.H{
		"number": total,
	})
}

func GetNftSecondPriceFlushSurplus(c *gin.Context) {
	req := form.GetSecondPriceFlushSurplusReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.ValidatorQuery(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	realUserSurplus, err := models.Nft{Ctx: c}.GetNftSecondPriceFlushSurplus(req)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	nftProductSizeId := req.NftProductSizeId
	nftCountData, err := models.NewBusinessNftMarketWarehouseTotalCount(c).GetByProductIdAndSizeId(int64(req.ProductId), int64(*nftProductSizeId))
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	nftCount := nftCountData.NftCount
	var percent float64
	if nftCount != 0 {
		percent = float64(realUserSurplus) / float64(nftCount)
	}
	response.Responses(200, errno.GetMsg(errno.Success), gin.H{
		"number":  realUserSurplus,
		"percent": percent,
	})
}

func NftRecyclingByCount(c *gin.Context) {
	req := form.NftRecyclingByCountReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	userId, _ := c.Get("user_id")
	admUser, err := models.User{Ctx: c}.FindAdmUserById(cast.ToUint64(userId))
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 防止频繁攻击
	propBoxOpenKey := "NFTAdmRecycling"
	st := cli.HotDogRedis.SetNX(c, propBoxOpenKey, "lock", 6*time.Second).Val()
	if !st {
		response.Responses(10112, "", nil)
		return
	}
	user, err := models.User{Ctx: c}.FindSysUserByMobile(req.Mobile)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 获取相应的订单
	orderList, err := models.AiMatchProductOrder{Ctx: c}.GetProductOrderList(map[string]any{
		"receiver_name":       cast.ToString(user.UserId),
		"product_type":        "NFT",
		"status":              2,
		"is_delete":           0,
		"apply_refund_status": 0,
		"nft_product_size_id": req.NftProductSizeID,
		"product_id":          req.ProductID,
	}, &req.Count)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	if len(orderList) != req.Count {
		response.Responses(10112, "数量不足", nil)
		return
	}
	for _, v := range orderList {
		if v.MarketType == "land" {
			response.ResponseFails("暂不支持回收土地", 1008)
			return

		} else if lo.Contains[string]([]string{"***", ""}, v.ReceiverCity) {
			response.ResponseFails("订单还没发编号，不能回收", 1008)
			return
		}
	}
	productId := req.ProductID
	productTitle := orderList[0].ProductTitle
	p := struct {
		Title        string      `json:"title"`
		Nums         int         `json:"nums"`
		UserInfoList [][2]string `json:"userInfoList"`
	}{
		Title: productTitle,
		Nums:  req.Count,
		UserInfoList: [][2]string{
			{cast.ToString(user.UserId), req.Mobile},
		},
	}
	jsonStr, _ := json.Marshal(p)

	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      int64(admUser.UserId),
		Username:    admUser.Name,
		Remark:      "回收藏品",
		Scenes:      67,
		AssociateId: int64(productId),
		RequestData: string(jsonStr),
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	recycleRecordData := models.AiMatchProductNftRecycleRecord{
		ProductTitle:       productTitle,
		ProductId:          productId,
		NftProductSizeId:   req.NftProductSizeID,
		RecycleTargetCount: req.Count,
		Status:             0,
		UserId:             int(user.UserId),
		AdmUserName:        admUser.Name,
		Type:               constant.BATCH_RECORD_RECYCLE_TYPE,
		OperatorId:         admUser.UserId,
	}
	err = models.AiMatchProductNftRecycleRecord{}.Create(c, &recycleRecordData)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 发送消息到飞书
	msgType := "批量回收"
	msg := fmt.Sprintf("%s(%d)  回收 %s (mobile: %s) 如下藏品: 《%s》%d 份", admUser.Name, admUser.UserId, user.RealName, user.Mobile, productTitle, req.Count)
	httpReq.FeiShuWithUrlRootBot(constant.FEI_SHU_RECYCLE_URL, msgType, msg)
	response.ResponseSuccess(recycleRecordData.ID)
}

func HandleRecycleOrder(c *gin.Context, orders []models.AiMatchProductOrderModel) error {
	tx := cli.HotDogGormDB.WithContext(c).Begin()
	orderIds := make([]int64, 0)
	for _, order := range orders {
		orderIds = append(orderIds, order.ID)

		err := tx.Model(&models.AiMatchProductOrderModel{}).Where("id", order.ID).
			Updates(map[string]any{
				"status":    67,
				"note":      "RECYCLING",
				"is_delete": 1,
			}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		var count int64
		err = tx.Model(&models.AiMatchRestNftProductModel{}).Where(map[string]any{
			"product_id":          order.ProductId,
			"nft_product_size_id": order.NftProductSizeId,
			"receiver_city":       order.ReceiverCity,
			"receiver_province":   order.ReceiverProvince,
			"is_release":          0,
			"is_delete":           0,
		}).Count(&count).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if count > 0 {
			continue
		}
		if order.NftProductSizeId > 0 {
			iosSourceFile := ""
			androidSourceFile := ""
			if lo.Contains[string](PFP_LIST, order.SourceType) {
				var productSize models.SaleProductNftSizePfpModel
				err := tx.Where(map[string]any{
					"product_id":          order.ProductId,
					"nft_product_size_id": order.NftProductSizeId,
					"receiver_city":       order.ReceiverCity,
					"is_delete":           0,
				}).First(&productSize).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return RecycleNoNftSizeError
					}
					tx.Rollback()
					return err
				}
				iosSourceFile = productSize.OriginMedia
				androidSourceFile = productSize.OriginMedia

			} else {
				var productSize models.SaleProductNftSizeModel
				err := tx.Where(map[string]any{
					"nft_product_size_id": order.NftProductSizeId,
				}).First(&productSize).Error
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return RecycleNoNftSizeError
					}
					tx.Rollback()
					return err
				}
				iosSourceFile = productSize.IosSourceFile
				androidSourceFile = productSize.AndroidSourceFile
			}
			err := tx.Create(&models.AiMatchRestNftProductModel{
				ProductId:         order.ProductId,
				ProductTitle:      order.ProductTitle,
				SizeId:            order.SizeId,
				NftProductSizeId:  order.NftProductSizeId,
				ActivePicture:     order.ProductPicture,
				IosSourceFile:     iosSourceFile,
				AndroidSourceFile: androidSourceFile,
				ReceiverProvince:  order.ReceiverProvince,
				ReceiverCity:      order.ReceiverCity,
				ReceiverRegion:    order.ReceiverRegion,
				CombineInfo:       order.LogisticsInfo,
				CombineActiveId:   order.NewFlashId,
				PlayMethod:        order.PlayMethod,
				CouponIds:         datatypes.JSON([]byte(`[]`)),
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			err := tx.Create(&models.AiMatchRestNftProductModel{
				ProductId:         order.ProductId,
				ProductTitle:      order.ProductTitle,
				SizeId:            order.SizeId,
				NftProductSizeId:  order.NftProductSizeId,
				ActivePicture:     order.ProductPicture,
				IosSourceFile:     "",
				AndroidSourceFile: "",
				ReceiverProvince:  order.ReceiverProvince,
				ReceiverCity:      order.ReceiverCity,
				ReceiverRegion:    order.ReceiverRegion,
				BoxContent:        order.BoxContent,
				CouponIds:         datatypes.JSON([]byte(`[]`)),
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}
	// 更新ProductNftSecond
	err := tx.Model(&models.AiMatchProductNftSecondModel{}).
		Where("order_id IN (?)", orderIds).
		Where("is_delete", 0).
		Where("status", "on_shelf").
		Update("status", "off_shelf").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func GetSimpleList(c *gin.Context) {
	req := form.NftSimpleListReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}

	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 50
	}
	var marketType []string
	productList, err := models.SaleCalendarProduct{Ctx: c}.ListByCondition(models.ListByConditionReq{
		Page:         req.PageNumber,
		PageSize:     req.PageSize,
		ProductId:    req.Id,
		ProductName:  req.Name,
		IsBusiness:   int32(req.IsBusiness),
		Fields:       "id,product_name,market_type,nft_type,new_pictures,is_delete,on_sale_status,price",
		OnSaleStatus: 0,
		MarketType:   marketType,
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	var (
		productIds      []uint64
		productNft      = make(map[int64][]*dto.SimpleProductNft)
		productMapCount = make(map[int64]int64)
	)
	for _, v := range productList {
		productIds = append(productIds, v.Id)
	}
	var productNftSizeIds []uint64
	// 获取nft藏品信息
	nftProductList, err := models.NewSaleProductNftSize().NftProductListByCondition(c, models.SaleProductNftSizeListByConditionReq{
		ProductIds: productIds,
	})
	for _, v := range nftProductList {
		productNftSizeIds = append(productNftSizeIds, uint64(v.NftProductSizeId))
	}
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	restResp, err := models.AiMatchRestNftProduct{Ctx: c}.RestNftProductAvailableStock(models.RestNftProductAvailableStockReq{
		ProductIds:    productIds,
		ProductSizIds: productNftSizeIds,
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	for _, item := range nftProductList {
		var availableStock uint32
		for _, ii := range restResp {
			if int64(ii.ProductId) == item.ProductId && int64(ii.ProductSizId) == item.NftProductSizeId {
				availableStock = ii.AvailableStock
			}
		}
		productNft[item.ProductId] = append(productNft[item.ProductId], &dto.SimpleProductNft{
			Id:             uint64(item.NftProductSizeId),
			Name:           item.ProductTitle,
			AvailableStock: availableStock,
			Image:          item.NftCoverPic,
		})
		productMapCount[item.ProductId] += item.TotalCount
	}
	simpleProductList := make([]dto.SimpleProductList, 0)
	for _, item := range productList {
		if item.IsDelete == 1 { // 删除
			continue
		}

		image := ""
		if item.NewPictures != "" {
			var pictures []string
			err = json.Unmarshal([]byte(item.NewPictures), &pictures)
			if err == nil && len(pictures) > 0 {
				image = pictures[0]
			}
		}
		child, ok := productNft[int64(item.Id)]
		if !ok {
			child = []*dto.SimpleProductNft{}
		}
		simpleProductList = append(simpleProductList, dto.SimpleProductList{
			Id:         item.Id,
			Name:       item.ProductName,
			Image:      image,
			MarketType: item.MarketType,
			NftType:    item.NftType,
			TotalCount: cast.ToUint32(productMapCount[int64(item.Id)]),
			Child:      child,
			Price:      item.Price,
		})
	}
	response.ResponseSuccess(simpleProductList)
}

func NftRestList(c *gin.Context) {
	req := form.NftRestListReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	where := make(map[string]any)
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	where["is_release"] = 0
	where["is_delete"] = 0
	where["product_id"] = req.ProductId
	where["nft_product_size_id"] = req.NftProductSizeId
	list, count, err := models.AiMatchRestNftProduct{Ctx: c}.GetRestNftProductList(where, []string{"id desc"}, req.PageNumber, req.PageSize)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccessWithList(list, int(count))
}

func NftDestroyByCount(c *gin.Context) {
	req := form.NftDestroyByCountReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	userId, _ := c.Get("user_id")
	admUser, err := models.User{Ctx: c}.FindAdmUserById(cast.ToUint64(userId))
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 防止频繁攻击
	if !until.RequestRateLimit(c, fmt.Sprintf("destroy:%d_%d", req.NftProductSizeID, req.ProductID), time.Minute*5) {
		response.Responses(errno.RequestTooOfter, errno.MsgFlags[errno.RequestTooOfter], nil)
		return
	}
	user, err := models.User{Ctx: c}.FindSysUserByMobile(req.Mobile)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 获取相应的订单
	orderList, err := models.AiMatchProductOrder{Ctx: c}.GetProductOrderList(map[string]any{
		"receiver_name":       user.UserId,
		"product_type":        "NFT",
		"status":              2,
		"is_delete":           0,
		"apply_refund_status": 0,
		"nft_product_size_id": req.NftProductSizeID,
		"product_id":          req.ProductID,
	}, &req.Count)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	if len(orderList) != req.Count {
		response.Responses(10112, "数量不足", nil)
		return
	}
	for _, v := range orderList {
		if v.MarketType == "land" {
			response.ResponseFails("暂不支持摧毁土地", 1008)
			return

		} else if lo.Contains[string]([]string{"***", ""}, v.ReceiverCity) {
			response.ResponseFails("订单还没发编号，不能摧毁", 1008)
			return
		}
	}
	orderIds := []int64{}
	for _, v := range orderList {
		orderIds = append(orderIds, v.ID)
	}
	err = models.AiMatchProductOrder{Ctx: c}.UpdateProductOrderByParams(map[string]any{
		"id": orderIds,
	}, map[string]any{
		"status": 76,
		"note":   req.Remark,
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	productId := req.ProductID
	productTitle := orderList[0].ProductTitle
	p := struct {
		Title        string      `json:"title"`
		Nums         int         `json:"nums"`
		UserInfoList [][2]string `json:"userInfoList"`
	}{
		Title: productTitle,
		Nums:  req.Count,
		UserInfoList: [][2]string{
			{cast.ToString(user.UserId), req.Mobile},
		},
	}
	jsonStr, _ := json.Marshal(p)

	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      int64(admUser.UserId),
		Username:    admUser.Name,
		Remark:      req.Remark,
		Scenes:      76,
		AssociateId: int64(productId),
		RequestData: string(jsonStr),
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	response.ResponseSuccess(nil)
}

func NftSecondPriceFlushUserPercentage(c *gin.Context) {
	req := form.NftSecondPriceFlushUserPercentageReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	userId, _ := c.Get("user_id")
	flushKey := fmt.Sprintf("NFTFlushUserPercentageAdm:%d:%d-%d", userId, req.ProductId, req.NftProductSizeId)

	st := cli.HotDogRedis.SetNX(c, flushKey, "lock", 2*time.Second).Val()
	if !st {
		response.Responses(10112, "", nil)
		return
	}
	realUserSurplus, nftCount, _, err := models.Nft{Ctx: c}.NftSecondPriceFlushUserPercentage(req)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.ResponseFail(err.Error())
		return
	}
	if nftCount == 0 {
		response.ResponseFail("该藏品还未确定展示剩余份数，请确认后刷新")
		return
	}
	if realUserSurplus == 0 {
		response.ResponseFail("暂无用户拥有该藏品")
		return
	}
	response.ResponseSuccess(nil)
}

func NftAirdropByCount(c *gin.Context) {
	req := form.NftAirdropByCountReq{}
	validate := until.ApiValidator{ReqData: &req}
	response := until.NewResponse(c)
	if !validate.Validator(c) {
		response.Responses(errno.Error, validate.GetError(), nil)
		return
	}
	userId, _ := c.Get("user_id")
	admUser, err := models.User{Ctx: c}.FindAdmUserById(cast.ToUint64(userId))
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// 防止频繁攻击
	propBoxOpenKey := "SendRestNftProductTOSomeone"
	st := cli.HotDogRedis.SetNX(c, propBoxOpenKey, "lock", 6*time.Second).Val()
	if !st {
		response.Responses(10112, "", nil)
		return
	}
	user, err := models.User{Ctx: c}.FindSysUserByMobile(req.Mobile)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// resp, err := httpReq.HotStoreInfo(os.Getenv("NO_EXPIRE_USER_TOKEN"), dto.HotStoreInfoReq{
	// 	ProductId:        req.ProductID,
	// 	NftProductSizeId: req.NftProductSizeID,
	// })
	// if err != nil {
	// 	response.ResponseFail(err.Error())
	// 	return
	// }
	// if resp.Data.StoreCount < req.Count {
	// 	response.ResponseFail(fmt.Sprintf("藏品库存不足，当前可用库存份数%d,请检查藏品可用库存份数(可用库存=藏品剩余份数-普通用户份数-锁定库存份数)", resp.Data.StoreCount))
	// 	return
	// }
	productData, err := models.NewAiMatchProductNftSecondPrice(c).GetByProductIdAndNftProductSizeId(req.ProductID, req.NftProductSizeID)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	p := struct {
		Title        string                    `json:"title"`
		Nums         int                       `json:"nums"`
		UserInfoList [][2]string               `json:"userInfoList"`
		RequestData  form.NftAirdropByCountReq `json:"request_data"`
	}{
		Title: productData.ProductTitle,
		Nums:  req.Count,
		UserInfoList: [][2]string{
			{cast.ToString(user.UserId), req.Mobile},
		},
		RequestData: req,
	}
	jsonStr, _ := json.Marshal(p)

	err = models.OperateRecord{Ctx: c}.CreateRecord(models.AiMatchBackendOperateRecord{
		UserId:      int64(admUser.UserId),
		Username:    admUser.Name,
		Remark:      "空投藏品",
		Scenes:      69,
		AssociateId: int64(req.ProductID),
		RequestData: string(jsonStr),
	})
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	recycleRecordData := models.AiMatchProductNftRecycleRecord{
		ProductTitle:       productData.ProductTitle,
		ProductId:          req.ProductID,
		NftProductSizeId:   req.NftProductSizeID,
		RecycleTargetCount: req.Count,
		Status:             0,
		UserId:             int(user.UserId),
		AdmUserName:        admUser.Name,
		Type:               constant.BATCH_RECORD_AIRDROP_TYPE,
		OperatorId:         admUser.UserId,
	}
	err = models.AiMatchProductNftRecycleRecord{}.Create(c, &recycleRecordData)
	if err != nil {
		response.ResponseFail(err.Error())
		return
	}
	// // 发送消息到飞书
	// httpReq.FeiShuRootBot("%s(%d) 空投 %s (mobile: %s) 如下藏品: 《%s》%d 份", admUser.Name, admUser.UserId, user.RealName, user.Mobile, productData.ProductTitle, req.Count)

	response.ResponseSuccess(recycleRecordData.ID)
}
