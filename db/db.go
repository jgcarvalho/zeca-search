package db

import (
	// "github.com/jgcarvalho/zeca/ca"
	"encoding/json"
	"log"
	"math/rand"
	"reflect"
	"strings" // "gopkg.in/mgo.v2"
	"time"

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

func LoadProteinsFromBoltDB(dirname, dbname, bucket string) []Protein {
	db, err := bolt.Open(dirname+dbname, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	var result []Protein
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		n := b.Stats().KeyN
		result = make([]Protein, n)
		i := 0
		b.ForEach(func(k, v []byte) error {
			err := json.Unmarshal(v, &result[i])
			if err != nil {
				fmt.Println("DB error:", err)
			}
			i++
			return nil
		})
		return nil
	})
	// fmt.Println(result)
	return result
}

func (c *Protein) getField(field string) []string {
	r := reflect.ValueOf(c)
	s := reflect.Indirect(r).FieldByName(field)
	return s.Interface().([]string)
}

func GetProteins(db Config) (start, end []string, e error) {
	rand.Seed(time.Now().UTC().UnixNano())
	proteins := LoadProteinsFromBoltDB(db.Dir, db.Name, db.Bucket)
	var get int
	n := 10
	start = []string{"#"}
	end = []string{"#"}
	// for i := 0; i < len(proteins); i++ {
	for i := 0; i < n; i++ {
		get = rand.Intn(len(proteins))
		start = append(start, proteins[get].getField(strings.Title(db.Init))...)
		start = append(start, "#")
		end = append(end, proteins[get].getField(strings.Title(db.Target))...)
		end = append(end, "#")
	}
	if len(start) != len(end) {
		e = fmt.Errorf("Error: Number of CA start cells is different from end cells")
	}
	fmt.Println(start)
	fmt.Println(end)
	return start, end, e
}
