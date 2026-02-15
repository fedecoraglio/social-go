package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetByID(ctx context.Context, id int64) (*Post, error)
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, userID int64, feedQuery PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		GetByID(ctx context.Context, id int64) (*User, error)
		Create(ctx context.Context, user *User) error
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetByPostId(ctx context.Context, postId int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
