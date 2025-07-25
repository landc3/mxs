package model

import (
	"context"
	"strconv"

	"gorm.io/gorm"

	"mxshop_srvs/goods_srv/global"

)

//类型， 这个字段是否能为null， 这个字段应该设置为可以为null还是设置为空， 0
//实际开发过程中 尽量设置为不为null
//https://zhuanlan.zhihu.com/p/73997266
//这些类型我们使用int32还是int
type Category struct{
	BaseModel
	Name  string `gorm:"type:varchar(20);not null" json:"name"`
	ParentCategoryID int32 `json:"parent"`
	ParentCategory *Category `json:"-"`
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level int32 `gorm:"type:int;not null;default:1" json:"level"` //1级 2级 3级
	IsTab bool `gorm:"default:false;not null" json:"is_tab"` //是否显示
}


type Brands struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null"`
	Logo  string `gorm:"type:varchar(200);default:'';not null"`
}

//商品和品牌关系
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands Brands
}
//商品和品牌关系表名自定义
func (GoodsCategoryBrand) TableName() string{
	return "goodscategorybrand"
}
//商品轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	Url string `gorm:"type:varchar(200);not null"`
	Index int32 `gorm:"type:int;default:1;not null"`
}

type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category Category
	BrandsID int32 `gorm:"type:int;not null"`
	Brands Brands

	OnSale bool `gorm:"default:false;not null"` //是否上架
	ShipFree bool `gorm:"default:false;not null"`//是否包邮
	IsNew bool `gorm:"default:false;not null"`//是否新品
	IsHot bool `gorm:"default:false;not null"` //是否热销

	Name  string `gorm:"type:varchar(50);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"`  // 商品编号
	ClickNum int32 `gorm:"type:int;default:0;not null"` //点击数
	SoldNum int32 `gorm:"type:int;default:0;not null"` //销量
	FavNum int32 `gorm:"type:int;default:0;not null"`//收藏数
	MarketPrice float32 `gorm:"not null"`//商品市场价
	ShopPrice float32 `gorm:"not null"`//商品销售价
	GoodsBrief string `gorm:"type:varchar(100);not null"`//商品简单描述
	Images GormList `gorm:"type:varchar(1000);not null"`//商品图片
	DescImages GormList `gorm:"type:varchar(1000);not null"`//商品描述图片
	GoodsFrontImage string `gorm:"type:varchar(200);not null"`//商品图片
}
//商品添加成功后，同步到es中
func (g *Goods) AfterCreate(tx *gorm.DB) (err error){
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) (err error){
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}

	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).
		Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (g *Goods) AfterDelete(tx *gorm.DB) (err error){
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}