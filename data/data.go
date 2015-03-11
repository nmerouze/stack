package data

type Params map[string]interface{}

// type Entity interface {
// 	Collection() string
// }

type Finder interface {
	Find(Params) error
}

type FinderFunc func(Params) error

func (f FinderFunc) Find(p Params) error {
	return f(p)
}
