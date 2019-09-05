/**
Golang upper v1.11
with echo, gorm modules with mysql Database sample
by Yoshi Sakai
 */
package main
import (
    "io"
    "net/http"
    "html/template"
    "github.com/castaneai/gomodtest/banner"
    "github.com/labstack/echo"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    )

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

var Database *gorm.DB

/*
** Config for MySQL setting
 */
func GormConnect() *gorm.DB {
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
func returnJson(c echo.Context, b *banner.UserParam) error {
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

func indexController(c echo.Context) error {
    // get banner when exist
    bannerData := banner.GetActiveBanner(Database, c)
    // Template rendering
    return c.Render(http.StatusOK, "index", bannerData)
}

func insertController(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    bannerData := banner.Insert(Database, u)
    return c.JSON(http.StatusOK, bannerData)
}

func findController(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    b := banner.Find(Database, u)
    return returnJson(c, b)
}

func updateController(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    b := banner.Update(Database, u)
    return returnJson(c, b)
}

func deleteController(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    banner.Delete(Database, u)
    return nil
}
/*
** main loop
 */
func main() {
    t := &Template{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
    e := echo.New()
    e.Renderer = t
    Database = GormConnect()
    Database.AutoMigrate(&banner.Banner{})

    // Routes
    e.GET("/", indexController)
    e.POST("/insert", insertController)
    e.POST("/find", findController)
    e.POST("/update", updateController)
    e.POST("/delete", deleteController)
    // Start server
    e.Logger.Fatal(e.Start(":8080"))
    defer Database.Close()
}
