package gprace

// RaceDriver struct
type RaceDriver struct {
	ID   int
	Name string
}

// NewDriver is the RaceDriver constructor
func NewDriver(id int, name string) (driver *RaceDriver) {

	d := new(RaceDriver)
	d.ID = id
	d.Name = name
	return d
}
