package main

import (
	"io"
	"log"
	"os"
	"wallpaper-filter/router"
)

func SetLog() {
	f, err := os.OpenFile("./logs/wallpaper-filter.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	multiWriter := io.MultiWriter(f, os.Stdout)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	SetLog()

	router.InitRouter()
}
