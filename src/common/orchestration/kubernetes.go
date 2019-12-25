package orchestration

type Namespace struct {
	name string
}
type Namespaces []Namespace

type WorkLoad struct{
	id string
	name string
	Type string
	createdAt string
}

