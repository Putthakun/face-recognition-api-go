package entity

import "time"

type Employee struct {
	EmpID         int64          `gorm:"primaryKey;column:EmpId"`
	Name          string         `gorm:"column:Name;not null"`
	Credential    *Credential    `gorm:"foreignKey:EmpID"`
	FaceEmbeddeds []FaceEmbedded `gorm:"foreignKey:EmpID"`
	Transactions  []Transaction  `gorm:"foreignKey:EmpID"`
}

func (Employee) TableName() string { return "Employees" }

type Credential struct {
	CredentialID int64     `gorm:"primaryKey;column:CredentialId;autoIncrement"`
	EmpID        int64     `gorm:"column:EmpId;not null;uniqueIndex"`
	PasswordHash string    `gorm:"column:PasswordHash;not null"`
	Role         string    `gorm:"column:Role;not null;default:Employee"`
	IsActive     bool      `gorm:"column:IsActive;not null;default:true"`
	CreatedAt    time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	Employee     *Employee `gorm:"foreignKey:EmpID"`
}

func (Credential) TableName() string { return "Credentials" }

type FaceEmbedded struct {
	FaceEmbeddedID   int64     `gorm:"primaryKey;column:FaceEmbeddedId;autoIncrement"`
	EmpID            int64     `gorm:"column:EmpId;not null"`
	FaceEmbeddedData []byte    `gorm:"column:FaceEmbeddedData;type:varbinary(max)"`
	CreatedAt        time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	Employee         *Employee `gorm:"foreignKey:EmpID"`
}

func (FaceEmbedded) TableName() string { return "FaceEmbeddeds" }
