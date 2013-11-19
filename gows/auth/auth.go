/*
auth.go : GOWS authencation program
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
*/
package auth

import (
//    "github.com/jmoiron/sqlx"
	"log"
	"fmt"
	"github.com/gorilla/sessions"
    "github.com/gorilla/securecookie"
    "net/http"
    "strconv"
    "time"
//    "strings"
//    "encoding/base64"
//	"encoding/binary"
	"hash/crc64"
//	"errors"
	"text/template"
	"../db"
	"reflect"
//	"crypto/x509"
//	"crypto/rsa"
//	"encoding/asn1"
//	"crypto/x509/pkix"
)




var redirTempl *template.Template
var regstTempl *template.Template
var store = sessions.NewCookieStore([]byte("AuthKey"),[]byte("CryptKey"))



type usrInfo struct {
	uid		uint32
	usrloc	string
	name	string
	login	string
	keys	usrsec
}

type usrsec struct {
	Userkey	string
	Sesskey	string
	NextWin string
}

func init() {
	crctbl = crc64.MakeTable(crc64.ECMA)
	cryptgows.cryptKey = securecookie.GenerateRandomKey(32)
	cryptgows.hashKey = securecookie.GenerateRandomKey(32)
	cryptgows.authKey = securecookie.GenerateRandomKey(32)
	cryptgows.crcKey = securecookie.GenerateRandomKey(32)
	redirTempl = template.Must(template.ParseFiles("gows/redir.html"))
	regstTempl = template.Must(template.ParseFiles("gows/regst.html"))
}
	

func validate(u string, p string) *usrInfo {
	fmt.Println(u,p)
	var usrinfo usrInfo
	
	usrinfo.uid = 1
	usrinfo.usrloc = "LOCALKEY-1"
	usrinfo.name = "tstname"
	usrinfo.login = "tstlogin"
	usrinfo.keys.Userkey = "USERKEY"
	
	if p == "test" {
		return &usrinfo
	}
	return nil
}



func FgtHandler(w http.ResponseWriter, r *http.Request) {
	var us usrsec
	us.Userkey = "USERKEY"

	qry := r.URL.RawQuery
	fmt.Println("Forgot Req",qry,r,w)
	if len(qry) == 0 {  // no parameters		
		http.Redirect(w,r,"/login",307)
		return
	}
	val := r.URL.Query()
	fmt.Println(val)
}

func ReqHandler(w http.ResponseWriter, r *http.Request) {
	var us usrsec
	us.Userkey = "USERKEY"
	us.NextWin = "/login"

	err := r.ParseForm()
	val := r.Form

	fmt.Println("Regst Req",val)
	if len(r.URL.RawQuery) == 0 {  // no parameters		
		http.Redirect(w,r,"/login",307)
		//return
	}
	
/*	if (val["oldloc"] != nil) {
		fmt.Println(val["oldloc"][0])
		return		// later add crypto
	}
	if val["keygen"] != nil {
		_, err := base64.StdEncoding.DecodeString(val["keygen"][0]) //[]byte
		if err !=  nil {
			fmt.Println("bad b64", err)
			return
		}
	}
*/	
	if (val["pass1"][0] != val["pass2"][0] ) {
		fmt.Println("pass not match")
		return
	}
	
//	if len(val["pass1"][0]) < 6 || len(val["rname"][0]) > 32 || len(val["pass1"][0]) > 32 || len(val["remail"][0]) > 32 || len(val["rlogin"][0]) >20 {
//		return
//	}
	
	db := db.GetDB()		
	stmt, err := db.Preparex("SELECT useradd($1,$2,$3,$4)")
	if err != nil {
		fmt.Println("bad prepare ", err)
		http.Redirect(w,r,"/login",307)
		//regstTempl.Execute(w, us)
		return
	}
	type lastid struct {
		Useradd		int64
	}
	idx := lastid{}
	err = stmt.Get(&idx,val["rname"][0], val["rlogin"][0], val["remail"][0], val["pass1"][0])
	if err != nil {
		fmt.Println("bad user add ", err)
		http.Redirect(w,r,"/login",307)
		//regstTempl.Execute(w, us)
		// send something back stating choose different login
		return
	}
	
	usridx := idx.Useradd
	fmt.Println(usridx)
	
	vals := map[interface{}]interface{}{}
	vals["sidx"] = idx.Useradd
	str, err := EnCook(&vals)
	if err != nil {
		fmt.Println("bad crypt ", err)
		http.Redirect(w,r,"/login",307)		
	}

	us.Userkey = str[0]
	us.NextWin = "/login"
	
	regstTempl.Execute(w, us)
	return
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	val := r.Form

	fmt.Println("Auth Req",val)
	if len(r.URL.RawQuery) == 0 {  // no parameters		
		log.Printf("No Query: %s\n", err)
		http.Redirect(w,r,"/login",307)
		//return
	}
	vals := map[interface{}]interface{}{}
	
	ip := r.Header.Get("X-Real-Ip")
	port := r.Header.Get("X-Real-Port")
	fmt.Println(val,val["username"])
	usrinfo := validate(val["username"][0], val["password"][0])
	
	if usrinfo == nil {
		http.Redirect(w,r,"/login",307)
		fmt.Println("No User")
		return
	}
	uid := usrinfo.uid
	vals["uid"] = strconv.FormatUint(uint64(uid),10)
	
	suid := strconv.FormatUint(uint64(uid), 36)
	
	tim := time.Now().UnixNano()
	ctim := tim + (1000000*60*15)
	stim := strconv.FormatUint(uint64(ctim), 10)
	
	fmt.Println(suid,reflect.TypeOf(suid), reflect.TypeOf(r))
	ses, err := store.Get(r,suid)	
	if err != nil {
		log.Printf("Error new session: %s\n", err)
		http.Redirect(w,r,"/login",307)
		return
    }
	if ses.IsNew {
		//ses.Options.Domain = r.Host
		ses.Options.Path = "/gows"
		ses.Options.HttpOnly = false
		ses.Options.Secure = false // change to true for production
		ses.Values = vals
	} else {  // reconnect
		vals = ses.Values
	}
	ses.Options.MaxAge = 60*30 //change to config
	
	db := db.GetDB()	//db = gows.GetDB()		
	
	type sessidx struct {
		SessIdx uint64
	}
	si := sessidx{}
	err = db.Get(&si,"INSERT INTO sess (uid,ip,svalid,port,intval) VALUES ($1,$2,true,$3,$4) RETURNING sessidx",uid,ip,port, 60)
	if err != nil {
		log.Printf("Error new session: %s\n", err)
		http.Redirect(w,r,"/login",307)
		return
	}
	sessionidx := si.SessIdx
	fmt.Println("SI >",si)
	vals["sidx"] = strconv.FormatUint(uint64(si.SessIdx),36)
	
	str, err := EnCook(&vals)
	if err != nil {
		log.Printf("New cookie err: %s\n", err)
		http.Redirect(w,r,"/login",307)
		return
	}

	cook := sessions.NewCookie("gows",str[1]+":"+str[0], ses.Options)
	fmt.Println("COOKIE SEND",cook)	
	http.SetCookie(w,cook)
	
	fmt.Println(uid, sessionidx, stim,vals, str)
	ses.Values = vals
	store.Save(r, w, ses)
	
	fmt.Println("EXIT")
	var us usrsec
	//us.Userkey = "USERKEY"
	us.Sesskey = str[0]
	us.NextWin = "/gows"
	redirTempl.Execute(w, us)
	return
}

