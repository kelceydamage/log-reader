// Package main is a server for emitting new log lines over a socket.
package main
    
import (
    "fmt"
    "net"
    "os"

    "github.com/hpcloud/tail"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

var (
    ioconfig = tail.Config{
        Follow: true,
        Poll: false,
        ReOpen: true,
    }
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    check(err)
    
    defer l.Close()
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        go handleRequest(conn)
    }
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
}

func handleRequest(conn net.Conn) {
  t, err := tail.TailFile("metrics.log", ioconfig)
  check(err)

  // Send the lines from the log and stay open for new lines.
  for line := range t.Lines {
    fmt.Println(line.Text)
    conn.Write([]byte(line.Text + "\n"))
  }

  conn.Close()
}