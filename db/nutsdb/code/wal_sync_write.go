package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/xujiajun/nutsdb"
)

const (
	_batchCount = 100_000
)

var value = []byte(`HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890HelloWorld1234567890`)

func main() {
	opt := nutsdb.DefaultOptions
	opt.Dir = "/tmp/nutsdb" // 这边数据库会自动创建这个目录文件
	opt.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := batchTestWrite(db); err != nil {
		panic(err)
	}

	if err := batchTestRead(db); err != nil {
		panic(err)
	}
}

func batchTestRead(db *nutsdb.DB) error {

	return db.View(func(tx *nutsdb.Tx) error {

		for i := 0; i < _batchCount; i++ {

			var (
				key = []byte(fmt.Sprintf("key-%d", i))
			)

			entry, err := tx.Get("bucket-test", key)
			if err != nil {
				return err
			}

			// log.Printf("%s - %s", entry.Key, entry.Value)

			if !bytes.Equal(value, entry.Value) {
				panic("value data is wrong!!!")
			}

		}

		return nil
	})
}

func batchTestWrite(db *nutsdb.DB) error {
	defer func(begin time.Time) {

		log.Printf("write time_used: %d", time.Since(begin).Milliseconds())

	}(time.Now())

	return db.Update(func(tx *nutsdb.Tx) error {

		for i := 0; i < _batchCount; i++ {
			var (
				key = []byte(fmt.Sprintf("key-%d", i))
				// value = []byte(fmt.Sprintf("HelloWorld1234567890-%d", i))
			)

			if err := tx.Put("bucket-test", key, value, 0); err != nil {
				return err
			}
		}

		return nil
	})
}

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}
