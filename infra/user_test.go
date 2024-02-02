package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tusmasoma/campfinder/domain/model"
	"gotest.tools/assert"
)

type CheckIfUserExistsArg struct {
	ctx   context.Context
	email string
}

type GetUserByIDArg struct {
	ctx context.Context
	id  string
}

type GetUserByEmailArg struct {
	ctx   context.Context
	email string
}

type UserCreateArg struct {
	ctx  context.Context
	user *model.User
}

type UserUpdateArg struct {
	ctx  context.Context
	user model.User
}

func TestUserRepo_CheckIfUserExists(t *testing.T) {
	var err error
	patterns := []struct {
		name       string
		setup      func(db *sql.DB)
		in         CheckIfUserExistsArg
		wantExists bool
	}{
		{
			name: "Success",
			in: CheckIfUserExistsArg{
				ctx:   context.Background(),
				email: "test@gmail.com",
			},
			wantExists: true,
		},
		{
			name: "Fail",
			in: CheckIfUserExistsArg{
				ctx:   context.Background(),
				email: "test@icloud.com",
			},
			wantExists: false,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewUserRepository(db)

			var exists bool
			exists, err = repo.CheckIfUserExists(tt.in.ctx, tt.in.email)

			require.NoError(t, err)
			assert.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestUserRepo_GetUserByID(t *testing.T) {
	var err error

	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetUserByIDArg
		want  struct {
			user model.User
			err  error
		}
	}{
		{
			name: "Success",
			in: GetUserByIDArg{
				ctx: context.Background(),
				id:  "5fe0e237-6b49-11ee-b686-0242c0a87001",
			},
			want: struct {
				user model.User
				err  error
			}{
				user: model.User{
					ID:       uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: GetUserByIDArg{
				ctx: context.Background(),
				id:  "5fe0e237-6b49-11ee-b686-0242c0a00000",
			},
			want: struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  sql.ErrNoRows,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewUserRepository(db)

			var user model.User
			user, err = repo.GetUserByID(tt.in.ctx, tt.in.id)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByID() \n got = %v,\n want %v", user, tt.want.user)
			}
		})
	}
}

func TestUserRepo_GetUserByEmail(t *testing.T) {
	var err error

	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    GetUserByEmailArg
		want  struct {
			user model.User
			err  error
		}
	}{
		{
			name: "Success",
			in: GetUserByEmailArg{
				ctx:   context.Background(),
				email: "test@gmail.com",
			},
			want: struct {
				user model.User
				err  error
			}{
				user: model.User{
					ID:       uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					Name:     "test",
					Email:    "test@gmail.com",
					Password: "password123",
					IsAdmin:  false,
				},
				err: nil,
			},
		},
		{
			name: "Fail",
			in: GetUserByEmailArg{
				ctx:   context.Background(),
				email: "test@icloud.com",
			},
			want: struct {
				user model.User
				err  error
			}{
				user: model.User{},
				err:  sql.ErrNoRows,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewUserRepository(db)

			var user model.User
			user, err = repo.GetUserByEmail(tt.in.ctx, tt.in.email)

			if (err != nil) != (tt.want.err != nil) {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.want.err)
			}

			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByEmail() \n got = %v,\n want %v", user, tt.want.user)
			}
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	var err error

	patterns := []struct {
		name    string
		setup   func(db *sql.DB)
		in      UserCreateArg
		wantErr error
	}{
		{
			name: "Success",
			in: UserCreateArg{
				ctx: context.Background(),
				user: &model.User{
					Name:     "new_test",
					Email:    "new_test@gmail.com",
					Password: "password123",
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewUserRepository(db)

			err = repo.Create(tt.in.ctx, tt.in.user)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				var exists bool
				exists, err = repo.CheckIfUserExists(tt.in.ctx, tt.in.user.Email)
				if !exists {
					t.Errorf("Create() is successful but there is not user in db:  err = %v", err)
				}
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	var err error

	patterns := []struct {
		name  string
		setup func(db *sql.DB)
		in    UserUpdateArg
		want  struct {
			name string
			err  error
		}
	}{
		{
			name: "Success",
			in: UserUpdateArg{
				ctx: context.Background(),
				user: model.User{
					ID:       uuid.MustParse("5fe0e237-6b49-11ee-b686-0242c0a87001"),
					Name:     "updated_user",
					Email:    "test@gmail.com",
					Password: "password123",
				},
			},
			want: struct {
				name string
				err  error
			}{
				name: "updated_user",
				err:  nil,
			},
		},
	}

	for _, tt := range patterns {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(db)
			}

			repo := NewUserRepository(db)

			err = repo.Update(tt.in.ctx, tt.in.user)

			if (err != nil) != (tt.want.err != nil) {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.want.err)
			} else if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.want.err)
			}

			var user model.User
			user, err = repo.GetUserByID(tt.in.ctx, tt.in.user.ID.String())
			if err != nil {
				t.Errorf("After Update(), GetUserByID() error = %v, want no error", err)
			}
			if user.Name != tt.want.name {
				t.Errorf("After Update(), GetUserByID() got name = %v, want name = %v", user.Name, tt.want.name)
			}
		})
	}
}
