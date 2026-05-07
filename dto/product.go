package dto

type ProductInfoData struct {
	Id                   uint64  `protobuf:"varint,1,opt,name=id,proto3" json:"id"`                                             // 商品ID
	BrandId              uint32  `protobuf:"varint,2,opt,name=brand_id,json=brandId,proto3" json:"brand_id"`                    // 品牌ID
	Name                 string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`                                          // 商品名称
	ArtNo                string  `protobuf:"bytes,4,opt,name=art_no,json=artNo,proto3" json:"art_no"`                           // 货号
	AfterSales           string  `protobuf:"bytes,5,opt,name=after_sales,json=afterSales,proto3" json:"after_sales"`            // 售后
	SaleTime             uint64  `protobuf:"varint,6,opt,name=sale_time,json=saleTime,proto3" json:"sale_time"`                 // 发售时间
	StartSaleTime        string  `protobuf:"bytes,7,opt,name=start_sale_time,json=startSaleTime,proto3" json:"start_sale_time"` // 开售时间
	IsBusiness           uint32  `protobuf:"varint,8,opt,name=is_business,json=isBusiness,proto3" json:"is_business"`
	IsMerchantProvide    uint32  `protobuf:"varint,9,opt,name=is_merchant_provide,json=isMerchantProvide,proto3" json:"is_merchant_provide"`
	ShipmentTime         string  `protobuf:"bytes,10,opt,name=shipment_time,json=shipmentTime,proto3" json:"shipment_time"`      // 商品发货时间
	FreightAmount        float32 `protobuf:"fixed32,11,opt,name=freight_amount,json=freightAmount,proto3" json:"freight_amount"` // 邮费
	Weight               uint32  `protobuf:"varint,12,opt,name=weight,proto3" json:"weight"`                                     // 权重
	Supplier             string  `protobuf:"bytes,13,opt,name=supplier,proto3" json:"supplier"`                                  // 供应商
	DetailDesc           string  `protobuf:"bytes,14,opt,name=detail_desc,json=detailDesc,proto3" json:"detail_desc"`            // 详情
	Pictures             string  `protobuf:"bytes,15,opt,name=pictures,proto3" json:"pictures"`                                  // 商品图片
	WearPictures         string  `protobuf:"bytes,16,opt,name=wear_pictures,json=wearPictures,proto3" json:"wear_pictures"`      // 穿搭图片
	StylePicture         string  `protobuf:"bytes,17,opt,name=style_picture,json=stylePicture,proto3" json:"style_picture"`      // 版型图片
	SizePictures         string  `protobuf:"bytes,18,opt,name=size_pictures,json=sizePictures,proto3" json:"size_pictures"`      // 尺码对照表
	Specification        string  `protobuf:"bytes,19,opt,name=specification,proto3" json:"specification"`                        // 规格参数
	Type                 string  `protobuf:"bytes,20,opt,name=type,proto3" json:"type"`
	ShowType             string  `protobuf:"bytes,21,opt,name=show_type,json=showType,proto3" json:"show_type"`
	NftType              string  `protobuf:"bytes,22,opt,name=nft_type,json=nftType,proto3" json:"nft_type"`
	AuthorName           string  `protobuf:"bytes,23,opt,name=author_name,json=authorName,proto3" json:"author_name"`
	AuthorAvatar         string  `protobuf:"bytes,24,opt,name=author_avatar,json=authorAvatar,proto3" json:"author_avatar"`
	Links                string  `protobuf:"bytes,25,opt,name=links,proto3" json:"links"`
	Price                string  `protobuf:"bytes,26,opt,name=price,proto3" json:"price"`
	Currency             string  `protobuf:"bytes,27,opt,name=currency,proto3" json:"currency"`
	ProductAbbreviation  string  `protobuf:"bytes,28,opt,name=product_abbreviation,json=productAbbreviation,proto3" json:"product_abbreviation"`
	BrandName            string  `protobuf:"bytes,29,opt,name=brand_name,json=brandName,proto3" json:"brand_name"`                // 品牌名称
	Status               uint32  `protobuf:"varint,30,opt,name=status,proto3" json:"status"`                                      // 状态 0=未上架 1=上架
	IsDelete             uint32  `protobuf:"varint,31,opt,name=is_delete,json=isDelete,proto3" json:"is_delete"`                  // 是否删除 0=否
	ThingId              uint64  `protobuf:"varint,32,opt,name=thing_id,json=thingId,proto3" json:"thing_id"`                     // 版权实物商品ID
	IsCanRetrieve        uint32  `protobuf:"varint,33,opt,name=is_can_retrieve,json=isCanRetrieve,proto3" json:"is_can_retrieve"` // 是否可以取回
	OnSaleStatus         uint32  `protobuf:"varint,34,opt,name=on_sale_status,json=onSaleStatus,proto3" json:"on_sale_status"`
	NewPictures          string  `protobuf:"bytes,35,opt,name=new_pictures,json=newPictures,proto3" json:"new_pictures"`              // 商品banner图片
	FreightDiscount      uint64  `protobuf:"varint,36,opt,name=freight_discount,json=freightDiscount,proto3" json:"freight_discount"` // 运费折扣
	CreateTime           string  `protobuf:"bytes,37,opt,name=create_time,json=createTime,proto3" json:"create_time"`
	UpdateTime           string  `protobuf:"bytes,38,opt,name=update_time,json=updateTime,proto3" json:"update_time"`
	IsTest               int32   `protobuf:"varint,39,opt,name=is_test,json=isTest,proto3" json:"is_test"`
	Style                string  `protobuf:"bytes,40,opt,name=style,proto3" json:"style"`
	MarketType           string  `protobuf:"bytes,41,opt,name=market_type,json=marketType,proto3" json:"market_type"`
	Point                uint64  `protobuf:"varint,42,opt,name=point,proto3" json:"point"`
	Extend               string  `protobuf:"bytes,43,opt,name=extend,proto3" json:"extend"`
	IsHighSpider         uint32  `protobuf:"varint,44,opt,name=is_high_spider,json=isHighSpider,proto3" json:"is_high_spider"`
	IsFar                int32   `protobuf:"varint,45,opt,name=is_far,json=isFar,proto3" json:"is_far"`
	SaleMethod           string  `protobuf:"bytes,46,opt,name=sale_method,json=saleMethod,proto3" json:"sale_method"`
	ProductName          string  `protobuf:"bytes,47,opt,name=product_name,json=productName,proto3" json:"product_name"`
	NftCoverPic          string  `protobuf:"bytes,48,opt,name=nft_cover_pic,json=nftCoverPic,proto3" json:"nft_cover_pic"`
	IsOneCanOpen         uint32  `protobuf:"varint,49,opt,name=is_one_can_open,json=isOneCanOpen,proto3" json:"is_one_can_open"`
	BlindboxOpenAmount   uint64  `protobuf:"varint,50,opt,name=blindbox_open_amount,json=blindboxOpenAmount,proto3" json:"blindbox_open_amount"`
	OpenNftNum           uint32  `protobuf:"varint,51,opt,name=open_nft_num,json=openNftNum,proto3" json:"open_nft_num"`
	ProductId            uint64  `protobuf:"varint,52,opt,name=product_id,json=productId,proto3" json:"product_id"`
	NftProductSizeId     uint64  `protobuf:"varint,53,opt,name=nft_product_size_id,json=nftProductSizeId,proto3" json:"nft_product_size_id"`
	IsSelfPickUp         uint32  `protobuf:"varint,54,opt,name=is_self_pick_up,json=isSelfPickUp,proto3" json:"is_self_pick_up"` // 是否自提
	ConsumeNftConfig     string  `protobuf:"bytes,55,opt,name=consume_nft_config,json=consumeNftConfig,proto3" json:"consume_nft_config"`
	DetailHeadDisplayPic string  `protobuf:"bytes,56,opt,name=detail_head_display_pic,json=detailHeadDisplayPic,proto3" json:"detail_head_display_pic"`
	DressingPictures     string  `protobuf:"bytes,57,opt,name=dressing_pictures,json=dressingPictures,proto3" json:"dressing_pictures"`
}

type SimpleProductNft struct {
	Id             uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id"`                                               // nft藏品ID
	Name           string `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`                                            // nft藏品名称
	Image          string `protobuf:"bytes,3,opt,name=image,proto3" json:"image"`                                          // 图片
	AvailableStock uint32 `protobuf:"varint,4,opt,name=available_stock,json=availableStock,proto3" json:"available_stock"` // 可用库存
}

type SimpleProductList struct {
	Id         uint64              `protobuf:"varint,1,opt,name=id,proto3" json:"id"`                                  // 藏品ID
	Name       string              `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`                               // 藏品名称
	Image      string              `protobuf:"bytes,3,opt,name=image,proto3" json:"image"`                             // 藏品图片
	MarketType string              `protobuf:"bytes,4,opt,name=market_type,json=marketType,proto3" json:"market_type"` // 藏品类型
	NftType    string              `protobuf:"bytes,5,opt,name=nft_type,json=nftType,proto3" json:"nft_type"`          // nft类型
	TotalCount uint32              `protobuf:"varint,6,opt,name=total_count,json=totalCount,proto3" json:"total_count"`
	Child      []*SimpleProductNft `protobuf:"bytes,7,rep,name=child,proto3" json:"child"`
	Price      string              `protobuf:"bytes,8,opt,name=price,proto3" json:"price"`
}

type SimpleProductListResp struct {
	List []*SimpleProductList `protobuf:"bytes,1,rep,name=list,proto3" json:"list"` // 列表
}
