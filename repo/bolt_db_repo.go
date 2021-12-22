package repo

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/prathik/tracker/domain"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"sort"
	"time"
)

const (
	sessionsBucket = "sessions"
	timeFormat     = time.RFC3339
	inbox = "inbox"
)

type boltDbRepo struct {
	db *bolt.DB
}

func (b *boltDbRepo) Store(item *domain.InboxItem) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		w, err := tx.CreateBucketIfNotExists([]byte(inbox))
		if err != nil {
			log.Fatal(err)
		}

		timeSlot := make([]byte, 8)
		binary.LittleEndian.PutUint64(timeSlot, uint64(time.Now().UnixNano()))

		itemMarshalled, err := msgpack.Marshal(item)
		if err != nil {
			return err
		}
		return w.Put(timeSlot, itemMarshalled)
	})

	return err
}

func (b *boltDbRepo) Query(duration time.Duration) ([]*domain.Session, error) {
	var sessions []*domain.Session

	err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(sessionsBucket))

		if bkt == nil {
			return errors.New("sessions bucket does not exist")
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
		w, err := tx.CreateBucketIfNotExists([]byte(sessionsBucket))
		if err != nil {
			log.Fatal(err)
		}
		timeSlot := GetSlotFromTime(session.Time)
		data, err := msgpack.Marshal(session)
		if err != nil {
			panic(err)
		}
		return w.Put(timeSlot, data)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (b *boltDbRepo) Pop() {
	err := b.db.Update(func(tx *bolt.Tx) error {
		w, err := tx.CreateBucketIfNotExists([]byte(sessionsBucket))
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

func (b *boltDbRepo) GetAllInbox() ([]*domain.InboxItem, error) {
	var items []*domain.InboxItem
	err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(inbox))

		if bkt == nil {
			return errors.New("inbox bucket does not exist")
		}

		return bkt.ForEach(func(k, v []byte) error {
			var item domain.InboxItem
			err := msgpack.Unmarshal(v, &item)
			if err != nil {
				return err
			}
			items = append(items, &item)
			return nil
		})
	})

	sort.Slice(items, func(i, j int) bool {
		return items[i].CapturedTime.Before(items[j].CapturedTime)
	})

	return items, err
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
