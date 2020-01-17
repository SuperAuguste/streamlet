package streamlet

import (
	"fmt"
	"os"
	"testing"
)

var database Streamlet
var databaseFile *os.File

func TestMain(m *testing.M) {

	databaseFile, _ = os.OpenFile("test_database", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	database = New(databaseFile)

	os.Exit(m.Run())

}

func BenchmarkInsertBulk(b *testing.B) {

	var arr []map[string]interface{}
	for i := 0; i < b.N; i++ {

		arr = append(arr, map[string]interface{}{
			"a_string": "Hello World!",
			"a_number": 111,
		})

	}
	database.InsertBulk(arr)

}

func BenchmarkInit(b *testing.B) {

	database.Init()

	fmt.Println(database.Documents[database.Keys()[0]].Data)

}
