package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Inspost struct {
	Id        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Footprint string `json:"footprint" binding:"required"`
	Position  string `json:"position" binding:"required"`
	Number    string `json:"number" binding:"required"`
	Remarks   string `json:"remarks" binding:"required"`
}

func cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        method := c.Request.Method

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
        c.Header("Access-Control-Allow-Credentials", "true")
        if method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
        }
        c.Next()
    }
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("[Casapi-error] error:", err.Error())
	}
}

func DoQuery(db *sql.DB, sqlInfo string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(sqlInfo, args...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{}
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{})
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}

func Searchcomponent(c *gin.Context, db *sql.DB) {
	method := c.Param("method")
	value := c.Param("value")
	page := c.Query("page")
	limit := c.Query("limit")

	page_t, err := strconv.Atoi(page)
	limit_t, err := strconv.Atoi(limit)

	if limit_t == 0 || page_t == 0 || err != nil {
		limit_t = 1000
		page_t = 1
	}

	start := page_t*limit_t - limit_t
	end := page_t*limit_t - start

	var count int
	var res []map[string]interface{}

	switch method {
	case "name":
		sqlstr := "SELECT * FROM ELEC WHERE name like \"%" + value + "%\" LIMIT ?,?"
		rows := db.QueryRow("SELECT count(*) FROM ELEC WHERE name like \"%"+value+"%\"", value)
		err = rows.Scan(&count)
		checkErr(err)

		res, err = DoQuery(db, sqlstr, start, end)
		checkErr(err)
		break
	case "id":
		rows := db.QueryRow("SELECT count(*) FROM ELECWHERE id=?", value)
		err = rows.Scan(&count)
		checkErr(err)

		res, err = DoQuery(db, "SELECT * FROM ELEC WHERE id=? LIMIT ?,?", value, start, end)
		checkErr(err)
		break
	case "footprint":
		rows := db.QueryRow("SELECT count(*) FROM ELECWHERE footprint=?", value)
		err = rows.Scan(&count)
		checkErr(err)

		res, err = DoQuery(db, "SELECT * FROM ELEC WHERE footprint=? LIMIT ?,?", value, start, end)
		checkErr(err)
		break
	case "value":
		if value == "all" {
			res, err = DoQuery(db, "SELECT DISTINCT value FROM ELEC")
			checkErr(err)
		} else {
			rows := db.QueryRow("SELECT count(*) FROM ELECWHERE value=?", value)
			err = rows.Scan(&count)
			checkErr(err)

			res, err = DoQuery(db, "SELECT * FROM ELEC WHERE value=? LIMIT ?,?", value, start, end)
			checkErr(err)
		}
		break
	case "position":
		rows := db.QueryRow("SELECT count(*) FROM ELECWHERE position=?", value)
		err = rows.Scan(&count)
		checkErr(err)

		res, err = DoQuery(db, "SELECT * FROM ELEC WHERE position=? LIMIT ?,?", value, start, end)
		checkErr(err)
		break
	default:
		rows := db.QueryRow("SELECT count(*) FROM ELEC")
		err = rows.Scan(&count)
		checkErr(err)

		sqlstr := "SELECT * FROM ELEC LIMIT ?,?"
		res, err = DoQuery(db, sqlstr, start, end)
		checkErr(err)
	}

	m := make(map[string]interface{}, 4)
	m["code"] = 0
	m["msg"] = ""
	m["count"] = count
	m["data"] = res
	res_json, err := json.MarshalIndent(m, "", "    ")
	checkErr(err)
	c.String(http.StatusOK, string(res_json))

}

func Insertcomponent(c *gin.Context, db *sql.DB) {

	var InsPostT Inspost
	contentType := c.Request.Header.Get("Content-Type")
	fmt.Println(contentType)
	err := c.BindJSON(&InsPostT)
	if err != nil {
		checkErr(err)
		c.JSON(200, gin.H{"errcode": 400, "description": "Post Data Err"})
	} else {
		stmt, err := db.Prepare("INSERT INTO ELEC (id,name,value,footprint,position,number,remarks) VALUES (?,?,?,?,?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(InsPostT.Id, InsPostT.Name, InsPostT.Value, InsPostT.Footprint, InsPostT.Position, InsPostT.Number, InsPostT.Remarks)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		c.String(200, "successful insert"+string(id))
	}
}

func Deletecomponent(c *gin.Context, db *sql.DB) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM ELEC WHERE id=?")
	checkErr(err)
	res, err := stmt.Exec(id)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	if affect == 0 {
		c.String(404, "not found")
	} else {
		c.String(200, "successful")
	}
}

func main() {
	r := gin.Default()

        r.USE(cors())

	db, err := sql.Open("sqlite3", "./elec.db")
	checkErr(err)

	//默认路由
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome Component Archive System Api")
	})

	//查询元件
	r.GET("/v1/search/:method/:value", func(c *gin.Context) {
		Searchcomponent(c, db)
	})

	//添加元件
	r.POST("/v1/insert", func(c *gin.Context) {
		Insertcomponent(c, db)
	})

	//删除元件
	r.GET("/v1/delete/:id", func(c *gin.Context) {
		Deletecomponent(c, db)
	})

	r.Run(":3000") // listen and serve on 0.0.0.0:3000
}
