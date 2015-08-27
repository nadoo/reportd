package main

import (
	"log"
	"net/http"
	"reflect"
	"time"

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
		Title  string
		Params bool `default:false`
		Sql    string
	}
)

func (self *Reportd) getIndexData(c *gin.Context) gin.H {
	db := self.db
	data := gin.H{"PageTitle": "Report", "ReportTime": time.Now().Format("2006-01-02 15:04:05")}

	var args []interface{}                    // sql query args
	var params = make(map[string]interface{}) // url parameters

	for k, v := range c.Request.URL.Query() {
		params[k] = v
	}

	var results = []Result{}
	for _, v := range conf.Reports {
		sqlStr := v.Sql
		var err error

		if v.Params {
			sqlStr, args, err = sqlx.Named(v.Sql, params)
			if err != nil {
				log.Println(err)
			}
			sqlStr, args, err = sqlx.In(sqlStr, args...)
			if err != nil {
				log.Println(err)
			}
			sqlStr = db.Rebind(sqlStr)
		}

		rows, err := db.Queryx(sqlStr, args...)
		if err != nil {
			log.Println(err)
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
	return data
}

func (self *Reportd) getIndex(c *gin.Context) {
	data := self.getIndexData(c)
	c.HTML(http.StatusOK, "index", data)
}

func (self *Reportd) getIndexJson(c *gin.Context) {
	data := self.getIndexData(c)
	c.JSON(http.StatusOK, data)
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
	router.GET("/json", reportd.getIndexJson)
	router.Run(conf.Listen)
}
