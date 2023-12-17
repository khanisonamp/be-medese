package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Model struct
type Model struct {
	Id        uuid.UUID  `json:"id" gorm:"primary_key;index"`
	Seq       int64      `json:"seq" gorm:"primary_key;auto_increment:false;"`
	CreatedAt time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// BeforeCreate hook table
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if uuid.Equal(m.Id, uuid.Nil) {
		m.Id = uuid.NewV4()
	}
	m.Seq = time.Now().UnixNano()
	return nil
}

type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
