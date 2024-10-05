package sanitizer

import (
	"fmt"
)

func ErrInvalidPath(path string, message string) error {
	return fmt.Errorf("invalid path: %s - %s", path, message)
}
