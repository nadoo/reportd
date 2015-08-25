package main

import (
	"log"
	"net/http"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/koding/multiconfig"
)

var conf = new(Config)

type Reportd struct {
	db *sqlx.DB
}

type Row []interface{}

type Result struct {
	Title   string   `json:"title"`
	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`
}

type (
	Config struct {
		Debug   bool   `default:false`
		Listen  string `default:":8080"`
		DBType  string `default:"mysql"`
		DBConn  string
		Reports []Report
	}

	Report struct {
		Title string
		Sql   string
	}
)

func (self *Reportd) getIndex(c *gin.Context) {
	db := self.db
	data := gin.H{"PageTitle": "Report"}

	var results = []Result{}
	for _, v := range conf.Reports {
		rows, err := db.Queryx(v.Sql)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		cols, _ := rows.Columns()
		result := Result{Title: v.Title, Columns: cols}

		for rows.Next() {
			obj, err := rows.SliceScan()
			for i, item := range obj {
				if item == nil {
					obj[i] = nil
				} else {
					t := reflect.TypeOf(item).Kind().String()

					if t == "slice" {
						obj[i] = string(item.([]byte))
					}
				}
			}

			if err == nil {
				result.Rows = append(result.Rows, obj)
			}
		}

		results = append(results, result)
	}

	data["Results"] = results
	c.HTML(http.StatusOK, "index", data)
}

func main() {
	m := multiconfig.NewWithPath("config.toml")
	err := m.Load(conf)
	if err != nil {
		log.Fatal(err)
	}

	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := sqlx.Open(conf.DBType, conf.DBConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	reportd := Reportd{db: db}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", reportd.getIndex)
	router.Run(conf.Listen)
}
