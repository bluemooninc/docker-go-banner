package banner

import (
	"fmt"
	"net/http"
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
type BannerItem struct {
	PromotionCode string `json:"PromotionCode"`
	ContentUrl string `json:"ContentUrl"`
	StartedAt string `json:"StartedAt"`
	ExpiredAt string `json:"ExpiredAt"`
}
/*time.Time
** Config for MySQL setting
 */
func gormConnect() *gorm.DB {
	DBMS     := "mysql"
	USER     := "docker"
	PASS     := "docker"
	PROTOCOL := "tcp(mysql_host:3306)"
	DBNAME   := "test_database"
	// add parseTime option
	CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME+"?parseTime=true"
	db,err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db
}
/*
** Common return of Json
 */
func returnJson(c echo.Context, b *BannerItem) error {
	return c.JSON(
		http.StatusOK,
		struct {
			Id int    `json:"Id"`
			PromotionCode string `json:"PromotionCode"`
			ContentUrl string `json:"ContentUrl"`
			StartedAt string `json:"StartedAt"`
			ExpiredAt string `json:"ExpiredAt"`
		}{
			Id: http.StatusOK,
			PromotionCode: b.PromotionCode,
			ContentUrl: b.ContentUrl,
			StartedAt: b.StartedAt,
			ExpiredAt: b.ExpiredAt,
		},
	)
}
/*
** Insert record by Json data
 */
func Insert(c echo.Context) error {
	db := gormConnect()
	db.AutoMigrate(&Banner{})
	u := new(BannerItem)
	if err := c.Bind(u); err != nil {
		return err
	}
	// Convert to timezone offset-aware UTC
	st,_ := httpdate.Str2Time(u.StartedAt, nil)
	ex,_ := httpdate.Str2Time(u.ExpiredAt, nil)
	db.Create(&Banner{
		PromotionCode: u.PromotionCode,
		ContentUrl: u.ContentUrl,
		StartedAt: st,
		ExpiredAt: ex,
	})
	fmt.Println(ex)
	defer db.Close()
	return returnJson(c, u)
}
/*
** Find record by Json data
 */
func Find(c echo.Context) error {
	var Banner Banner

	db := gormConnect()
	u := new(BannerItem)
	if err := c.Bind(u); err != nil {
		return err
	}
	db.First(&Banner, "promotion_code = ?", u.PromotionCode)
	defer db.Close()
	return returnJson(c, u)
}
/*
** Update record by Json data
 */
func Update(c echo.Context) error {
	db := gormConnect()
	db.AutoMigrate(&Banner{})
	u := new(BannerItem)
	if err := c.Bind(u); err != nil {
		return err
	}
	var Banner Banner
	// Convert to timezone offset-aware UTC
	st,_ := httpdate.Str2Time(u.StartedAt, nil)
	ex,_ := httpdate.Str2Time(u.ExpiredAt, nil)
	db.LogMode(true)
	db.Model(&Banner).Where("promotion_code = ?", u.PromotionCode).
		Updates(map[string]interface{}{"ContentUrl": u.ContentUrl, "StartedAt": st, "ExpiredAt": ex})
	defer db.Close()
	return returnJson(c, u)
}
/*
** Delete record by Json data
 */
func Delete(c echo.Context) error {
	db := gormConnect()
	db.AutoMigrate(&Banner{})
	u := new(BannerItem)
	if err := c.Bind(u); err != nil {
		return err
	}
	var Banner Banner
	db.First(&Banner, "promotion_code = ?", u.PromotionCode)
	db.Delete(&Banner)
	defer db.Close()
	return returnJson(c, u)
}
/*
** Find record by Json data
 */
func GetActiveBanner(c echo.Context) *Banner {
	var Banner Banner
	atNow := time.Now()

	fmt.Println(atNow)
	fmt.Println(atNow.In(time.UTC))
	db := gormConnect()
	db.LogMode(true)
	// After a banner expires, it should not be displayed again.
	// There may be occasions where two banners are considered active. In this case, the banner with the earlier expiration should be displayed.
	if c.RealIP() == "10.0.0.1" || c.RealIP() == "172.29.0.1" {
		// if the user has an internal IP address (10.0.0.1, 10.0.0.2), even if the current time is before the display period of the banner.
		db.Where("expired_at > ?", atNow.In(time.UTC)).Order("expired_at asc").First(&Banner)
	} else {
		db.Where("started_at < ? AND expired_at > ?", atNow.In(time.UTC), atNow.In(time.UTC)).Order("expired_at asc").First(&Banner)
	}
	defer db.Close()
	return &Banner
}
