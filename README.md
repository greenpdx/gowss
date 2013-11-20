GOWSS
====

GO language WebSocket Web Application Framework

This is a websocket framework designed to work with nginx. It servers 
static pages/web apps and MVC data.  It is designed to support Postgres
Database.  New packages can be written and added to the main program.

The server requires nginx web server. The earlier routing and 
authorization the quicker the response, less server work.

The server request are sorted at the web server layer.
request 
"/login"	show login in page
"/gows/ws"	if closed, reopen
					if cookie expired -> /login
					else process command.
"/gows"		if cookie expired -> /login 
			else serve static web application
			
During login a crypto cookie is generated.



I have written other MVC frameworks but this is my first attempt
to learn and write GO.  
