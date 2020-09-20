package db

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// CreateDB creates Bolt db
func CreateDB(bucketName string) *bolt.DB {
	// Create db
	db, err := bolt.Open("url.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// WriteRecord saves record in the bucket 'bucketname'
func WriteRecord(db *bolt.DB, bucketName string, hashValue string, url string) {

	if err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		if err := bucket.Put([]byte(hashValue), []byte(url)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

// ReadRecord reads record with key 'hashValue' from bucket 'bucketname'
func ReadRecord(db *bolt.DB, bucketName string, hashValue string) string {
	var redirectURL string

	if err := db.View(func(tx *bolt.Tx) error {
		url := tx.Bucket([]byte(bucketName)).Get([]byte(hashValue))
		redirectURL = fmt.Sprintf("http://%s", url)
		return nil

	}); err != nil {
		log.Fatal(err)
	}

	return redirectURL
}
