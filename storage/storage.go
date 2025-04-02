package storage

type Storage interface {
	Append(message string) error
	InsertBeforeLast(message string, isCalculation bool) error
}
