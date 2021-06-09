package repo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prathik/tracker/repo"
	"github.com/prathik/tracker/service"
	"os"
	"time"
)

const repoFile = "repo.db"
var _ = Describe("BoltDbRepo", func() {
	Describe("Storing data", func() {
		Context("A new entry is added", func() {
			It("Stores the entry which can later be queried", func() {
				boltDbRepo := repo.NewBoltDbRepo(repoFile)
				boltDbRepo.Create(&service.Item{Time: time.Now(), Work: "test", Importance: 10, Joy: 10, Notes: "test"})
				data := boltDbRepo.QueryData(1 * time.Hour)
				Expect(data.DayDataCollection).Should(Not(BeEmpty()))
				boltDbRepo.Close()
				err := os.Remove(repoFile)
				if err != nil {
					panic(err)
				}
			})
		})
	})
})
