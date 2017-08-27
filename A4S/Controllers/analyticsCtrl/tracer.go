package analytics

import (
	"log"
	"os"
	"runtime/trace"
)

func MonitorTracer() {
	/*
					- net: network blocking profile
					- sync: synchronization blocking profile
					- syscall: syscall blocking profile
					- sched: scheduler latency profile

		Generate a pprof-like profile from the trace:
				go tool trace -pprof=TYPE trace.out > TYPE.pprof
		Then, you can use the pprof tool to analyze the profile:
				go tool pprof TYPE.pprof
	*/
	log.Println(`Monitor Tracer on ... 📈`)

	f, err := os.Create("./Analytics/TracerFiles/trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
}
