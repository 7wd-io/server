package pg

import "fmt"

var (
	errRoot = func(wrap error) error {
		return fmt.Errorf("pg suite: %w", wrap)
	}
)
