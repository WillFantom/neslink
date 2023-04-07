package neslink

// Action represents an entity that has a name and some function (act) that can
// return an error.
type Action interface {
	name() string
	act() error
}
