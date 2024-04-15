package models

//foreignKey外键  如果是表名称加上Id的话默认也可以不配置   如果不是我们需要通过foreignKey配置外键
//references表示的是主键    默认就是Id   如果是Id的话可以不配置
type ArticleCate struct {
	Id      int //主键
	Title   string
	State   int
	Article []Article `gorm:"foreignKey:CateId;references:Id"`
}

//表示配置操作数据库的表名称
func (ArticleCate) TableName() string {
	return "article_cate"
}
