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
	Describe("Inbox", func() {
		Context("Store item", func() {
			It("Stores the item in the file", func() {
				boltDbRepo, _ := repo.NewBoltDbRepo(repoFile)
				defer boltDbRepo.Close()
				item := domain.NewInboxItem(time.Now(), "test", boltDbRepo)
				err := item.Save()
				Expect(err).Should(BeNil())
				items, err := boltDbRepo.GetAllInbox()
				Expect(err).Should(BeNil())
				Expect(items).ShouldNot(BeEmpty())
				Expect(items[0].Text).Should(Equal("test"))
				err = os.Remove(repoFile)
				if err != nil {
					panic(err)
				}
			})
		})
	})
	Describe("Query data", func() {
		Context("File not present", func() {
			It("Gets empty data", func() {
				boltDbRepo, err := repo.NewBoltDbRepo("unknown-file")
				Expect(err).Should(BeNil())
				defer boltDbRepo.Close()
				_, err = boltDbRepo.Query(1 * time.Hour)
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
				boltDbRepo.Save(&domain.Session{Time: time.Now(), Challenge: "PERFECT"})
				sessions, _ := boltDbRepo.Query(1 * time.Hour)
				Expect(sessions).Should(HaveLen(1))
				item := sessions[0]
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
