package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// GormList 自定义类型，存储为 json 字符串，将图片数组存储为 json 字符串
type GormList []string

func (g GormList) Value() (driver.Value, error){
	return json.Marshal(g) // 将 gorm.List 转换为 json
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g) // 将 value 转换为 []byte
}
//BaseModel 为所有模型提供了统一的字段定义和命名规范。在不同的微服务中，虽然业务逻辑可能不同，
//但对于通用字段的处理应该保持一致。通过使用 BaseModel，可以确保所有模型的主键、创建时间、更新时间等字段的定义和命名都是相同的
type BaseModel struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"` //为什么使用int32， bigint
	CreatedAt time.Time `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDeleted bool `json:"-"`
}
