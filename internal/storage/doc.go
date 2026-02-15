// Package storage provides thread-safe in-memory storage for user data.
//
// The storage package offers a simple key-value store that maps user IDs
// to integer values. It is designed for concurrent access and uses read-write
// locks to ensure data consistency across multiple goroutines.
//
// Example usage:
//
//	store := storage.NewStorage()
//	store.Add(123, 100)
//	value, exists := store.Get(123)
//	if exists {
//		fmt.Printf("User 123 has value: %d\n", value)
//	}
package storage
