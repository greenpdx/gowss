/*
main.go : GOWS main program
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
package gows
 
import (
	"fmt"
//    "strings"
    "reflect"
)

var kindStr = []string{
	"Invalid",
	"Bool",
	"Int",
	"Int8",
	"Int16",
	"Int32",
	"Int64",
	"Uint",
	"Uint8",
	"Uint16",
	"Uint32",
	"Uint64",
	"Uintptr",
	"Float32",
	"Float64",
	"Complex64",
	"Complex128",
	"Array",
	"Chan",
	"Func",
	"Interface",
	"Map",
	"Ptr",
	"String",
	"Struct",
	"UnsafePointer",
}

func DeReflect( unk interface{} ) {
	typ := reflect.TypeOf(unk)
	k := typ.Kind()
	ps := reflect.ValueOf(&unk)
	s := ps.Elem()
	ks := kindStr[s.Kind()]
	//nt := reflect.TypeOf(typ.(book.gbRet))
	nt := typ.PkgPath()
	fmt.Println(ps, typ, s, k, ks, nt)
	
	switch (k) {
	case reflect.Bool:
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
	case reflect.Uintptr:
	case reflect.Float32:
	case reflect.Float64:
	case reflect.Complex64:
	case reflect.Complex128:
	case reflect.Chan:
	case reflect.Func:
	case reflect.Array:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.String:
	case reflect.UnsafePointer:
	default:
		break;
	case reflect.Struct:
		deStruct(typ, ps)
		break;
	}
}

func deStruct (typ reflect.Type, ps reflect.Value ) {
	fcnt := typ.NumField()
	for i := 0; i < fcnt; i++ {
		ftyp :=  typ.Field(i)
		//f := ps.FieldByName(ftyp.Name)
		//v := reflect.ValueOf(&f).Elem()
		fmt.Println(ftyp)
	}

}
