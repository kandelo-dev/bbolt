package bbolt_test

import (
	"path/filepath"
	"testing"

	bolt "go.etcd.io/bbolt"
)

func TestKandeloReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "kandelo.db")
	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("test"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte("key"), []byte("value"))
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = bolt.Open(path, 0o600, &bolt.Options{ReadOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.View(func(tx *bolt.Tx) error {
		got := string(tx.Bucket([]byte("test")).Get([]byte("key")))
		if got != "value" {
			t.Fatalf("reopened value: got %q, want value", got)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
