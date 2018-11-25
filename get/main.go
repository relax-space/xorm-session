package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Fruit struct {
	Name string `json:"name" xorm:"unique 'name'"`
	Code string
}

//expect:query result and total count
//actual: err:Table not found
func main() {
	db, err := xorm.NewEngine("mysql", "root:1234@tcp(127.0.0.1:3306)/fruit")
	if err != nil {
		panic(err)
	}
	db.Sync(
		new(Fruit),
	)
	//db.ShowSQL()
	session := db.NewSession()
	for index := 0; index < 100; index++ {
		totalCount, fruits, err := GetAll(session, 0, 2)
		if err != nil {
			fmt.Println(totalCount, fruits, err)
		}
	}
	fmt.Println("done")
}

func GetAll(session *xorm.Session, offset, limit int) (totalCount int64, fruits []Fruit, err error) {
	errc := make(chan error)

	go func(q xorm.Session) {
		v, err := q.Count(&Fruit{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v
		errc <- nil

	}(*session)

	go func(q xorm.Session) {
		items := make([]Fruit, 0)
		if err := q.Cols("Name,Code").Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		fruits = items
		errc <- nil
	}(*session)

	if err := <-errc; err != nil {
		return 0, nil, err
	}
	if err := <-errc; err != nil {
		return 0, nil, err
	}
	return
}
