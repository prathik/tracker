package repo

import (
	"bytes"
	"github.com/boltdb/bolt"
	"github.com/prathik/tracker/domain"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"sort"
	"time"
)

const (
	bucket = "sessions"
	timeFormat = time.RFC3339
)

type boltDbRepo struct {
	db *bolt.DB
}

func (b *boltDbRepo) QueryData(duration time.Duration) *domain.Days {
	dayCount := map[string]*domain.Day{}
	_ = b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		now := time.Now()
		nowStartOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		minDate := nowStartOfDay.Add(-duration).Format(timeFormat)

		min := []byte(minDate)
		max := []byte(now.Format(timeFormat))

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			t, _ := time.Parse(timeFormat, string(k))
			dayData := dayCount[t.Format("06-01-02")]
			if dayData == nil {
				dayData = &domain.Day{Time: t, Count: 0}
			}
			var item domain.Session
			err := msgpack.Unmarshal(v, &item)
			if err != nil {
				panic(err)
			}
			dayData.Count = dayData.Count + 1
			dayData.Sessions = append(dayData.Sessions, &item)
			dayCount[t.Format("06-01-02")] = dayData
		}
		return nil
	})

	var dayData []*domain.Day
	for _, v := range dayCount {
		dayData = append(dayData, v)
	}

	sort.Slice(dayData, func(i, j int) bool {
		return dayData[i].Time.Before(dayData[j].Time)
	})

	return &domain.Days{Days: dayData}
}

func (b *boltDbRepo) Save(item *domain.Session) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		w, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatal(err)
		}
		timeSlot := GetSlotFromTime(item.Time)
		b, err := msgpack.Marshal(item)
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

func NewBoltDbRepo(dbFilePath string) *boltDbRepo {
	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &boltDbRepo{db: db}
}
