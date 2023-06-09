package cmd

// Builder ...
type Builder interface {
	Build(*Answer) (string, error)
}
