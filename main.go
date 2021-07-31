package sqllogger

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	proxy "github.com/shogo82148/go-sql-proxy"
)

var sqlLogFileName string
var sqlLogFile *os.File

func init() {
	sqlLogFileName = os.Getenv("SQL_LOGFILE")
	if sqlLogFileName == "" {
		sqlLogFileName = "/tmp/sql.log"
	}

	log.Printf("Go SQL Logger: Log File -> %s\n", sqlLogFileName)
	var err error
	if sqlLogFile, err = os.OpenFile(sqlLogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		log.Printf("Go SQL Logger: Open File Error: %s\n", err.Error())
		return
	}
	registerTraceDBDriver()
}

func registerTraceDBDriver() {
	regexCutSpace := regexp.MustCompile(`[ \r\n\t]{1,}`)
	regexTagComment := regexp.MustCompile(`(/\* *(.*?) *\*/)`)

	PreFunc := func(c context.Context, stmt *proxy.Stmt, args []driver.NamedValue) (interface{}, error) {
		return time.Now().UnixNano(), nil
	}
	PostFunc := func(c context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, err error) error {
		if sqlLogFile != nil && err != driver.ErrSkip {
			now := time.Now()
			startTime := ctx.(int64)
			timeDelta := now.UnixNano() - startTime
			query := regexCutSpace.ReplaceAllString(stmt.QueryString, " ")
			posList := regexTagComment.FindStringSubmatchIndex(query)
			tag := ""
			if posList != nil {
				tag = query[posList[4]:posList[5]]
				query = query[:posList[1]]
			}
			fmt.Fprintf(sqlLogFile, "%d\t%d\t%s\t%s\n", startTime, timeDelta, tag, query)
		}
		return nil
	}

	for _, driverName := range sql.Drivers() {
		if strings.Contains(driverName, ":logger") {
			continue
		}
		db, _ := sql.Open(driverName, "")
		defer db.Close()
		newDriverName := driverName + ":logger"
		log.Printf("Go SQL Logger: SQL Driver Regist -> %s\n", newDriverName)
		sql.Register(driverName+":logger", proxy.NewProxyContext(db.Driver(), &proxy.HooksContext{
			PreExec: PreFunc,
			PostExec: func(c context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, result driver.Result, err error) error {
				return PostFunc(c, ctx, stmt, args, err)
			},
			PreQuery: PreFunc,
			PostQuery: func(c context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, rows driver.Rows, err error) error {
				return PostFunc(c, ctx, stmt, args, err)
			},
		}))
	}
}
