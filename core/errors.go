package core

import "fmt"

// NoMatch is returned by Endpoint if Input is not matched.
type NoMatch struct{}

func (err NoMatch) Error() string {
	return fmt.Sprintf("No Match")
}
