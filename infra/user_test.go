package infra

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tusmasoma/campfinder/domain/model"
	"github.com/tusmasoma/campfinder/domain/repository"
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

			ValidateErr(t, err, tt.want.err)

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

			ValidateErr(t, err, tt.want.err)

			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("GetUserByEmail() \n got = %v,\n want %v", user, tt.want.user)
			}
		})
	}
}

func TestUserRepo_Create(t *testing.T) {
	patterns := []struct {
		name    string
		setup   func(tx repository.SQLExecutor)
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
			txRepo := NewTransactionRepository(db)
			err := txRepo.Transaction(func(tx repository.SQLExecutor) error {
				if tt.setup != nil {
					tt.setup(tx)
				}

				repo := NewUserRepository(db)

				err := repo.Create(tt.in.ctx, tt.in.user, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.wantErr)

				if err == nil {
					var exists bool
					exists, err = repo.CheckIfUserExists(tt.in.ctx, tt.in.user.Email, repository.QueryOptions{Executor: tx})
					if !exists {
						t.Errorf("Create() is successful but there is not user in db:  err = %v", err)
					}
				}
				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	patterns := []struct {
		name  string
		setup func(tx repository.SQLExecutor)
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
			txRepo := NewTransactionRepository(db)
			err := txRepo.Transaction(func(tx repository.SQLExecutor) error {
				if tt.setup != nil {
					tt.setup(tx)
				}

				repo := NewUserRepository(db)

				err := repo.Update(tt.in.ctx, tt.in.user, repository.QueryOptions{Executor: tx})

				ValidateErr(t, err, tt.want.err)

				var user model.User
				user, err = repo.GetUserByID(tt.in.ctx, tt.in.user.ID.String(), repository.QueryOptions{Executor: tx})
				if err != nil {
					t.Errorf("After Update(), GetUserByID() error = %v, want no error", err)
				}
				if user.Name != tt.want.name {
					t.Errorf("After Update(), GetUserByID() got name = %v, want name = %v", user.Name, tt.want.name)
				}
				return nil
			})
			if err != nil {
				t.Error("Transaction failed")
			}
		})
	}
}
