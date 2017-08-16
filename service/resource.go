package service

/* Example Resource */
type Resource struct {
	Name  string
	Usage string
}

func (r *Resource) SumLength() int {
	return len(r.Name) + len(r.Usage)
}
