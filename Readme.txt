This tool can be used to send a file over network sockets to
a specified target (for example a payload)...

Once the payload is sent, the tool listens for connections o two
port (usually 30 and 31).

Everything received on port 30, is printed to the console.
Everything received on port 31, is written/appended to a file called dump.bin
