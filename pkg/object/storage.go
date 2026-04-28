package object

type Storage interface {
	Set(key string, obj any) error
	Get(key string) (any, bool)
	Del(key string) error
}
