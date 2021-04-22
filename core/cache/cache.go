package cache

// CacheKey ...
type CacheKey interface{}

// CacheValue ...
type CacheValue interface{}

// Cache ...
type Cache interface {
	GetCache(key CacheKey) (CacheValue, error)
	PutCache(key CacheKey, val CacheValue) error
	DeleteCache(key CacheKey) error
	CheckExist(key CacheKey) bool
	ExecuteForEachItem(f func(key CacheKey, value CacheValue, param ...interface{}), param ...interface{}) error
}

// LRU_MAX ...
const LRU_MAX = 20

// NewCache ...
func NewCache(maxArr ...int) Cache {
	cache := NewGeneralCache()
	return cache
}
