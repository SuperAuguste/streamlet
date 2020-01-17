package streamlet

import (
	// "io"
	"bufio"
	"io"

	"github.com/SuperAuguste/fsi/fsi"
	jsoniter "github.com/json-iterator/go"
)

type Streamlet struct {
	ReadWriter io.ReadWriter
	Writer     *bufio.Writer
	FSI        fsi.FSI
	Documents  map[string]StreamletDocument
}

type StreamletDocument struct {
	Id   string
	Data map[string]interface{}
}

var jsond = jsoniter.ConfigFastest

// Create / open a new Streamlet instance.
func New(readWriter io.ReadWriter) Streamlet {

	return Streamlet{
		ReadWriter: readWriter,
		Writer:     bufio.NewWriter(readWriter),
		FSI:        fsi.New(),
		Documents:  make(map[string]StreamletDocument),
	}

}

func (db *Streamlet) readLine(line string) {

	id := line[:33]
	data := line[34:]

	if data == "DELETE" {

		delete(db.Documents, id)

	} else {

		var decoded map[string]interface{}
		// jsond.b().
		jsond.UnmarshalFromString(data, &decoded)

		db.addDocument(id, decoded)

	}

}

// Reads the contents of the database.
func (db *Streamlet) Init() {

	scanner := bufio.NewScanner(db.ReadWriter)
	for scanner.Scan() {

		db.readLine(scanner.Text())

	}

}

func (db *Streamlet) addDocument(id string, data map[string]interface{}) {

	db.Documents[id] = StreamletDocument{
		Id:   id,
		Data: data,
	}

}

func (db *Streamlet) writeString(data string) {

	db.Writer.WriteString(data)

}

func (db *Streamlet) save() {

	db.Writer.Flush()

}

// Inserts a document into the database.
func (db *Streamlet) Insert(document map[string]interface{}) {

	id := db.FSI.Generate()
	j, _ := jsond.Marshal(document)

	db.addDocument(id, document)
	db.writeString(id + "-" + string(j) + "\n")
	db.save()

}

// Inserts more than one line of documents into the database.
func (db *Streamlet) InsertBulk(documents []map[string]interface{}) {

	for _, i := range documents {

		id := db.FSI.Generate()
		j, _ := jsond.Marshal(i)

		db.addDocument(id, i)
		db.writeString(id + "-" + string(j) + "\n")

	}

	db.save()

}

// Edits a document.
func (db *Streamlet) Edit(id string, document map[string]interface{}) {

	j, _ := jsond.Marshal(document)
	db.addDocument(id, document)

	db.writeString(id + "-" + string(j) + "\n")
	db.save()

}

func (db *Streamlet) Update(document StreamletDocument) {

	db.Edit(document.Id, document.Data)

}

// Deletes a document from the database.
func (db *Streamlet) Delete(id string) {

	db.writeString(id + "-DELETE\n")
	delete(db.Documents, id)
	db.save()

}

// Deletes more than one line of documents from the database.
func (db *Streamlet) DeleteBulk(ids []string) {

	for _, id := range ids {

		db.writeString(id + "-DELETE\n")
		delete(db.Documents, id)

	}

	db.save()

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

	return db.Documents[id]

}

func (db *Streamlet) Keys() []string {

	keys := make([]string, 0)

	for key := range db.Documents {

		keys = append(keys, key)

	}

	return keys

}

// Deletes deleted documents and edits edited documents in the database file itself - WARNING: this will rewrite the database file entirely, use this sparingly or on small databases.
// func (db *Streamlet) Clean() {

// 	os.Truncate(db.File.Name(), 0)
// 	for _, doc := range db.Documents {

// 		j, _ := json.Marshal(doc.Data)
// 		db.File.WriteString(doc.Id + "-" + string(j) + "\n")

// 	}

// 	db.File.Sync()

// }
