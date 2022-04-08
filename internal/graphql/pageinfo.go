package graphql

type pageInfoResolver struct {
	hasPreviousPage bool
	hasNextPage     bool
	startCursor     *string
	endCursor       *string
}

func (r *pageInfoResolver) HasNextPage() (bool, error) {
	return r.hasNextPage, nil
}

func (r *pageInfoResolver) HasPreviousPage() (bool, error) {
	return r.hasPreviousPage, nil
}

func (r *pageInfoResolver) StartCursor() (*string, error) {
	return r.startCursor, nil
}

func (r *pageInfoResolver) EndCursor() (*string, error) {
	return r.endCursor, nil
}
