package loggingdrain

import "context"

type persistenceHandler interface {
	Save(context.Context, *TemplateMiner) error
	Load(context.Context) (*TemplateMiner, error)
}
