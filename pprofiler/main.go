package pprofiler

import (
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func Start(seconds time.Duration) {
	go func() {
		logfilename := os.Getenv("CPU_PROFILE_FILE")
		if logfilename != "" {
			logfile, err := os.Create(logfilename)
			if err == nil {
				log.Printf("CPU Profiler start (%dsec): %s\n", seconds, logfilename)
				defer logfile.Close()
				pprof.StartCPUProfile(logfile)
				defer pprof.StopCPUProfile()
				time.Sleep(seconds * time.Second)
				log.Println("CPU Profiler finish", logfilename)
			} else {
				log.Println("CPU Profiler error: ", err)
			}
		}
	}()
}
