package repo

import (
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/prathik/tracker/domain"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"time"
)

const (
	bucket = "sessions"
	timeFormat = time.RFC3339
)

type boltDbRepo struct {
	db *bolt.DB
}

func (b *boltDbRepo) Query(duration time.Duration) ([]*domain.Session, error) {
	var sessions []*domain.Session

	err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucket))

		if bkt == nil {
			return errors.New("bucket does not exist")
		}

		cursor := bkt.Cursor()
		now := time.Now()
		nowStartOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		minDate := nowStartOfDay.Add(-duration).Format(timeFormat)
		min := []byte(minDate)
		max := []byte(now.Format(timeFormat))

		for k, v := cursor.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = cursor.Next() {

			var session domain.Session
			err := msgpack.Unmarshal(v, &session)
			if err != nil {
				return err
			}
			sessions = append(sessions, &session)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (b *boltDbRepo) Save(session *domain.Session) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		w, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatal(err)
		}
		timeSlot := GetSlotFromTime(session.Time)
		b, err := msgpack.Marshal(session)
		if err != nil {
			panic(err)
		}
		return w.Put(timeSlot, b)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (b *boltDbRepo) Pop() {
	err := b.db.Update(func(tx *bolt.Tx) error {
		w, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatal(err)
		}
		key, _ := w.Cursor().Last()
		return w.Delete(key)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (b *boltDbRepo) Close() {
	_ = b.db.Close()
}

func GetSlotFromTime(currentTime time.Time) []byte {
	timeSlot := []byte(currentTime.Format(timeFormat))
	return timeSlot
}

func NewBoltDbRepo(dbFilePath string) (*boltDbRepo, error) {
	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &boltDbRepo{db: db}, nil
}
