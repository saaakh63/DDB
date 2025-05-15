package main

import (
  "bufio"
  "database/sql"
  "fmt"
  "log"
  "net"
  "os"
  "strings"

  _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectToMaster(masterIP string) net.Conn {
  conn, err := net.Dial("tcp", masterIP+":8080")
  if err != nil {
    log.Fatalf("Failed to connect to Master: %v", err)
  }
  log.Printf("Connected to Master at %s\n", masterIP)
  return conn
}

func listenForMaster(conn net.Conn) {
  reader := bufio.NewReader(conn)
  for {
    cmd, err := reader.ReadString('\n')
    if err != nil {
      log.Println("Lost connection to Master.")
      return
    }
    cmd = strings.TrimSpace(cmd)
    log.Printf("Command from Master: %s\n", cmd)
    executeLocalSQL(cmd)
  }
}

func executeLocalSQL(cmd string) {
  _, err := db.Exec(cmd)
  if err != nil {
    log.Printf("Local SQL error: %v\n", err)
  }
}

func snapTerminal(conn net.Conn) {
  scanner := bufio.NewScanner(os.Stdin)
  for {
    fmt.Print("Snap> ")
    scanner.Scan()
    cmd := strings.TrimSpace(scanner.Text())

    // Restrict schema-altering commands
    restrictedCommands := []string{
      "create database", "drop database",
      "create table", "drop table",
      "alter table", "truncate table",
    }

    lowerCmd := strings.ToLower(cmd)
    for _, restricted := range restrictedCommands {
      if strings.HasPrefix(lowerCmd, restricted) {
        fmt.Println("Permission denied: This operation is restricted to the master.")
        return
      }
    }

    // Send to Master
    conn.Write([]byte(cmd + "\n"))

    // Execute locally
    executeLocalSQL(cmd)
  }
}


func main() {
  var err error
  db, err = sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/") 
  if err != nil {
    log.Fatal(err)
  }

  masterIP := "192.168.137.61" // Change as needed
  conn := connectToMaster(masterIP)
  defer conn.Close()

  go listenForMaster(conn)
  snapTerminal(conn)
}