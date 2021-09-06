package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type SubscribeModel struct {
	ProposalId string `gorm:"primaryKey" json:"proposalId"`
	Email string `gorm:"primaryKey" json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SubscribeModel) TableName() string {
	return "subscribes"
}

type SqlClient struct {
	db *gorm.DB
}

func (c *SqlClient) Init() error {
	var err error
	c.db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	err = c.db.AutoMigrate(&SubscribeModel{})
	if err != nil{
		return err
	}

	return  nil
}

func (c *SqlClient) NewSubscribe(subscribe SubscribeModel) error {
	res := c.db.Find(&subscribe)
	if res.Error != nil{
		return res.Error
	}
	if res.RowsAffected != 0{
		return nil
	}
	if err := c.db.Create(&subscribe).Error; err != nil{
		return err
	}
	return nil
}

func (c *SqlClient) GetSubscribeEmail(proposalId string) ([]string, error) {
	var subscribes []SubscribeModel
	if err := c.db.Select("email").Where("proposal_id = ?", proposalId).Find(&subscribes).Error; err != nil{
		return nil, err
	}

	emails := []string{}
	for _, value := range subscribes {
		emails = append(emails, value.Email)
	}
	return emails, nil
}