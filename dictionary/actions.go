package dictionary

import (
	"bytes"
	"encoding/gob"
	"sort"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

//Add ajoute un mot et la definition dans le dictionaire,
//renvoi une erreur
func (d *Dictionary) Add(word string, definition string) error {
	entry := Entry{
		Word:       strings.Title(word),
		Definition: definition,
		CreatedAt:  time.Now(),
	}

	//Encode le struct en bytes pour l'utiliser dans le update
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(entry)
	if err != nil {
		return err
	}
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(word), buffer.Bytes())
	})
}

//Get recupere tout de qui est en lien avec le mot passé en parametre
func (d *Dictionary) Get(word string) (Entry, error) {
	var entry Entry
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(word))
		if err != nil {
			return err
		}

		entry, err = getEntry(item)
		return err
	})

	return entry, err
}

//List recupere tout le contenue dans le dictionaire,
//[]string est un tableau trier alphabetique avec unique les mots,
//[string]Entry est une map de mot avec une definition
func (d *Dictionary) List() ([]string, map[string]Entry, error) {
	entries := make(map[string]Entry)
	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			entry, err := getEntry(item)
			if err != nil {
				return err
			}
			entries[entry.Word] = entry
		}
		return nil
	})

	return sortedKeys(entries), entries, err
}

//Remove supprime la ligne de la base de données avec un mot,
//word correspond au mot a supprimé
func (d *Dictionary) Remove(word string) error {

	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(word))
	})

}

func sortedKeys(entries map[string]Entry) []string {
	keys := make([]string, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func getEntry(item *badger.Item) (Entry, error) {
	var entry Entry
	var buffer bytes.Buffer
	item.Value(func(val []byte) error {
		_, err := buffer.Write(val)
		return err
	})

	dec := gob.NewDecoder(&buffer)
	err := dec.Decode(&entry)
	return entry, err

}
