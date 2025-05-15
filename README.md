# Distributed Database System in Go

##  Overview

This project demonstrates a basic **Distributed Database System** using the **Go programming language**. It features a **master-slave architecture**, where one master node controls multiple slave (Snap) nodes. The system replicates SQL commands from the master to the slaves to maintain data consistency.

---

## Key Features

-  **Master-Slave Architecture**
-  **Data Replication** from master to slaves
-  **Restricted Operations** on slaves (no schema-altering operations allowed)
-  **TCP-based Communication**
-  **Manual SQL Query Execution via Terminal**
-  **MySQL** used as the underlying database

---

##  Master Node
- Listens for Snap (slave) connections via TCP.
- Accepts SQL commands from a terminal.
- Executes queries on its own MySQL instance.
- Broadcasts these queries to all connected Snap nodes.

### Snap (Slave) Node
- Connects to the master node via TCP.
- Listens for SQL commands from the master.
- Executes received queries on its local MySQL instance.
- Can also send permitted queries to the master for execution.

> Snap nodes are **restricted** from performing schema-altering operations like:
> - `CREATE DATABASE`
> - `DROP DATABASE`
> - `CREATE TABLE`
> - `DROP TABLE`
> - `ALTER TABLE`
> - `TRUNCATE TABLE`

---

##  Technologies Used

- **Go (Golang)** for concurrent networking and system logic.
- **MySQL** for relational database management.
- **TCP/IP** for communication between master and slave nodes.

---

##  Getting Started

###  Prerequisites

- Go 1.17+
- MySQL Server running locally on all nodes
- Properly configured user credentials for MySQL

###  Master Setup

1. Update your MySQL credentials in the master code:
   ```go
   db, err = sql.Open("mysql", "shery:password@tcp(127.0.0.1:3306)/")
