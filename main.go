/**
Golang upper v1.11
with echo, gorm modules with mysql Database sample
by Yoshi Sakai
 */
package main
import (
    "io"
    "html/template"
    "net/http"
    "github.com/bluemooninc/docker-go-banner/configs"
    "github.com/bluemooninc/docker-go-banner/banner"
    "github.com/labstack/echo"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    )
var Database *gorm.DB

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func indexController(c echo.Context) error {
    conf := configs.LoadConfig()
    // get banner when exist
    bannerData := banner.GetActiveBanner(Database, c, conf.InternalIps)
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
    return banner.ReturnJson(c, b)
}

func updateController(c echo.Context) error{
    u := new(banner.UserParam)
    if err := c.Bind(u); err != nil {
        return err
    }
    b := banner.Update(Database, u)
    return banner.ReturnJson(c, b)
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
    Database = configs.GormConnect()
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
