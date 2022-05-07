package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

type botUser struct {
	Id              int64
	Sid             string
	Name            string
	Subscribe       bool
	Silently        bool
	RecentMessageId int
	RecentKeys      []string
	Parameters      map[string]string
	Messages        map[string]string
	Db              *DataBase
}

type Action struct {
	Db         *DataBase
	Result     chan *UpdateReturn
	Command    chan string
	Commands   *Commands
	Update     *Updates
	Bolt       *bolt.DB
	ProxyUsage bool
	ProxyUrl   string
	Autoload   []byte
	Cid        int
	Sleep      time.Duration
	Q          *Query
}

type DataBase struct {
	Db      *bolt.DB
	Err     error
	Backets *[]DataBaseBucket
	Lock    chan bool
	Timeout time.Duration
}

type DataBaseBucket struct {
	Parent  *bolt.DB
	Err     error
	Name    string
	Lock    chan bool
	Timeout time.Duration
	Inc     int
}

// New(*DataBase, int64) *botUser
// pair user and database
func (o *botUser) New(db *DataBase, key int64) *botUser {
	o.Db = db
	o.Id = key
	o.Sid = sprintf("%d", key)
	o.getUserInfo()
	return o
}

// getDbVal(string) string
// get user value by key from database
func (o *botUser) getDbVal(key string) string {
	return o.Db.FindCreate(o.Sid).Print(key)
}

// getDbVals(string, string) map[string]string
// get user values by key from database. Salt - additional value for bucket
func (o *botUser) getDbVals(salt, key string) map[string]string {
	if key == "" {
		return o.Db.FindCreate(o.Sid + salt).PrintAll()
	} else {
		return o.Db.FindCreate(o.Sid + salt).PrintAllPrefix(key)
	}
}

// getUserInfo()
// load user data from database
func (o *botUser) getUserInfo() {
	o.Subscribe = sBool(o.getDbVal("subscribe"))
	o.Silently = sBool(o.getDbVal("silently"))
	o.Parameters = make(map[string]string)
}

// pUM(string)
// get json value from database by key
func (o *botUser) pUM(key string) {
	json.Unmarshal([]byte(o.Db.FindCreate(o.Sid).Print(key)), &o.Parameters)
}

// pM(string)
// set json value to database by key
func (o *botUser) pM(key string) {
	data, _ := json.Marshal(o.Parameters)
	o.Db.FindCreate(o.Sid).Put(key, string(data))
}

// setParameter(string, string, interface{})
// upload parameter to database
func (o *botUser) setParameter(basket, key string, val interface{}) {
	o.pUM(basket)
	o.Parameters[key] = ssBool(val)
	o.pM(basket)
}

// getParameter(string, string)
// get parameter from database by key
func (o *botUser) getParameter(basket, key string) string {
	o.pUM(basket)
	return o.Parameters[key]
}

// turnParameter(string, string)
// values like "+", "-" convert reverse value in database
func (o *botUser) turnParameter(basket, key string) {
	o.pUM(basket)
	o.Parameters[key] = ssBool(!sBool(o.Parameters[key]))
	o.pM(basket)
}

// Open(string)
// new database with timeout checkout
func (o *DataBase) Open(name string) {
	o.Db, _ = bolt.Open(name+".db", 0600, &bolt.Options{Timeout: 1 * time.Second})
}

// Close()
// close database
func (o *DataBase) Close() {
	o.Db.Close()
}

// FindCreate(string) *DataBaseBucket
// new bucket in databae with name. If exist get it from database
func (o *DataBase) FindCreate(name string) *DataBaseBucket {
	ob := new(DataBaseBucket)
	ob.Name = name
	ob.Parent = o.Db
	o.Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ob.Name))
		if err != nil {
			return fmt.Errorf("-")
		}
		return nil
	})
	return ob
}

// Delete(string)
// delete backet from database
func (o *DataBase) Delete(name string) {
	o.Db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("-")
		}
		return nil
	})
}

// Add(string, string)
// add uniq value in backet
func (o *DataBaseBucket) Add(key, value string) {
	o.Parent.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(o.Name))
		c := tx.Bucket([]byte(o.Name)).Cursor()
		last, _ := c.Last()
		i := 1
		if last != nil {
			i, _ = strconv.Atoi(strings.Split(string(last), "_")[1])
			i++
		}
		increm := func(j *int) int { *j++; return *j }
		addnull := func(j int) string { ret := strings.TrimLeft(sprintf("%d", j+1000000000), "1"); return ret }
		for b.Get([]byte(key+"_"+addnull(i))) != nil {
			increm(&i)
		}
		err := b.Put([]byte(key+"_"+addnull(i)), []byte(value))
		return err
	})
}

// Put(string, string)
// add value in backet
func (o *DataBaseBucket) Put(key, value string) {
	o.Parent.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(o.Name))
		var err error
		if b != nil {
			err = b.Put([]byte(key), []byte(value))
		}
		return err
	})
}

// PrintAll() map[string]string
// get all values from backet. Debug
func (o *DataBaseBucket) PrintAll() map[string]string {
	m := make(map[string]string)
	i := 0
	o.Parent.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(o.Name))
		if b != nil {
			b.ForEach(func(k, v []byte) error {
				m[string(k)] = string(v)
				i++
				return nil
			})
		}
		return nil
	})
	return m
}

// PrintAllPrefix(string) map[string]string
// get values from backet by part key
func (o *DataBaseBucket) PrintAllPrefix(part string) map[string]string {
	m := make(map[string]string)
	o.Parent.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(o.Name)).Cursor()
		prefix := []byte(part)
		for k, v := b.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = b.Next() {
			m[string(k)] = string(v)
		}
		return nil
	})
	return m
}

// Print(string) string
// get value from backet by full key
func (o *DataBaseBucket) Print(key string) string {
	value := ""
	o.Parent.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(o.Name))
		if b != nil {
			v := b.Get([]byte(key))
			if v != nil {
				value = string(v)
			}
		}
		return nil
	})
	return value
}
