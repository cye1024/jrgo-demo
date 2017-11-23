package main

import "jrgo-demo/server"

func main() {
	server.ServerHTTP()
	select {}
}
