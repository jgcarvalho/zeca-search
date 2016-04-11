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

func LoadProteinsFromBoltDB(dirname, dbname, bucket string) []Protein {
	db, err := bolt.Open(dirname+dbname, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}

	var result []Protein
	// var prot Protein
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
			// } else {
			// 	fmt.Println(prot)
			// 	result = append(result, prot)
			// }

			// fmt.Printf("key=%s, value=%s\n", k, v)
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

// TODO change string to []string
func GetProteins(db Config) (start, end []string, e error) {
	proteins := LoadProteinsFromBoltDB(db.Dir, db.Name, db.Bucket)
	start = []string{"#"}
	end = []string{"#"}
	for i := 0; i < len(proteins); i++ {
		// start = append(start, strings.Split(proteins[i].getField(strings.Title(db.Init)), "")...)
		start = append(start, proteins[i].getField(strings.Title(db.Init))...)
		start = append(start, "#")
		end = append(end, proteins[i].getField(strings.Title(db.Target))...)
		end = append(end, "#")
		// start += proteins[i].getField(strings.Title(db.Init)) + "#"
		// end += proteins[i].getField(strings.Title(db.Target)) + "#"
		// start += proteins[i].Chains[0].getField(strings.Title(db.Init)) + "#"
		// end += proteins[i].Chains[0].getField(strings.Title(db.Target)) + "#"
	}
	if len(start) != len(end) {
		e = fmt.Errorf("Error: Number of CA start cells is different from end cells")
	}
	fmt.Println(start)
	fmt.Println(end)
	return start, end, e
}
