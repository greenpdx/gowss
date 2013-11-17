main.go : GOWS main program
Copyright (C) 2013 Shaun Savage <savages@savages.com>

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Lesser General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of  MERCHANTABILITY or FITNESS FOR
A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License along with
this program.  If not, see <http://www.gnu.org/licenses/>.
-->
package gows
 
import (
	"fmt"
)

import _ "github.com/jbarham/gopgsqldriver"
type Module struct {
	Name string
	Procs map[string]Command
	Data interface{}
}

var modules = make(map[string]*Module)
var modlst = [16]string{}

func (m Module) Add (mod *Module) {
	Modules[mod.Name] = mod
	fmt.Println("Mod add",mod.Name,len(mod.Procs))
	for nam, proc := range mod.Procs {
		fmt.Println(nam)
		cmdtbl[nam] = proc
	}
	m.List()
}
/*
func ddmods() {
	modules["stroke"] = &stroke.Module
	modules["book"] = &book.Module
	modules["edit"] = &edit.Module

	stroke.Module.Add(&stroke.Module)
	book.Module.Add(&book.Module)
	edit.Module.Add(&edit.Module)
}
*/
//func Del
func (m Module)List() {
	fmt.Println(len(cmdtbl))
//	for k,v := range cmdtbl {
//	}
}

var Modules = make(map[string]*Module)
ZZ
