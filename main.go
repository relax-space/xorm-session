package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Fruit struct {
	Name string
	Code string
}

//Looking forward to session support concurrency, just like golang's context, it will not change when multi-threaded
func main() {
	db, err := xorm.NewEngine("mysql", "root:1234@tcp(127.0.0.1:3306)/fruit")
	if err != nil {
		panic(err)
	}
	db.Sync(
		new(Fruit),
	)
	//db.ShowSQL()
	for index := 0; index < 100; index++ {
		session := db.NewSession()
		fmt.Println(GetAll(session, 0, 2))
	}
	//output
	//2 [{1 1} {1 2}] <nil>
	//0 [] Table not found
}

func GetAll(session xorm.Interface, offset, limit int) (totalCount int64, fruits []Fruit, err error) {
	errc := make(chan error)

	go func() {
		v, err := session.Count(&Fruit{})
		if err != nil {
			errc <- err
			return
		}
		totalCount = v
		errc <- nil

	}()

	go func() {
		items := make([]Fruit, 0)
		if err := session.Cols("Name,Code").Limit(limit, offset).Find(&items); err != nil {
			errc <- err
			return
		}
		fruits = items
		errc <- nil
	}()

	if err := <-errc; err != nil {
		return 0, nil, err
	}
	if err := <-errc; err != nil {
		return 0, nil, err
	}
	return
}
