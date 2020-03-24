package aggregate

import "fmt"

// SnapshotError ...
type SnapshotError struct {
	Err       error
	StoreName string
}

func (err SnapshotError) Error() string {
	return fmt.Sprintf("%s snapshot: %s", err.StoreName, err.Err)
}

// Unwrap ...
func (err SnapshotError) Unwrap() error {
	return err.Err
}
