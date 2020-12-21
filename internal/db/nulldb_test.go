package db_test

import (
	"alertstore/internal/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNullDBObject(t *testing.T) {
	a := assert.New(t)

	n := db.NullDB{}
	a.Equal(n.String(), "null database driver")

	a.Nil(n.Save(nil))
	a.NoError(n.Ping())
	a.NoError(n.CheckModel())
}
