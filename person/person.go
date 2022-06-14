package person

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"mux_test/db"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var redisDbConn db.RedisDbConn = db.RedisDbConn{
	RedisClient: db.GetRedisConn(),
}

func AddPerson(p Person) (int64, error) {
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
	if err == nil {
		lI, _ := res.LastInsertId()
		redisId := "id" + fmt.Sprint(lI)
		err = redisDbConn.SetEntry(redisId, p, 10)
		if err != nil {
			fmt.Printf("Error occured while adding person details in the redis: %v", err)
		}
	}
	lastId, _ := res.LastInsertId()
	return lastId, err
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
	var p Person
	redRec, err := redisDbConn.GetEntry("id" + fmt.Sprint(id))
	if err != nil {
		fmt.Printf("Error while fetching person (id = %d) records from redis.\nError:%v", id, err)
	} else if redRec != "" {
		fmt.Println("redRec" + redRec)
		json.Unmarshal([]byte(redRec), &p)
		return p, nil
	}
	sqlConn, err := db.NewSqlConn("mysql")

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
