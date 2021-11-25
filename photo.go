package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-redis/redis/v8"
)

// global variables
var debug = false
var useImageBorder = true // distinct images with border color
var borderStyleSQL = "10px solid red"
var borderStyleRedis = "10px solid green"

// mssql variables
var mssqlConnection *sql.DB
var mssqlServer = "myserver"
var mssqlPort = 1433
var mssqlUser = "myusername"
var mssqlPassword = "mypassword"
var mssqlSelect = "select photo from mydatabase.dbo.photos where id='%s'" // ensure only one row (image sql data) is selected
var connString string
var errdb error

// redis variables
var ctx = context.Background()
var keyExpiration time.Duration = 10 * time.Second
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0, // default database
})

func fromSQL(id string, one bool) string {

	query := fmt.Sprintf(mssqlSelect, strings.ReplaceAll(id, "'", "")) // Remove simple quotes preventing SQL injection
	if debug {
		log.Printf("sqlQuery: %s", query)
	}
	stmt, err := mssqlConnection.Prepare(query)
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()
	row := stmt.QueryRow()
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	var rawImage string
	err = row.Scan(&rawImage)
	if err != nil {
		log.Println("Scan failed:", err.Error())
		return ""
	}
	imageBase64 := base64.StdEncoding.EncodeToString([]byte(rawImage))
	return "<img src='data:image/jpg;base64," + imageBase64 + "'"
}

func getImage(id string) string {
	retImage, err := rdb.Get(ctx, "image:"+id).Result()

	if err != nil { //we do not have a redis image, getting from sql
		retImage = fromSQL(id, true)
		if len(retImage) == 0 {
			return ""
		}
		borderRedis := ""
		borderSQL := ""
		if useImageBorder {
			borderRedis = " style='border: " + borderStyleRedis + "'"
			borderSQL = " style='border:" + borderStyleSQL + "'"
		}
		rdb.Set(ctx, "image:"+id, retImage+borderRedis+">", keyExpiration)
		if debug {
			log.Printf("Got image id %s from SQL", id)
		}
		return retImage + borderSQL + ">"
	}
	if debug {
		log.Printf("Got image id %s from REDIS", id)
	}
	return retImage
}

func getOne(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	id := strings.ReplaceAll(r.URL.Query().Get("id"), "'", "")
	if len(id) == 0 {
		fmt.Fprintf(w, "Missing id! Usage:\n\n")
		fmt.Fprintf(w, " - For one id call -> http://localhost/?id=2874]\n")
		fmt.Fprintf(w, " - For simulate multiple individual id calls -> http://localhost/?id=<3341|2442|1007|...>]\n")
		return
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, "<h3>Time elapsed: <time>...</time>ms</h3>")

	ids := strings.Split(id, "|")
	for i := range ids {
		fmt.Fprintf(w, "%s\n", getImage(ids[i]))
	}
	fmt.Fprintf(w, "<script language='javascript'>document.getElementsByTagName('time')[0].innerText='%d';</script>", time.Since(start).Milliseconds())
}

func main() {
	connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", mssqlServer, mssqlUser, mssqlPassword, mssqlPort)
	if debug {
		fmt.Printf("ConnectionString: %s\n", connString)
	}
	mssqlConnection, errdb = sql.Open("mssql", connString)
	if errdb != nil {
		log.Fatal("Unable to connect to server:", errdb.Error())
	}
	defer mssqlConnection.Close()
	http.HandleFunc("/", getOne)
	http.ListenAndServe(":80", nil)
}
