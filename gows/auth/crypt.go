/*
crypt.go : GOWS crytography program
Copyright (C) 2013 Shaun Savage <savages@savages.com>

This program is free software: you can redistribute it and/or modify it under
the terms of the GNU Leser General Public License as published by the Free Software
Foundation, either version 3 of the License, or (at your option) any later
version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of  MERCHANTABILITY or FITNESS FOR
A PARTICULAR PURPOSE. See the GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License along with
this program.  If not, see <http://www.gnu.org/licenses/>.
*/
package auth

import (
//	"log"
	"fmt"
    "github.com/gorilla/securecookie"
    "encoding/base64"
	"encoding/binary"
	"hash/crc64"
	"errors"
	"reflect"
)

type Cryptcfg struct {
	authKey		[]byte
	cryptKey	[]byte
	crcKey		[]byte
	digKey		[]byte
	hashKey		[]byte
}
var cryptgows Cryptcfg
var crctbl *crc64.Table

func EnCrypt(vals *map[interface{}]interface{}, keys *Cryptcfg ) ([]string, error) {
	
	scook := securecookie.New(
		keys.authKey,
		keys.cryptKey)
	str,err := scook.Encode("gows",vals)
	if err != nil {
		fmt.Println("Bad encode ",err)
		return nil, err
	}
	fmt.Println(str)	
	cs := crc64.Checksum([]byte(str+string(keys.crcKey)),crctbl)
	bcs := make([]byte, 16)
	len := binary.PutUvarint(bcs,cs)

	crcstr := base64.StdEncoding.EncodeToString(bcs[:len])
	return []string{str, crcstr}, nil
}

func DeCrypt(str []string, keys *Cryptcfg) (*map[interface{}]interface{}, error) {
    vals := make(map[interface{}]interface{})
	crctbl := crc64.MakeTable(crc64.ECMA)
	crcsecret := keys.crcKey
	cs := crc64.Checksum([]byte(str[1]+string(crcsecret)),crctbl)
	bcs := make([]byte, 16)
	slen := binary.PutUvarint(bcs,cs)
	
	crcstr := base64.StdEncoding.EncodeToString(bcs[:slen])
	//fmt.Println("COOK",cs,bcs[:slen],reflect.TypeOf(crcstr),len(crcstr),reflect.TypeOf(str[0]),crcstr,len(str[0]),str[0])
	fmt.Println("COOK",cs,bcs[:slen],crcstr,str[0])
	if str[0] != string(crcstr) { //invalid cookie
		fmt.Println(crcstr, str[0], reflect.TypeOf(crcstr),len(crcstr),reflect.TypeOf(str[0]),len(str[0]))
		fmt.Println("bad cookie")
		return nil, errors.New("Invalid Cookie")
	}
	scook := securecookie.New(
		keys.authKey,
		keys.cryptKey)
	fmt.Println(str[1], vals, reflect.TypeOf(str[1]),reflect.TypeOf(vals))
	err := scook.Decode("gows",str[1],&vals)
	if err != nil {
		fmt.Println("bad cookie crypt",err)
		return nil, errors.New("Invalid Cookie")		
	}
	return &vals,nil
}


func DeCook(str []string) (*map[interface{}]interface{}, error) {
	return DeCrypt(str, &cryptgows)
}

func EnCook(vals *map[interface{}]interface{}) ([]string, error) {
	return EnCrypt(vals, &cryptgows)
}
