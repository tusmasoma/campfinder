package infra

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

type CacheItem struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Text      string    `db:"text" json:"text"`
	Count     int       `db:"count" json:"count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (item *CacheItem) Serialize() (string, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// DeserializeCacheItem は JSON 文字列から CacheItem インスタンスを生成します。
func DeserializeCacheItem(data string) (*CacheItem, error) {
	var item CacheItem
	err := json.Unmarshal([]byte(data), &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
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
	serialized, _ := items[0].Serialize()
	err := repo.Set(ctx, "item0", serialized)
	ValidateErr(t, err, nil)

	// get
	temp, err := repo.Get(ctx, "item0")
	ValidateErr(t, err, nil)

	getItem, _ := DeserializeCacheItem(temp)
	if d := cmp.Diff(getItem, &items[0], cmpopts.IgnoreFields(CacheItem{}, "CreatedAt", "UpdatedAt")); len(d) != 0 {
		t.Errorf("Get()differs: (-got +want)\n%s", d)
	}

	// exists
	exists := repo.Exists(ctx, "item0")
	if !reflect.DeepEqual(exists, true) {
		t.Errorf("Exists() \n got = %v,\n want = %v", exists, true)
	}

	// scan
	keys, err := repo.Scan(ctx, "item*")
	ValidateErr(t, err, nil)
	if !reflect.DeepEqual(keys, []string{"item0"}) {
		t.Errorf("Scan() \n got = %v,\n want = %v", keys, []string{"item0"})
	}

	// delete
	err = repo.Delete(ctx, "item0")
	ValidateErr(t, err, nil)
	exists = repo.Exists(ctx, "item0")
	if !reflect.DeepEqual(exists, false) {
		t.Errorf("Delete() \n got = %v,\n want = %v", exists, false)
	}
}
