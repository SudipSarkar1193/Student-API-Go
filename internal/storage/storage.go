package storage

type storage interface {
	CreateStudent(name string, email string) (int64, error)
}
