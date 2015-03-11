package data_test

import (
	"testing"

	"../data"
)

type tea struct {
	title string
}

type teaModel struct {
	db     *Db // TODO: Use an adapter
	entity *tea
}

// TODO: Better mapping
func (m teaModel) testFinder(p data.Params) error {
	v := m.db.Get(p["id"].(string))
	m.entity.title = v.(string)
	return nil
	// TODO: Test with the sql package
	// return c.db.QueryRow("SELECT * FROM teas WHERE id = $1", p["id"]).Scan(&e.(*tea).title)
}

type Db struct {
	data map[string]interface{}
}

func (d *Db) Get(k string) interface{} {
	return d.data[k]
}

func (d *Db) Set(k string, v interface{}) {
	d.data[k] = v
}

func NewDb() *Db {
	return &Db{map[string]interface{}{}}
}

func TestFinder(t *testing.T) {
	db := NewDb()
	db.Set("123", "sencha")

	e := new(tea)
	teaModel{db, e}.testFinder(data.Params{"id": "123"})

	if e.title != "sencha" {
		t.Fatalf("Tea's title should be sencha, got: %#v", e.title)
	}
}
