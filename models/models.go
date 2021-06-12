package models
type AreaLogs struct{
	ID int64 `gorm:"primary_key;auto_increment;not_null" json:"id"`
	AreaID int64 `gorm:"not_null" json:"area_id"`//库房
	WID int64 `json:"wid"`//如果有
	PID string `json:"pid"`
	Content string `json:"content"` //说明
	Time string `json:"time"`
}