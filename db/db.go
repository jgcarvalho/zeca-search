package db

import (
	// "github.com/jgcarvalho/zeca/ca"
	"encoding/json"
	"log"
	"reflect"
	"strings" // "gopkg.in/mgo.v2"

	"github.com/boltdb/bolt"
	"github.com/jgcarvalho/zeca-search/rules"
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

func GetProteins(db Config) (start, end []rules.State, e error) {
	dbase, err := bolt.Open(db.Dir+db.Name, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	strStart := []string{"#"}
	strEnd := []string{"#"}

	var result Protein
	dbase.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db.Bucket))
		b.ForEach(func(k, v []byte) error {
			err := json.Unmarshal(v, &result)
			if err != nil {
				fmt.Println("DB error:", err)
			}
			strStart = append(strStart, result.getField(strings.Title(db.Init))...)
			strStart = append(strStart, "#")
			strEnd = append(strEnd, result.getField(strings.Title(db.Target))...)
			strEnd = append(strEnd, "#")
			return nil
		})
		return nil
	})
	if len(strStart) != len(strEnd) {
		e = fmt.Errorf("Error: Number of CA strStart cells is different from strEnd cells")
	}
	fmt.Println(strStart)
	fmt.Println(strEnd)

	start = make([]rules.State, len(strStart))
	end = make([]rules.State, len(strEnd))

	for i := range strStart {
		start[i] = rules.String2State(strStart[i])
		end[i] = rules.String2State(strEnd[i])
	}
	return start, end, e
}
