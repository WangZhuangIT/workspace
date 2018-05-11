package models

import (
	db "gin-demo/database"
	"log"
)

type Person struct {
	Id        int    `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
}

func (p *Person) AddPerson() (id int64, err error) {
	rs, err := db.SqlDB.Exec("insert into person (first_name, last_name)values(?, ?)", p.FirstName, p.LastName)
	if err != nil {
		log.Fatalln(err)
	}
	id, err = rs.LastInsertId()
	return
}

func (p *Person) GetPerson() (persons []Person, err error) {
	persons = make([]Person, 0)
	rows, err := db.SqlDB.Query("select id, first_name, last_name from person")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var person Person
		rows.Scan(&person.Id, &person.FirstName, &person.LastName)
		persons = append(persons, person)
	}
	if err = rows.Err(); err != nil {
		log.Fatalln(err)
	}
	return
}

func (p *Person) GetPersonOne() (person Person, err error) {
	err = db.SqlDB.QueryRow("select id, first_name, last_name from person where id=?", p.Id).Scan(
		&person.Id, &person.FirstName, &person.LastName,
	)
	return
}

func (p *Person) UpdatePerson() (rowid int64, err error) {
	stmt, err := db.SqlDB.Prepare("update person set first_name = ?, last_name = ? where id = ?")
	defer stmt.Close()
	if err != nil {
		log.Fatalln(err)
	}

	re, err := stmt.Exec(p.FirstName, p.LastName, p.Id)
	if err != nil {
		log.Fatalln(err)
	}

	rowid, err = re.RowsAffected()
	return
}

func (p *Person) DelPerson() (err error, id int, rowid int64) {
	rs, err := db.SqlDB.Exec("delete from person where id = ?", p.Id)

	if err != nil {
		log.Fatalln(err)
	}
	rowid, err = rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	return
}
