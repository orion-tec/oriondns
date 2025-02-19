package ai

type AI interface {
	Query(query string) (string, error)
}

type ai struct {
}

func New() AI {
	return &ai{}
}

func (a *ai) Query(query string) (string, error) {
	return "", nil
}
