package dictionary

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type Dictionary struct {
	db *badger.DB
}

type Entry struct {
	Word       string
	Definition string
	CreatedAt  time.Time
}

//String formate le texte
func (e Entry) String() string {
	created := e.CreatedAt.Format(time.Stamp)
	return fmt.Sprintf("%-10v\t%-50v%-6v", e.Word, e.Definition, created)
}

//New crée un Dictionary avec le dossier passé en parametre
func New(dir string) (*Dictionary, error) {
	opts := badger.DefaultOptions(dir)
	opts.Logger = nil
	opts.Dir = dir
	opts.ValueDir = dir

	db, err := badger.Open(opts)

	if err != nil {
		return nil, err
	}

	dict := &Dictionary{
		db: db,
	}
	return dict, nil
}

//Close ferme la base de données
func (d *Dictionary) Close() {
	d.db.Close()
}
