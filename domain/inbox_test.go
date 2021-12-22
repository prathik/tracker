package domain

import (
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestStoreInboxItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockInboxRepo(ctrl)
	item := NewInboxItem(time.Now(), "test", repo)
	repo.EXPECT().Store(item).Times(1)
	err := item.Save()
	if err != nil {
		t.Fail()
	}
}