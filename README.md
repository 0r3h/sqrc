# SQRC

Rudimentary Squad RCON CLI. Reads `stdin` for commands, outputs command responses and chat to `stdout`.

## Installation

Either use release binary or compile from source.

    go get github.com/0r3h/sqrc

## Launching

Pass Squad server IP, RCON port and password as arguments when executing binary.

    sqrc <ip:port> <password>