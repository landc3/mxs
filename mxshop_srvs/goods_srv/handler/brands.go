package handler

import (
	"context"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// 用于处理客户端发送的品牌列表查询请求，并返回相应的品牌信息列表
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	//首先创建了一个空的 BrandListResponse 对象，这是根据 .proto 文件中定义的消息类型生成的，用来封装返回给客户端的数据。
	brandListResponse := proto.BrandListResponse{}

	//获取品牌列表
	var brands []model.Brands //定义了一个 brands 切片来存储从数据库中查询到的品牌数据。
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}

	//计算总数
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	var brandResponses []*proto.BrandInfoResponse //用于存放转换后的品牌信息。
	for _, brand := range brands {
		//将每一个从数据库中获取的 brand 实例的信息（如 ID、名称、Logo）转换成 gRPC 服务定义中的 BrandInfoResponse 消息格式，并添加到 brandResponses 切片中。
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	//将 brandResponses 赋值给 brandListResponse.Data 字段，表示本次查询返回的品牌列表数据。
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}
func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//新建品牌
	if result := global.DB.Where("name=?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)

	return &proto.BrandInfoResponse{Id: brand.ID}, nil
}
func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &emptypb.Empty{}, nil
}
func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	global.DB.Save(&brands)

	return &emptypb.Empty{}, nil
}
