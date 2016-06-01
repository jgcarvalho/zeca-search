package db

import (
	// "github.com/jgcarvalho/zeca/ca"
	"encoding/json"
	"log"
	"reflect"
	"strings" // "gopkg.in/mgo.v2"

	"github.com/boltdb/bolt"
	//"labix.org/v2/mgo/bson"
	"fmt"
)

// type Protein struct {
// 	Pdb_id string  "pdb_id"
// 	Chains []Chain "chains_data"
// }

type Config struct {
	Dir    string `toml:"db-dir"`
	Name   string `toml:"db-name"`
	Bucket string `toml:"bucket-name"`
	Init   string `toml:"init"`
	Target string `toml:"target"`
}

func (c *Protein) getField(field string) []string {
	r := reflect.ValueOf(c)
	s := reflect.Indirect(r).FieldByName(field)
	return s.Interface().([]string)
}

func GetProteins(db Config) (start, end []string, e error) {
	dbase, err := bolt.Open(db.Dir+db.Name, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	start = []string{"#"}
	end = []string{"#"}

	var result Protein
	dbase.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.Bucket))
		b.ForEach(func(k, v []byte) error {
			err := json.Unmarshal(v, &result)
			if err != nil {
				fmt.Println("DB error:", err)
			}
			start = append(start, result.getField(strings.Title(db.Init))...)
			start = append(start, "#")
			end = append(end, result.getField(strings.Title(db.Target))...)
			end = append(end, "#")
			return nil
		})
		return nil
	})
	if len(start) != len(end) {
		e = fmt.Errorf("Error: Number of CA start cells is different from end cells")
	}
	fmt.Println(start)
	fmt.Println(end)
	return start, end, e
}
