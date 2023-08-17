package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// docker build --tag sqlite_server_image . - create image
// docker run --publish 9997:9998 sqlite_server_image - запустить контейнер и связать локальный порт 9997 c портом 9998 в докере

func main() {
	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte("Hello, world!"))
		},
	)

	http.ListenAndServe(":9998", nil)

	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	request :=
		`
		DROP TABLE IF EXISTS test_table;
		CREATE TABLE test_table (
			id INTEGER PRIMARY KEY,
			name TEXT,
			x INTEGER
		);
		INSERT INTO test_table (name, x) VALUES ('egor', 11);
		INSERT INTO test_table (id, name, x) VALUES (2, 'egor', 12);
		INSERT INTO test_table (name, x) VALUES ('egor', 14);
		INSERT INTO test_table (id, name, x) VALUES (11, 'egor', 13);
		INSERT INTO test_table (name, x) VALUES ('egor', 14);
		`
	_, err = db.Exec(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.Query("SELECT * FROM test_table WHERE x > 12")
	defer rows.Close()
	fmt.Printf("%T \n", rows)

	fmt.Println("\nROWS INSIDE:")

	for rows.Next() {
		var id, x int
		var name string
		err = rows.Scan(&id, &name, &x)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(id, name, x)
	}
}
