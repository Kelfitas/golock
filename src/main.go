package main

func main() {
	authDone := make(chan bool)
	go watchAuth(authDone)
	go startApp(authDone)

	<-authDone
}
