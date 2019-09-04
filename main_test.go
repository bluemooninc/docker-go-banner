package main

import (
	"testing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"

)

/*
** Config for Database table
 */
type Banner struct {
	gorm.Model
	PromotionCode string
	ContentUrl string
	StartedAt time.Time
	ExpiredAt time.Time
}

var (
	userJSON = `{"PromotionCode":"unitTest","ContentUrl":"https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png"}`
)
// a successful case
func TestCreateBanner(t *testing.T) {
	// Setup
	db := GormConnect()
	now := time.Now()
	st := now.Add(-time.Hour)
	ex := now.Add(time.Hour)
	// Create
	db.Create(&Banner{
		PromotionCode: "PromotionCode",
		ContentUrl: "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_92x30dp.png",
		StartedAt: st,
		ExpiredAt: ex,
	})

	// Assertions
//	if assert.NoError(t, h.createBanner(c)) {
//		assert.Equal(t, http.StatusCreated, rec.Code)
//		assert.Equal(t, userJSON, rec.Body.String())
//	}
}
