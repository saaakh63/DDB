package main

import (
  "bufio"
  "database/sql"
  "fmt"
  "log"
  "net"
  "os"
  "strings"
  "sync"

  _ "github.com/go-sql-driver/mysql"
)

type SnapClient struct {
  conn net.Conn
  id   string
}

var (
  snaps     = make(map[string]SnapClient)
  snapsLock = sync.Mutex{}
  db        *sql.DB
)

func handleSnap(conn net.Conn) {
  defer conn.Close()
  snapID := conn.RemoteAddr().String()

  snapsLock.Lock()
  snaps[snapID] = SnapClient{conn, snapID}
  snapsLock.Unlock()

  log.Printf("Snap connected: %s\n", snapID)

  reader := bufio.NewReader(conn)
  for {
    msg, err := reader.ReadString('\n')
    if err != nil {
      log.Printf("Snap disconnected: %s\n", snapID)
      snapsLock.Lock()
      delete(snaps, snapID)
      snapsLock.Unlock()
      return
    }
    msg = strings.TrimSpace(msg)
    log.Printf("Received from %s: %s\n", snapID, msg)

    executeSQL(msg)                // Master executes
    broadcastToOthers(snapID, msg) // Forward to other Snaps
  }
}

func broadcastToOthers(senderID, cmd string) {
  snapsLock.Lock()
  defer snapsLock.Unlock()
  for id, snap := range snaps {
    if id != senderID {
      _, err := snap.conn.Write([]byte(cmd + "\n"))
      if err != nil {
        log.Printf("Error forwarding to %s: %v\n", id, err)
      }
    }
  }
}

func executeSQL(cmd string) {
  _, err := db.Exec(cmd)
  if err != nil {
    log.Printf("SQL Error: %v\n", err)
  }
}

func terminalControl() {
  scanner := bufio.NewScanner(os.Stdin)
  for {
    fmt.Print("Master> ")
    scanner.Scan()
    cmd := strings.TrimSpace(scanner.Text())
    if strings.HasPrefix(strings.ToLower(cmd), "show snaps") {
      snapsLock.Lock()
      for id := range snaps {
        fmt.Println("Connected Snap:", id)
      }
      snapsLock.Unlock()
      continue
    }
    executeSQL(cmd)
    broadcastToOthers("MASTER", cmd) // "MASTER" acts as pseudo ID
  }
}

func main() {
  var err error
  db, err = sql.Open("mysql", "shery:password@tcp(127.0.0.1:3306)/") // Master MySQL
  if err != nil {
    log.Fatal(err)
  }

  ln, err := net.Listen("tcp", ":8080")
  if err != nil {
    log.Fatal(err)
  }
  defer ln.Close()

  go terminalControl()

  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Println("Error accepting connection:", err)
      continue
    }
    go handleSnap(conn)
  }
}
