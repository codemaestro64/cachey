package store

type Store interface {
	Has(key string) bool
	Get(key string) any
	Put(key string, data any, duration int)
	Delete(key string)
	Flush()
	FlushExpired()
}
