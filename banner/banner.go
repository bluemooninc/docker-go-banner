package banner

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/Songmu/go-httpdate"
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
/*
** Config for Json data
 */
type UserParam struct {
	PromotionCode string `json:"PromotionCode"`
	ContentUrl string `json:"ContentUrl"`
	StartedAt string `json:"StartedAt"`
	ExpiredAt string `json:"ExpiredAt"`
}

/*
** Insert record by Json data
 */
func Insert(DATABASE *gorm.DB, u *UserParam) error {
	// Convert to timezone offset-aware UTC
	st,_ := httpdate.Str2Time(u.StartedAt, nil)
	ex,_ := httpdate.Str2Time(u.ExpiredAt, nil)
	DATABASE.Create(&Banner{
		PromotionCode: u.PromotionCode,
		ContentUrl: u.ContentUrl,
		StartedAt: st,
		ExpiredAt: ex,
	})
	return nil
}
/*
** Find record by Json data
 */
func Find(DATABASE *gorm.DB, u *UserParam) *Banner {
	var Banner Banner
	DATABASE.First(&Banner, "promotion_code = ?", u.PromotionCode)
	return &Banner
}
/*
** Update record by Json data
 */
func Update(DATABASE *gorm.DB,u *UserParam) *Banner {
	var Banner Banner
	// Convert to timezone offset-aware UTC
	st,_ := httpdate.Str2Time(u.StartedAt, nil)
	ex,_ := httpdate.Str2Time(u.ExpiredAt, nil)
	DATABASE.LogMode(true)
	DATABASE.Model(&Banner).Where("promotion_code = ?", u.PromotionCode).
		Updates(map[string]interface{}{"ContentUrl": u.ContentUrl, "StartedAt": st, "ExpiredAt": ex})
	return &Banner
}
/*
** Delete record by Json data
 */
func Delete(DATABASE *gorm.DB,u *UserParam)  *Banner {
	var Banner Banner
	DATABASE.First(&Banner, "promotion_code = ?", u.PromotionCode)
	DATABASE.Delete(&Banner)
	return &Banner
}
/*
** Find record by Json data
 */
func GetActiveBanner(DATABASE *gorm.DB, c echo.Context) *Banner {
	var Banner Banner
	atNow := time.Now()

	fmt.Println(atNow)
	fmt.Println(atNow.In(time.UTC))
	DATABASE.LogMode(true)
	// After a banner expires, it should not be displayed again.
	// There may be occasions where two banners are considered active. In this case, the banner with the earlier expiration should be displayed.
	if c.RealIP() == "10.0.0.1" || c.RealIP() == "172.29.0.1" {
		// if the user has an internal IP address (10.0.0.1, 10.0.0.2), even if the current time is before the display period of the banner.
		DATABASE.Where("expired_at > ?", atNow.In(time.UTC)).Order("expired_at asc").First(&Banner)
	} else {
		DATABASE.Where("started_at < ? AND expired_at > ?", atNow.In(time.UTC), atNow.In(time.UTC)).Order("expired_at asc").First(&Banner)
	}
	return &Banner
}
