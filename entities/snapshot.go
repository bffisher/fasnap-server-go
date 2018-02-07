package entities

const (
	//OPSSave Saving opration
	OPSSave = iota
	//OPSDel Delete opration
	OPSDel
)

//Snapshot Save snapshot content
type Snapshot struct {
	Version   int64
	Operation int
	Date      string
	Content   string
}
