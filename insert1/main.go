package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Fruit struct {
	Name string `json:"name" xorm:"unique 'name'"`
	Code string
}

//expect: Transaction takes effect
//actual: Transaction takes effect
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
	session.Begin()
	sussess := true
	for index := 0; index < 1; index++ {
		err := Insert(session, 0, 2)
		if err != nil {
			sussess = false
			fmt.Println(err)
			break
		}
	}
	if sussess {
		session.Commit()
	} else {
		session.Rollback()
	}
	fmt.Println("done!")
	time.Sleep(10 * time.Second)
}

func Insert(session xorm.Interface, offset, limit int) (err error) {
	errc := make(chan error)

	go func() {
		f := Fruit{
			Name: "xiao",
		}
		_, err = session.Insert(&f)
		if err != nil {
			errc <- err
			return
		}
		errc <- nil

	}()

	go func() {
		f := Fruit{
			Name: "xiao",
		}
		_, err = session.Insert(&f)
		if err != nil {
			errc <- err
			return
		}
		errc <- nil
	}()

	if err := <-errc; err != nil {
		return err
	}
	if err := <-errc; err != nil {
		return err
	}
	return
}
