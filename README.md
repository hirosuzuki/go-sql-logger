# ISUCON Snippets

## Import SQL Logger Module

```go
import (
	_ "github.com/go-sql-driver/mysql"
	goSqlLogger "github.com/hirosuzuki/go-sql-logger"
   	"github.com/hirosuzuki/go-sql-logger/measure"
	"github.com/hirosuzuki/go-sql-logger/pprofiler"
)
```

## Initalize Handler with Go Prof

```go
func initializeHandler() {
	pprofiler.Start(70)
}
```

## Open MySQL Connection using sqlx

```go
goSqlLogger.RegisterTraceDBDriver()
sqlx.Open("mysql"+os.Getenv("MYSQL_DRIVER_POSTFIX"), dsn)
```

## Service Start Script

```sh:env.sh
#!/bin/sh
export MYSQL_DRIVER_POSTFIX=:logger
export SQL_LOGFILE=/tmp/sql.log
export CPU_PROFILE_FILE=/tmp/cpu.pprof
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

## MySQL Slow Query

```
SET GLOBAL long_query_time = 0;
SET GLOBAL slow_query_log = ON;
SET GLOBAL slow_query_log_file = "/tmp/slow.log";
```

