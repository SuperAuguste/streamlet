package streamlet

import (
	// "io"
	"bufio"
	"encoding/json"
	"os"

	"github.com/SuperAuguste/fsi/fsi"
)

type Streamlet struct {
	File        *os.File
	FSI         fsi.FSI
	Documents   []StreamletDocument
	DocumentIds []string
}

type StreamletDocument struct {
	Id   string
	Data map[string]interface{}
}

// Create / open a new Streamlet instance.
func New(file *os.File) Streamlet {

	return Streamlet{
		File:        file,
		FSI:         fsi.New(),
		Documents:   make([]StreamletDocument, 0),
		DocumentIds: make([]string, 0),
	}

}

func (db *Streamlet) deleteDocumentByIndex(idx int) {

	db.Documents[idx] = db.Documents[len(db.Documents)-1]
	db.Documents = db.Documents[:len(db.Documents)-1]

	db.DocumentIds[idx] = db.DocumentIds[len(db.DocumentIds)-1]
	db.DocumentIds = db.DocumentIds[:len(db.DocumentIds)-1]

}

func (db *Streamlet) findDocumentIndexById(id string) int {

	for index, _id := range db.DocumentIds {
		
		if _id == id {

			return index

		}

	}

	return -1

}

// Reads the contents of the database; it requires read permissions to function properly.
func (db *Streamlet) Init() {

	scanner := bufio.NewScanner(db.File)
	for scanner.Scan() {

		line := scanner.Text()
		if len(line) < 33 {

			continue

		}

		id := line[:33]
		data := line[34:]

		if data == "DELETE" {

			idx := db.findDocumentIndexById(id)
			db.deleteDocumentByIndex(idx)

		} else {

			var decoded map[string]interface{}
			json.Unmarshal([]byte(data), &decoded)

			if len(db.DocumentIds) != 0 && db.findDocumentIndexById(id) != -1 {

				idx := db.findDocumentIndexById(id)
				db.deleteDocumentByIndex(idx)

				db.Documents = append(db.Documents, StreamletDocument{
					Id:   id,
					Data: decoded,
				})
				db.DocumentIds = append(db.DocumentIds, id)

			} else {

				db.Documents = append(db.Documents, StreamletDocument{
					Id:   id,
					Data: decoded,
				})
				db.DocumentIds = append(db.DocumentIds, id)

			}

		}

	}

}

// Inserts a document into the database.
func (db *Streamlet) Insert(document interface{}) {

	j, _ := json.Marshal(document)
	db.File.WriteString(db.FSI.Generate() + "-" + string(j) + "\n")
	db.File.Sync()

}

// Inserts more than one line of documents into the database.
func (db *Streamlet) InsertBulk(documents []interface{}) {

	for _, i := range documents {

		j, _ := json.Marshal(i)
		db.File.WriteString(db.FSI.Generate() + "-" + string(j) + "\n")

	}

	db.File.Sync()

}

// Edits a document.
func (db *Streamlet) Edit(id string, document interface{}) {

	j, _ := json.Marshal(document)
	db.File.WriteString(id + "-" + string(j) + "\n")
	db.File.Sync()

}

func (db *Streamlet) Update(document StreamletDocument) {

	db.Edit(document.Id, document.Data)

}

// Deletes a document from the database.
func (db *Streamlet) Delete(id string) {

	db.File.WriteString(id + "-DELETE\n")
	db.deleteDocumentByIndex(db.findDocumentIndexById(id))
	db.File.Sync()

}

// Deletes more than one line of documents from the database.
func (db *Streamlet) DeleteBulk(ids []string) {

	for _, id := range ids {

		db.File.WriteString(id + "-DELETE\n")
		db.deleteDocumentByIndex(db.findDocumentIndexById(id))

	}

	db.File.Sync()

}

// Finds multiple documents that return true as a result of the callback function.
func (db *Streamlet) Find(callback func(document StreamletDocument) bool) []StreamletDocument {

	var docs []StreamletDocument
	for _, doc := range db.Documents {

		if callback(doc) {

			docs = append(docs, doc)

		}

	}

	return docs

}

// Finds one document that return true as a result of the callback function.
func (db *Streamlet) FindOne(callback func(document StreamletDocument) bool) StreamletDocument {

	for _, doc := range db.Documents {

		if callback(doc) {

			return doc

		}

	}

	return StreamletDocument{}

}

// Gets one document by id.
func (db *Streamlet) Get(id string) StreamletDocument {

	return db.FindOne(func(document StreamletDocument) bool {

		return document.Id == id

	})

}

// Deletes deleted documents and edits edited documents in the database file itself - WARNING: this will rewrite the database file entirely, use this sparingly or on small databases.
func (db *Streamlet) Clean() {

	os.Truncate(db.File.Name(), 0)
	for _, doc := range db.Documents {

		j, _ := json.Marshal(doc.Data)
		db.File.WriteString(doc.Id + "-" + string(j) + "\n")

	}

	db.File.Sync()

}
