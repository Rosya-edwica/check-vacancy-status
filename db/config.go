package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

)

type Database struct {
	Port       string
	Host       string
	User       string
	Password   string
	Name       string
	Connection *sql.DB
}

func (d *Database) Connect() {
	connection, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.User, d.Password, d.Host, d.Port, d.Name))
	checkErr(err)
	d.Connection = connection
}

func (d *Database) Close() {
	d.Connection.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}