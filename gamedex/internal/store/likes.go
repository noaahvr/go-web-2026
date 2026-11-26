package store

import "sync"

var (
	mu    sync.Mutex
	likes = map[int]int{}
)

func AddLike(gameID int) int {
	mu.Lock()
	defer mu.Unlock()
	likes[gameID]++
	return likes[gameID]
}

func Likes(gameID int) int {
	mu.Lock()
	defer mu.Unlock()
	return likes[gameID]
}
