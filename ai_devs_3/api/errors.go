package api

import "log"

func FatalOnError(err error) {
	if err == nil {
		return
	}
	log.Fatalf("%+v", err)
}
