/*
main.go : GOWS main demo program
Copyright (C) 2013 Shaun Savage <savages@savages.com>

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of  MERCHANTABILITY or FITNESS FOR
A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with
this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// rename ?
package main

import (
	"flag"
	//add gows
	"./gows"
	//add modules
	"./mod1"
	"fmt"
//	"reflect"
//	"runtime/debug"
)

// Add Module the modules
//	structure 
//  initalization code
func AddMods() {
	gows.Modules["mod1"] = &mod1.Module

	mod1.Module.Add(&mod1.Module)
}

// service address
var addr = flag.String("addr", ":8088", "http service address")
//var addr = flag.String("addr", "unix:/tmp/gows.sock", "http service address")

// basic test function
//  broadcast to everyone.
func test(conn *gows.Connection, jin map[string]interface{} ) map[string]interface{} {
	ret := gows.MsgNew(jin)
	ret["msg"] = jin["msg"]
	ret["sub"] = jin["sub"]
    gows.Broadcast(conn,ret)
	fmt.Println("main test")
    return nil
}

// main
func main() {
	flag.Parse()
	// init modules
	AddMods()
	// init go framework
	db := gows.Init()
	if db == nil {
		return
	}
 
	// service path
	gows.WHandle("/gows/ws", nil)
	//register "echo" function
	gows.Regfunc("echo", test)
	// different client code
	gows.Regfunc("store", test)

	// Have fun!
	gows.Run(addr, nil)
}


