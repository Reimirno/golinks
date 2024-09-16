package mapper

import (
	"fmt"
)

func ErrMapConfigSetup(message string) error {
	return fmt.Errorf("invalid map configs setup: %s", message)
}

func ErrOperationNotSupported(operation string) error {
	return fmt.Errorf("operation not supported: %s", operation)
}

func ErrInvalidMapper(name string) error {
	return fmt.Errorf("invalid mapper: %s", name)
}
