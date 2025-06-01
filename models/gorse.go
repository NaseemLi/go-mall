package models

import "time"

type Item struct {
	ItemId     string    `gorm:"primaryKey" mapstructure:"item_id"`
	IsHidden   bool      `mapstructure:"is_hidden"`
	Categories []string  `gorm:"serializer:json" mapstructure:"categories"`
	Timestamp  time.Time `gorm:"column:time_stamp" mapstructure:"timestamp"`
	Labels     any       `gorm:"serializer:json" mapstructure:"labels"`
	Comment    string    `mapsstructure:"comment"`
}

// User stores meta data about user.
type User struct {
	UserId    string   `gorm:"primaryKey" mapstructure:"user_id"`
	Labels    any      `gorm:"serializer:json" mapstructure:"labels"`
	Subscribe []string `gorm:"serializer:json" mapstructure:"subscribe"`
	Comment   string   `mapstructure:"comment"`
}

// FeedbackKey identifies feedback.
type FeedbackKey struct {
	FeedbackType string `gorm:"column:feedback_type" mapstructure:"feedback_type"`
	UserId       string `gorm:"column:user_id" mapstructure:"user_id"`
	ItemId       string `gorm:"column:item_id" mapstructure:"item_id"`
}

// Feedback stores feedback.
type Feedback struct {
	FeedbackKey `gorm:"embedded" mapstructure:",squash"`
	Timestamp   time.Time `gorm:"column:time_stamp" mapsstructure:"timestamp"`
	Comment     string    `gorm:"column:comment" mapsstructure:"comment"`
}

func (Feedback) TableName() string {
	return "feedback"
}
