package main

import (
	routerCtrl "github.com/tknott95/LCW-FC/Controllers"
	httpGlobals "github.com/tknott95/LCW-FC/Globals/http"
)

func main() {
	println("\n || LCW Golang Backend Server - TK || \n \n")
	println("\n🚀 Server Running on Port: " + httpGlobals.PortNumber + "\n")

	routerCtrl.InitServer()
}
