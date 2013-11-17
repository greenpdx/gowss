/*
mWeb.js : GOWS network JavaScript interface demo program
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
var wsconn = null;
ZZ
function maWeb(url, readyCB) {
	this.url = url;
	this._conn = null;
	this.Ready = false;
	self = this;
	//console.trace()
	if (readyCB == 'undefined' || readyCB == null ) {
		alert("Need CallBack");
		return;
	}
	// connect and reconnect
	this.connect = function(cb, data, msgcb) {
		if ( this._conn == null) {
			if (window["WebSocket"]) {
				self.cb2 = msgcb;
				self._conn = new WebSocket("ws://"+location.host+"/gofw/ws");
				//console.log("new sock")
				var state = this._conn.readyState;
				var buflen = this._conn.bufferedAmount;
				var ext = this._conn.extensions;
				var proto = this._conn.protocol;
				var bt = this._conn.binaryType;
        
				self._conn.onclose = function(evt) {
					var nm = wsconn;
					//var nm.Ready = false;
					wsconn.Ready = false;
					wsconn._conn = null;
				}
				self._conn.onmessage = function(evt) {
					var jmsg = JSON.parse(evt.data);
				}
				self._conn.onerror = function(er) {
					alert("ERROR " + er);
				}
				self._conn.onopen = function(evt) {
					self.Ready = true;
					cb(self, data, msgcb); 
				}
			} else {
			;
			}
		};
	};

	
	this.connect(readyCB, null)

	function _send(self, data, msgcb) {
		self._conn.onmessage = function( evt ) {
			var r = evt.target;
			var str = evt.data;
			try {
				var obj = JSON.parse( str );
			} catch (e) {
				if ( str.length > 0 )
					alert( e+" >"+str+"<" );
				throw eval(e+" >"+str+"<")  ;
			}
			msgcb( evt, obj );
			return false;
		};
		try {
			self._conn.send( data );
		} catch (e) {
			alert( e );
		}
	}
 
    this.send = function( data, msgcb ) {
		if (!this.Ready) {
			this.connect(_send, data, msgcb);
			return;
		}
		_send(this, data, msgcb)
	};

	this.mkCmd = function(cmdstr) {
		var cmd = {};
		cmd['ver'] = 0.01;
		cmd['seq'] = 12345;
		cmd['cmd'] = cmdstr;
		return cmd;
	};
 };
