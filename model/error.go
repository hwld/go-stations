package model

type ErrNotFound struct {
}

func (e *ErrNotFound) Error() string {
	return "データが見つかりませんでした"
}
