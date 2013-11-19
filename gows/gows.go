/*
gows.go : GOWS library entry main program
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
package gows
 
import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"os"
	"net"
	"log"
	"os/signal"
	"syscall"
	"net/http"
//	"runtime/debug"
//	"reflect"
	"github.com/jmoiron/sqlx"
	"./auth"
	"./db"
	"database/sql"
	"runtime/pprof"
	"flag"
)

func MsgNew(cmd map[string]interface{})  map[string]interface{} {
	msg:= map[string]interface{}{ "cmd": cmd["cmd"], "ver": cmd["ver"], "seq": cmd["seq"], "stat": 0}
	return msg
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")

func Init() *sqlx.DB {
//	go func() {
//		log.Println(http.ListenAndServe("localhost:6060", nil))
//	}()
	db := db.InitDB("gows","gows","gows")
	HHandle("/auth", auth.AuthHandler)
	HHandle("/login/regst", auth.ReqHandler)
	HHandle("/login/forgot", auth.FgtHandler)

	go h.Run()
	return db
}

func Setup( dir string, hndlr func(http.ResponseWriter, *http.Request)) {
    http.HandleFunc(dir, hndlr)
}


type Command func( conn *Connection, jin map[string]interface{}) map[string]interface{}

func Run(flag *string, errhndlr http.Handler) {
	proto := *flag
	if proto[:4] == "unix" {
		fil := proto[5:]
		err := os.Remove(fil)
		if err != nil {
			fmt.Println("No such file", err)
		}
		addr, err := net. ResolveUnixAddr("unix",fil)
		if err != nil {
			fmt.Println("Addr resolve ", err)
		}
		
		l, err := net.ListenUnix("unix", addr)
		if err != nil {
			log.Fatal("Open Socket ", err)
		}
		
		//lfil, err = l.File()
	
		err = os.Chmod(fil,0666)
		if err != nil {
			log.Fatal("chmod ", err)
		}
		//err  = fil.Chown(33,33)
		//if err != nil {
		//	log.Notice("chown ", err)
		//}
		
		//defer l.Close()
		fmt.Println("set up signal")
		// http://stackoverflow.com/questions/16681944/
			//how-to-reliably-unlink-a-unix-domain-socket-in-go-programming-language
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
		go func(c chan os.Signal) {
			// Wait for a SIGINT or SIGKILL:
			sig := <-c
			log.Printf("Caught signal %s: shutting down.", sig)
			// Stop listening (and unlink the socket if unix type):
			l.Close()
			fmt.Println("SIG")
			// And we're done:
			os.Exit(0)
		}(sigc)
		http.Serve(l, errhndlr)
	} else {
		if err := http.ListenAndServe(*flag, errhndlr); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}
}

func HHandle( dir string, hndlr func(http.ResponseWriter, *http.Request)) {

    http.HandleFunc(dir, hndlr)
}

func WHandle( dir string, hndlr websocket.Handler) {
	if hndlr == nil {
		hndlr = wsHandler
	}
	http.Handle(dir, websocket.Handler(hndlr))
}

func wsHandler(ws *websocket.Conn) {
	        if *cpuprofile != "" {
                f, err := os.Create(*cpuprofile)
                if err != nil {
                        log.Fatal(err)
                }
                pprof.StartCPUProfile(f)
                defer pprof.StopCPUProfile()
        }
	c := mkConn(ws)
	if c == nil {
		fmt.Println("Bad cookie")
		//defer ws.close()
		return
	}
	fmt.Println("WSHAND ",ws.PayloadType)
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader(ws)
	pprof.StopCPUProfile()
}

func echo(conn *Connection, jin map[string]interface{}) map[string]interface{} {
    Broadcast(conn,jin)
    return nil
}

func GetDB() *sqlx.DB {
	return db.GetDB()
}

func ConvNullStr( val sql.NullString) string  {
	return db.ConvNullStr(val)
}
func ConvNullInt( val sql.NullInt64) int64 {
	return db.ConvNullInt(val)
}






