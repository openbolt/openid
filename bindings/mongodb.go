package bindings

import (
	"github.com/openbolt/openid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// After Config, call Init(), on exit call Close()

// MongoDB is in MongoDB binding
type MongoDB struct {
	Host            string
	DBName          string
	CacheCollection string
	db              *mgo.Session
}

func (m *MongoDB) Init() (err error) {
	m.db, err = mgo.Dial(m.Host)
	if err != nil {
		return err
	}

	m.db.SetMode(mgo.Monotonic, true)
	return err
}

func (m *MongoDB) Close() {
	m.db.Close()
}

/*
 * Cacher
 */
func (m *MongoDB) Cache(val openid.Session) error {
	return m.db.DB(m.DBName).C(m.CacheCollection).Insert(val)
}

func (m *MongoDB) Retire(code string) {
	m.db.DB(m.DBName).C(m.CacheCollection).RemoveAll(bson.M{"Code": code})
}
