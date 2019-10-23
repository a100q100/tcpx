package configModel

type Config struct {
	Id         int    `gorm:"column:id;default:" json:"id" form:"id"`
	ConfigId   string `gorm:"column:config_id;default:" json:"config_id" form:"config_id"`
	Data       []byte `gorm:"column:data;default:" json:"data" form:"data"`
	DataString string `gorm:"column:data_string;default:" json:"data_string" form:"data_string"`
	Env        string `gorm:"column:env;default:" json:"env" form:"env"`
}

func (o Config) TableName() string {
	return "config"
}
