package cache

type CacheService interface {
	// Insert new entities into sorted sets with the offset number as the score
	Create(key string, delta string, score int64) error
	//Fetch a range of offsets using a specified offset number as the starting offset
	Read(key string, offset int64) ([]string, error)
}