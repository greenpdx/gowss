conn.go : GOWS connection 
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
	"code.google.com/p/go.net/websocket"
	"fmt"
	"encoding/json"
    "strings"
//    "runtime/debug"
	"./auth"
	"runtime/pprof"
	"os"
	"log"
)

type T struct {
    Msg string
    Count int
}

type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn
 
	// Buffered channel of outbound messages.
	Send chan map[string]interface{}
	hub
	Msg map[string]interface{}
}

func (c *Connection) SendMsg(msg map[string]interface{}) {
	c.Send <- msg
}

func (c *Connection) Reply( dat map[string]interface{}) {
	c.Send <- dat
}
 
type Cmdstruct struct {

    c *Connection
    msg map[string]interface{}
 }

func (c *Connection) reader(ws *websocket.Conn) {
	i := 0
	for {
		var data []byte
        var dat map[string]interface{}
        err := websocket.Message.Receive(c.Ws, &data)
//        err := websocket.JSON.Receive(c.Ws, &data)
        if ( err != nil) {
            fmt.Println("ws error ",err, )
            break
        }

        fmt.Println(i,  string(data))
		i = i + 1
			dec := json.NewDecoder(strings.NewReader(string(data)))
			fmt.Println("x1",dec)
			decerr := dec.Decode(&dat)
			fmt.Println("x2",decerr, dat)
			//err = websocket.JSON.Receive(c.Ws, &dat)
			fmt.Println("next")
			if decerr != nil {
				fmt.Println("JSON BAD ",decerr)
				break
			}
			h.command <-  &Cmdstruct{c:c,msg:dat}
    }
//	c.Ws.Close()
}
 
func (c *Connection) writer() {
	for dat := range c.Send {
		fmt.Println(dat)
		err := websocket.JSON.Send(c.Ws, dat)
		if err != nil {
            fmt.Println(err)
			break
		}
	}
//	c.Ws.Close()
}

// (c *Connection) Send () {
//}

/*
 * type Args struct {
 *      A int
 *      B int
 * }
 * 
 * func (t *Arith) Multiplt( args *Args, reply *int) error {
 *      *reply = args.A * Args.B
 *      return nil
 * }
*/


func mkConn(ws *websocket.Conn) *Connection {
        if *cpuprofile != "" {
                f, err := os.Create(*cpuprofile)
                if err != nil {
                        log.Fatal(err)
                }
                pprof.StartCPUProfile(f)
                defer pprof.StopCPUProfile()
        }
//conf := ws.Config()
	req := ws.Request()
	cookie, err := req.Cookie("hzc")
	if err != nil {
		fmt.Println("bad cookie")
		return nil
	}
	
	encstr := strings.Split(cookie.String()[5:],":")
	
	vals, err := auth.DeCook(encstr)
	if err != nil {
	}
	fmt.Println(vals)

	c := &Connection{ws, make(chan map[string]interface{}, 256), h, nil}
	//fmt.Println(cookie,c,conf,req)
	return c
}

