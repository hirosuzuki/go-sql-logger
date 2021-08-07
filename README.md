# ISUCON Snippets

## Import SQL Logger Module

```go
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/hirosuzuki/go-sql-logger"
   	"github.com/hirosuzuki/go-sql-logger/measure"
)
```

## Initalize Handler with Go Prof

```go
func initializeHandler() {
	go func() {
		logfilename := os.Getenv("CPU_PROFILE_FILE")
		if logfilename != "" {
			logfile, _ := os.Create(logfilename)
			defer logfile.Close()
			pprof.StartCPUProfile(logfile)
			defer pprof.StopCPUProfile()
			time.Sleep(70 * time.Second)
		}
	}()
}
```

## Open MySQL Connection using sqlx

```go
sqlx.Open("mysql"+os.Getenv("MYSQL_DRIVER_POSTFIX"), dsn)
```

## Service Start Script

```sh:start.sh
#!/bin/sh
export MYSQL_DRIVER_POSTFIX=:logger
export SQL_LOGFILE=/tmp/sql.log
export CPU_PROFILE_FILE=/tmp/cpu.pprof
exec /home/isucon/isuumo/webapp/go/isuumo >/dev/null 2>/dev/null
```

## Nginx Config

```nginx
log_format with_time '$remote_addr $uid_got $cookie_user [$time_local] '
    '"$request" $status $body_bytes_sent '
    '"$http_referer" "$http_user_agent" $request_time';

access_log /var/log/nginx/access.log with_time;
```

## Measure

```
defer measure.Start("APIPaymentToken").Stop()
```