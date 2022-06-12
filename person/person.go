package person

import (
	"database/sql"
	"fmt"
	"mux_test/db"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func AddPerson(p Person) error {
	sqlCon, err := db.NewSqlConn("mysql")
	if err != nil {
		fmt.Printf("Error while creating new mysql connection: %v", err)
	}
	defer sqlCon.Close()
	var res sql.Result
	q := fmt.Sprintf("insert into person (name, age) values ('%s', %d)", p.Name, p.Age)
	res, err = sqlCon.Exec(q)
	if err != nil {
		fmt.Printf("Error while insert query execution: %v", err)
	}
	rowNos, err := res.RowsAffected()
	fmt.Printf("No of rows inserted: %d\n", rowNos)
	return err
}

func GetPersons() []Person {
	var ps []Person
	sqlCon, err := db.NewSqlConn("mysql")
	if err != nil {
		fmt.Printf("Error while creating new mysql connection: %v", err)
	}
	defer sqlCon.Close()
	q := "select id, name, age from person"
	rows, err := sqlCon.Query(q)
	if err != nil {
		fmt.Printf("Error while select query execution: %v", err)
	}
	for rows.Next() {
		var p Person
		rows.Scan(&p.Id, &p.Name, &p.Age)
		ps = append(ps, p)
	}
	return ps
}

func GetPerson(id int) (Person, error) {
	sqlConn, err := db.NewSqlConn("mysql")
	var p Person
	if err != nil {
		fmt.Printf("Error while creating new mysql connection")
		return p, err
	}
	defer sqlConn.Close()
	q := fmt.Sprintf("select id, name, age from person where id = %d", id)
	rows, err := sqlConn.Query(q)
	if err != nil {
		fmt.Printf("Error while executing select query: %v", err)
		return p, err
	}
	for rows.Next() {
		rows.Scan(&p.Id, &p.Name, &p.Age)
	}
	return p, err
}
