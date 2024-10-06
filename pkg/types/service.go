package types

type Service interface {
	Start(errChan chan<- error)
	Stop() error
	GetName() string
}
