package cache

import (
	"context"

	"github.com/bentodev01/grello/internal/data"
	"github.com/go-redis/redis/v8"
)

type Caches struct {
	Board interface {
		Get(ctx context.Context, id string) data.BoardResult
		Set(ctx context.Context, b data.Board) error
	}
}

func NewCaches(redis *redis.Client) Caches {
	return Caches{
		Board: BoardCache{Cache: redis},
	}
}
