package go_sql_logger

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hirosuzuki/go-sql-logger/logger"
	proxy "github.com/shogo82148/go-sql-proxy"
)

func init() {
	// RegisterTraceDBDriver()
}

func RegisterTraceDBDriver() {
	regexCutSpace := regexp.MustCompile(`[ \r\n\t]{1,}`)
	regexTagComment := regexp.MustCompile(`(/\* *(.*?) *\*/)`)

	PreFunc := func(c context.Context, stmt *proxy.Stmt, args []driver.NamedValue) (interface{}, error) {
		return time.Now().UnixNano(), nil
	}
	PostFunc := func(c context.Context, ctx interface{}, stmt *proxy.Stmt, args []driver.NamedValue, err error) error {
		if err != driver.ErrSkip {
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
			if len(query) > 1000 {
				query = query[:1000]
			}
			logger.Log(startTime, timeDelta, tag, query)
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
