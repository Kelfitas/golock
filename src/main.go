package main

import "flag"

func main() {
	flag.Parse()
	authDone := make(chan bool)
	go watchAuth(authDone)
	go startApp(authDone)

	<-authDone
}
