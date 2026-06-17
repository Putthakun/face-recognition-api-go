package entity

import "time"

type Transaction struct {
	TransactionID int64     `gorm:"primaryKey;column:TransactionId;autoIncrement"`
	EmpID         *int64    `gorm:"column:EmpId"`
	CameraID      *int64    `gorm:"column:CameraId"`
	CreatedAt     time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	Employee      *Employee `gorm:"foreignKey:EmpID"`
	Camera        *Camera   `gorm:"foreignKey:CameraID"`
}

func (Transaction) TableName() string { return "Transactions" }
