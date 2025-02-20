package model


type CardTag struct {
	CardID    uint `json:"card_id" gorm:"primaryKey"`
	TagID     uint `json:"tag_id" gorm:"primaryKey"`
  }	
