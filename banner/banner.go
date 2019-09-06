package banner

import (
	"fmt"
	"net/http"
	"time"
	"github.com/labstack/echo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/Songmu/go-httpdate"
)

const dateLayout = "2006-01-02 15:04"

/*
** Struct Database table
 */
type Banner struct {
	gorm.Model
	PromotionCode string
	ContentUrl string
	StartedAt time.Time
	ExpiredAt time.Time
}
/*
** Struct Json data
 */
type UserParam struct {
	PromotionCode string `json:"PromotionCode"`
	ContentUrl    string `json:"ContentUrl"`
	StartedAt     string `json:"StartedAt"`
	ExpiredAt     string `json:"ExpiredAt"`
	RemoteAddr    string
}
/*
** Set Json data
 */
func setReturnData(Banner *Banner) *UserParam{
	// Set to view data
	d := UserParam{
		PromotionCode: Banner.PromotionCode,
		ContentUrl:    Banner.ContentUrl,
		StartedAt:     Banner.StartedAt.Format(dateLayout),
		ExpiredAt:     Banner.ExpiredAt.Format(dateLayout),
	}
	return &d
}
/*
 ** Insert record by Json data
 */
func Insert(Database *gorm.DB, u *UserParam) error {
	// Convert to timezone offset-aware UTC
	st,_ := httpdate.Str2Time(u.StartedAt, nil)
	ex,_ := httpdate.Str2Time(u.ExpiredAt, nil)
	Database.Create(&Banner{
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
func Find(Database *gorm.DB, u *UserParam) *UserParam {
	var Banner Banner
	Database.First(&Banner, "promotion_code = ?", u.PromotionCode)
	return setReturnData(&Banner)
}
/*
 ** Update record by Json data
 */
func Update(Database *gorm.DB,u *UserParam) *UserParam {
	var Banner Banner
	// Convert to timezone offset-aware UTC
	st,_ := httpdate.Str2Time(u.StartedAt, nil)
	ex,_ := httpdate.Str2Time(u.ExpiredAt, nil)
	Database.LogMode(true)
	Database.Model(&Banner).Where("promotion_code = ?", u.PromotionCode).
		Updates(map[string]interface{}{"ContentUrl": u.ContentUrl, "StartedAt": st, "ExpiredAt": ex})
	// Set to view data
	return setReturnData(&Banner)
}
/*
 ** Delete record by Json data
 */
func Delete(Database *gorm.DB,u *UserParam) error {
	var Banner Banner
	Database.First(&Banner, "promotion_code = ?", u.PromotionCode)
	Database.Delete(&Banner)
	return nil
}
/**
Check internal IP address
 */
func in_array(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if val == v {
			index = i
			exists = true
			return
		}
	}
	return
}
/*
 ** Find record by Json data
 */
func GetActiveBanner(Database *gorm.DB, c echo.Context, ipAddresses []string) *UserParam {
	var Banner Banner
	atNow := time.Now()

	fmt.Println(atNow)
	fmt.Println(atNow.In(time.UTC))
	Database.LogMode(true)
	// After a banner expires, it should not be displayed again.
	// There may be occasions where two banners are considered active. In this case, the banner with the earlier expiration should be displayed.
	ret,_ := in_array(c.RealIP(), ipAddresses)
	fmt.Println(c.RealIP())
	fmt.Println(ret)
	if ret == true  {
		// if the user has an internal IP address (10.0.0.1, 10.0.0.2), even if the current time is before the display period of the banner.
		Database.Where("expired_at > ?", atNow.In(time.UTC)).Order("expired_at asc").First(&Banner)
	} else {
		Database.Where("started_at < ? AND expired_at > ?", atNow.In(time.UTC), atNow.In(time.UTC)).Order("expired_at asc").First(&Banner)
	}
	// Set to view data
	u := UserParam{
		PromotionCode: Banner.PromotionCode,
		ContentUrl:    Banner.ContentUrl,
		StartedAt:     Banner.StartedAt.Format(dateLayout),
		ExpiredAt:     Banner.ExpiredAt.Format(dateLayout),
		RemoteAddr:    c.RealIP(),
	}
	return &u
}

/*
** Common return of Json
 */
func ReturnJson(c echo.Context, b *UserParam) error {
	if b != nil {
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
			}, )
	}
	return nil
}
