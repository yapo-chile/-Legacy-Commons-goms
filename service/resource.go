package service

// Resource is an example resource to be injected.
type Resource struct {
	Name  string
	Usage string
}

// SumLength is the sum of the lengths of Name and Usage of given resource.
func (r *Resource) SumLength() int {
	return len(r.Name) + len(r.Usage)
}
