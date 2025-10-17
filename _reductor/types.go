package reductor

// пустая модель если надо вернуть по умолчанию не известно что
type ModelGeneric struct{}

type ModelList map[ModelType]interface{}

// формат сообщения
type Message struct {
	Sender string
	Page   ModelType
	Model  interface{}
}

type IConfig interface {
	SetInConfig(key string, value interface{}, save ...bool) error
	GetKeyString(name string) string
	GetByName(name string) interface{}
}
