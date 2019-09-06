package configs
import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	)

type Config struct {
	InternalIps []string `json:"internalIps"`
}

func LoadConfig() *Config{
	bytes, err := ioutil.ReadFile("configs/config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var cnf Config
	json.Unmarshal(bytes, cnf)
	return &cnf
}

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
