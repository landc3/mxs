package router

import (
	"mxshop-api/goods-web/middlewares"

	"github.com/gin-gonic/gin"

	"mxshop-api/goods-web/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods").Use(middlewares.Trace())
	{
		/*
			新增、删除、修改都需要管理员权限
		*/
		GoodsRouter.GET("", goods.List)                                                            //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)          //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                                      //获取商品的详情
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete) //删除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                               //获取商品的库存

		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         //修改商品
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) //修改商品状态
	}
}
