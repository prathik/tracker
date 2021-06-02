package repo

import (
	"bytes"
	"github.com/boltdb/bolt"
	"github.com/prathik/tracker/service"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"sort"
	"time"
)

const timeFormat = time.RFC3339

type boltDbRepo struct {
	db *bolt.DB
}

func (b *boltDbRepo) QueryData(duration time.Duration) *service.DayDataCollection {
	dayCount := map[string]*service.DayData{}
	_ = b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("work")).Cursor()
		now := time.Now()
		nowStartOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		minDate := nowStartOfDay.Add(-duration).Format(timeFormat)

		min := []byte(minDate)
		max := []byte(now.Format(timeFormat))

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			t, _ := time.Parse(timeFormat, string(k))
			dayData := dayCount[t.Format("06-01-02")]
			if dayData == nil {
				dayData = &service.DayData{Time: t, Count: 0}
			}
			var item service.Item
			err := msgpack.Unmarshal(v, &item)
			if err != nil {
				panic(err)
			}
			dayData.Count = dayData.Count + 1
			dayData.WorkItem = append(dayData.WorkItem, &item)
			dayCount[t.Format("06-01-02")] = dayData
		}
		return nil
	})

	var dayData []*service.DayData
	for _, v := range dayCount {
		dayData = append(dayData, v)
	}

	sort.Slice(dayData, func(i, j int) bool {
		return dayData[i].Time.Before(dayData[j].Time)
	})

	return &service.DayDataCollection{DayDataCollection: dayData}
}

func (b *boltDbRepo) Create(item *service.Item) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		w, err := tx.CreateBucketIfNotExists([]byte("work"))
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
		w, err := tx.CreateBucketIfNotExists([]byte("work"))
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
