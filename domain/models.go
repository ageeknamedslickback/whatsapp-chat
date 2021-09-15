package domain

import (
	"time"

	"gorm.io/gorm"
)

// gorm.Model convention definition
type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Message represents Twilio's message data
type Message struct {
	gorm.Model
	AccountSid          string    `json:"account_sid"`
	From                string    `json:"from"`
	To                  string    `json:"to"`
	Body                string    `json:"body"`
	NumMedia            string    `json:"num_media"`
	NumSegments         string    `json:"num_segments"`
	APIVersion          string    `json:"api_version"`
	ProfileName         string    `json:"profile_name"`
	SmsMessageSid       string    `json:"sms_message_sid"`
	SmsSid              string    `json:"sms_sid"`
	SmsStatus           string    `json:"sms_status"`
	WaID                string    `json:"wa_id"`
	MediaContentType    string    `json:"media_content_type"`
	MediaURL            string    `json:"media_url"`
	TimeStamp           time.Time `json:"timestamp"`
	DateCreated         string    `json:"date_created"`
	DateSent            string    `json:"date_sent"`
	DateUpdated         string    `json:"date_updated"`
	Direction           string    `json:"direction"`
	ErrorCode           *string   `json:"error_code"`
	ErrorMessage        *string   `json:"error_message"`
	MessagingServiceSid string    `json:"messaging_service_sid"`
	Price               *string   `json:"price"`
	PriceUnit           *string   `json:"price_unit"`
	Sid                 string    `json:"sid"`
	Status              string    `json:"status"`
	URI                 string    `json:"uri"`
}

// Sender reps a new whatsapp sender, details of a person/number that sends the message
type Sender struct {
	PhoneNumber string `json:"phoneNumber"`
	ProfileName string `json:"profileName"`
}

// Response reps Twilio's webhook response
type Response struct{}
