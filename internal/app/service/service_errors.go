package service

type ErrValidationErr struct {
	Msg string
	Err error
}

func (e ErrValidationErr) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func (e ErrValidationErr) Unwrap() error {
	return e.Err
}

type ErrRepositoryErr struct {
	Msg string
	Err error
}

func (e ErrRepositoryErr) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func (e ErrRepositoryErr) Unwrap() error {
	return e.Err
}

type ErrCacheErr struct {
	Msg string
	Err error
}

func (e ErrCacheErr) Error() string {
	return e.Msg + ": " + e.Err.Error()
}

func (e ErrCacheErr) Unwrap() error {
	return e.Err
}
