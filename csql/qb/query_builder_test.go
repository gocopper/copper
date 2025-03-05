package qb_test

import (
	"testing"

	"github.com/gocopper/copper/csql/qb"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID      int    `db:"id,readonly"`
	Name    string `db:"name"`
	Email   string `db:"email"`
	Ignored string `db:"-"`
}

func TestPlaceholders(t *testing.T) {
	model := TestModel{ID: 1, Name: "Test", Email: "test@example.com"}
	result := qb.ValuePlaceholders(model)
	assert.Equal(t, "?, ?, ?", result)
}

func TestColumns(t *testing.T) {
	model := TestModel{}

	// Without alias
	result := qb.Columns(model)
	assert.Equal(t, "id, name, email", result)

	// With alias
	resultWithAlias := qb.Columns(model, "t")
	assert.Equal(t, "t.id, t.name, t.email", resultWithAlias)
}

func TestSetColumns(t *testing.T) {
	model := TestModel{}
	result := qb.SetColumns(model)
	assert.Equal(t, "name = ?, email = ?", result)
}

func TestValues(t *testing.T) {
	model := TestModel{ID: 1, Name: "Test", Email: "test@example.com"}
	result := qb.Values(model)

	assert.Len(t, result, 3)
	assert.Equal(t, 1, result[0])
	assert.Equal(t, "Test", result[1])
	assert.Equal(t, "test@example.com", result[2])
}

func TestSetValues(t *testing.T) {
	model := TestModel{ID: 1, Name: "Test", Email: "test@example.com"}
	result := qb.SetValues(model)

	assert.Len(t, result, 2)
	assert.Equal(t, "Test", result[0])
	assert.Equal(t, "test@example.com", result[1])
}

func TestWithEmbeddedStruct(t *testing.T) {
	type BaseModel struct {
		ID      int `db:"id,readonly"`
		Version int `db:"version"`
	}

	type UserModel struct {
		BaseModel
		Name string `db:"name"`
		Skip string `db:"-"`
	}

	model := UserModel{
		BaseModel: BaseModel{ID: 1, Version: 2},
		Name:      "Test",
		Skip:      "Ignored",
	}

	columns := qb.Columns(model)
	assert.Equal(t, "id, version, name", columns)

	setColumns := qb.SetColumns(model)
	assert.Equal(t, "version = ?, name = ?", setColumns)
}
