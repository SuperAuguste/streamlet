package streamlet

import (
	"os"
	"path/filepath"
	"testing"
)

type ExampleElementInsert struct {
	AString string
	ANumber int
}

func BenchmarkInsertBulk(b *testing.B) {

	exec, err := os.Executable()
	if err != nil {

		b.Fatal(err)

	}

	file, err := os.OpenFile(filepath.Join(filepath.Dir(exec), "database"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {

		b.Fatal(err)

	}
	defer file.Close()

	db := New(file)
	var arr []interface{}
	for i := 0; i < b.N; i++ {
		
		arr = append(arr, ExampleElementInsert{
			AString: "Hello World",
			ANumber: 111,
		})

	}
	db.InsertBulk(arr)

}
