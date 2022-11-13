package main

import (
	"balanceService/pkg/server"
)

func main() {
	server := server.Server{}
	server.Start()
}
