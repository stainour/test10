package domain

import "context"

type SequenceNumberGenerator interface {
	NextValue(context context.Context) (int64, error)
}
