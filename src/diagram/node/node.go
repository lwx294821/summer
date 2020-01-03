package node

type Node struct {
	ID string
	Category string
	Text string
	Status string
}
type Annotations struct {
	NodeId string
	metadata []MetaData
	metric []Metric
}
type Metric struct {
	name string
	value string
}
type MetaData struct {
	name string
	value string
}

