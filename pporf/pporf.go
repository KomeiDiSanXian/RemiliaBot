package pporf

import (
	"log"
	"net/http"

	_ "net/http/pprof"
)

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
