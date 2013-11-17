<!--
hub.go : GOWS  program the routing
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
import	(
    "encoding/json"
    "fmt"
//    "debug"
)

//import "errors"

type BaseReq struct {
	cmd string
	ver string
	seq	uint64
//	interface {}
}

type BaseRet struct {
	BaseReq
	stat int16
//	interface {}
}

type Cmditem struct {
    cmd string
    jout map[string]interface{}
}

type Cmdfunc struct {
    jout map[string]interface{}
    //Cmd
    //func(*hub,map[string]interface{}) map[string]interface{}
}


//ype CmdItem interface


/*
func (ci *Cmditem) test(h *Hub, jin map[string]interface{}, conn *Connection) map[string]interface{} {
    broadcast(h,jin)
    return nil
}

func (ci *Cmditem) test1(h *Hub, jin map[string]interface{}, conn *Connection) map[string]interface{} {
    broadcast(h,jin)
    return nil
}
*/

func Regfunc(nam string, cmd Command) {
    cmdtbl[nam] = cmd
}
 
var cmdtbl =  make(map[string]Command )


type hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	//Command chan map[string]interface{}
    command chan *Cmdstruct

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection
}

var h = hub{
	//command:   make(chan map[string]interface{}),
	command:   make(chan *Cmdstruct),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
}

func (h *hub) errsnd(conn *Connection, msg string) {
    var emsg map[string]interface{}
    str := []byte("{\"msg\":\"" + msg + "\",\"cmd\":\"error\", \"errval\":" + string("23") + "}")
    eerr := json.Unmarshal(str, &emsg)
    fmt.Println(eerr,msg, string(str));
    fmt.Println(len(cmdtbl))
     if ( eerr == nil ) {
        conn.Send <- emsg 
    }
}

func (h *hub) Run(/*pool *pgsql.Pool*/) {
	for {
		select {
		case c := <-h.register:
/*            conn, err := pool.Acquire()
            if err != nil {
                log.Printf("Error acquiring Connection: %s\n", err)
            } else {
                res, err := conn.Query("SELECT now()")
                if err != nil {
                        log.Printf("Error executing query: %s\n", err)
                } else {
                    if !res.Next() {
                        log.Println("Couldn't advance result cursor")
                    } else {
                        var now string
                        if err := res.Scan(&now); err != nil {
                            log.Println("scan %s",err)
                        } else {
                            fmt.Println("ts %s", now)
                        }
                    }
                }
            }*/
			h.connections[c] = true
			fmt.Println("connect  ",c)
		    //debug.PrintStack()

            //pool.Release(conn)
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.Send)
			fmt.Println("disconnect  ")
		    //debug.PrintStack()
		case m := <-h.command:
            var conn *Connection
            conn = m.c
            var msg map[string]interface{}
            msg = m.msg
            cmd := msg["cmd"].(string)
			fmt.Println("MSG",cmd,msg)
 		    //debug.PrintStack()

            f := cmdtbl[cmd]
            r := recover()
            if f == nil {
                //var emsg map[string]interface{}
                str := string("ERROR: cmd "+cmd+" not found.")
                fmt.Println(r, len(cmdtbl), str);
                h.errsnd(conn, str);
                continue
            }
            f(conn,msg)   //go
		}
	}
}

func Broadcast(conn *Connection, dat map[string]interface{}) {
	fmt.Println("BCAST",conn,dat)
    for c := range conn.connections {
		fmt.Println("BC0",)
        func(conn *Connection, dat map[string]interface{}) { //go
            select {
            case c.Send <- dat:
            default:
				fmt.Println("close",conn)
                delete(conn.connections, c)
                close(c.Send)
                go c.Ws.Close()
            }
        } (conn,dat)
    }
}

