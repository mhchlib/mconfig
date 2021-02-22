package cache

type CacheKey interface{}
type CacheValue interface{}

type Cache interface {
	GetCache(key CacheKey) (CacheValue, error)
	PutCache(key CacheKey, val CacheValue) error
	DeleteCache(key CacheKey) error
	CheckExist(key CacheKey) bool
	ExecuteForEachItem(f func(key CacheKey, value CacheValue, param ...interface{}), param ...interface{}) error
}

const LRU_MAX = 20

func NewCache(maxArr ...int) Cache {
	cache := NewGeneralCache()
	return cache
}
