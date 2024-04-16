package redis

import (
	"context"
	"testing"
)

func TestUserSession(t *testing.T) {
	ctx := context.Background()
	userID := "f6db2530-cd9b-4ac1-8dc1-38c795e6eec2"
	jti := "d3b07384-d113-4ec6-a7d7-9a3bb5d3c8f5"

	repo := NewUserRepository(client)

	// set user session
	err := repo.SetUserSession(ctx, userID, jti)
	ValidateErr(t, err, nil)

	// get user session
	getJTI, err := repo.GetUserSession(ctx, userID)
	ValidateErr(t, err, nil)
	if getJTI != jti {
		t.Errorf("GetUserSession() \n got = %v,\n want = %v", getJTI, jti)
	}
}
