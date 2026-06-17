package entity

type Camera struct {
	CameraID     int64         `gorm:"primaryKey;column:CameraId;autoIncrement"`
	Location     string        `gorm:"column:Location;not null"`
	Transactions []Transaction `gorm:"foreignKey:CameraID"`
}

func (Camera) TableName() string { return "Cameras" }
