package infra

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

type CacheItem struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Text      string    `db:"text"`
	Count     int       `db:"count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func TestCache(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	items := []CacheItem{
		{ID: uuid.NewString(), UserID: userID, Text: "bar", Count: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "bat", Text: "baz", Count: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: "qux", Text: "quux", Count: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	repo := NewRedisRepository(client)

	// set
	err := repo.Set(ctx, "item0", items[0])
	ValidateErr(t, err, nil)

	// get
	temp, err := repo.Get(ctx, "item0")
	ValidateErr(t, err, nil)

	getItem, _ := temp.(CacheItem)
	if !reflect.DeepEqual(getItem, items[0]) {
		t.Errorf("Get() \n got = %v,\n want %v", getItem, items[0])
	}

	// exists
	exists := repo.Exists(ctx, "item0")
	if !reflect.DeepEqual(exists, true) {
		t.Errorf("Exists() \n got = %v,\n want %v", exists, true)
	}

	// scan
	keys, err := repo.Scan(ctx, "item*")
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(keys, []string{"item0"}) {
		t.Errorf("Scan() \n got = %v,\n want %v", keys, []string{"item0"})
	}

	// delete
	err = repo.Delete(ctx, "item0")
	ValidateErr(t, err, nil)
	exists = repo.Exists(ctx, "item0")
	if !reflect.DeepEqual(exists, false) {
		t.Errorf("Delete() \n got = %v,\n want %v", exists, false)
	}
}
