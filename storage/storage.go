package storage

type Storage interface {
	Append(message string) error
}
