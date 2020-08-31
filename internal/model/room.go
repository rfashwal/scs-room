package model

import "time"

type Room struct {
	ID                 *uint      `gorm:"primary_key"`
	Name               *string    `gorm:"column:name"`
	AculatorValue      *uint      `gorm:"column:aculator_value"`
	TempratureRequired *float64   `gorm:"column:temprature_required"`
	TempratureCurrent  *float64   `gorm:"column:temprature_current"`
	CreatedAt          *time.Time `gorm:"column:created_at"`
	ModifiedAt         *time.Time `gorm:"column:modified_at"`
}
