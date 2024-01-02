// Package pprof 用于性能分析
package pprof

import (
	"log"
	"net/http"

	_ "net/http/pprof" // pprof 必须import
)

// Start 启动一个pprof服务
func Start() {
	log.Println("Starting pporf")
	go func() {
		err := http.ListenAndServe("127.0.0.1:6060", nil)
		if err != nil {
			log.Println("pporf error: ", err)
			return
		}
	}()
}
