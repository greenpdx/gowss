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
package mod1

import (
	"../gows"
	"fmt"
//	"reflect"
//	"runtime/debug"
)


type gbReq struct {
	gows.BaseReq
	btype uint8
}

var Module gows.Module

func mod1(conn *gows.Connection, cmd map[string]interface{} ) map[string]interface{} {
	fmt.Println("mod1")

	msg := gows.MsgNew(cmd)

	//db := gows.GetDB()


	conn.Send <- msg
	return msg
}


func init() {
	mod := gows.Module{
		Name:	"mod1",
		Data:	nil,
		Procs:  make(map[string]gows.Command),
	}
	mod.Procs["mod1"] = mod1

	Module = mod
	fmt.Println("mod1 init",mod.Name, mod.Procs)
	
}
