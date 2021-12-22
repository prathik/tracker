package repo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prathik/tracker/domain"
	"github.com/prathik/tracker/repo"
	"os"
	"time"
)

const repoFile = "./repo.db"
var _ = Describe("BoltDbRepo", func() {
	Describe("Query data", func() {
		Context("File not present", func() {
			It("Gets empty data", func() {
				boltDbRepo, err := repo.NewBoltDbRepo("unknown-file")
				Expect(err).Should(BeNil())
				_, err = boltDbRepo.QueryData(1 * time.Hour)
				Expect(err).ShouldNot(BeNil())
				err = os.Remove("./unknown-file")
				if err != nil {
					panic(err)
				}
			})
		})
	})

	Describe("Storing data", func() {
		Context("A new entry is added", func() {
			It("Stores the entry which can later be queried", func() {
				boltDbRepo, _ := repo.NewBoltDbRepo(repoFile)
				boltDbRepo.Save(&domain.Session{Time: time.Now(), Challenge: "PERFECT", Notes: "test note"})
				data, _ := boltDbRepo.QueryData(1 * time.Hour)
				Expect(data.Days).Should(Not(BeEmpty()))
				Expect(data.Days[0].Count).Should(Equal(1))
				Expect(data.Days[0].Sessions).Should(HaveLen(1))
				item := data.Days[0].Sessions[0]
				Expect(item.Notes).Should(Equal("test note"))
				Expect(item.Challenge).Should(Equal("PERFECT"))
				boltDbRepo.Close()
				err := os.Remove(repoFile)
				if err != nil {
					panic(err)
				}
			})
		})
	})
})
