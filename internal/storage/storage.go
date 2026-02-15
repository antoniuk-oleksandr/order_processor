package storage

import "sync"

// Storage provides thread-safe storage operations for user data.
// It maps user IDs to integer values and supports concurrent access.
type Storage interface {
	// Get retrieves the value associated with the given user ID.
	// Returns the value and true if found, or 0 and false if not found.
	Get(ID int) (int, bool)
	// Add increments the value for the given user ID by the specified amount.
	// If the ID doesn't exist, it will be created with the given value.
	Add(ID int, value int)
}

type storage struct {
	data map[int]int
	mu   sync.RWMutex
}

// NewStorage creates a new Storage instance with an empty data map.
// The returned storage is safe for concurrent use by multiple goroutines.
func NewStorage() Storage {
	return &storage{
		data: make(map[int]int),
	}
}

func (k *storage) Add(ID, value int) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.data[ID] += value
}

func (k *storage) Get(ID int) (int, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	val, ok := k.data[ID]
	return val, ok
}
