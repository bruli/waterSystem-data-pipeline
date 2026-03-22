package forecast

import (
	"context"
	"log/slog"
)

type Prediction struct {
	reader Reader
	repo   Repository
	log    *slog.Logger
}

func (p Prediction) Get(ctx context.Context, slot *Slot) error {
	weath, err := p.reader.Read(ctx, slot)
	if err != nil {
		return err
	}

	for _, w := range weath {
		if err = p.repo.Save(ctx, w); err != nil {
			p.log.ErrorContext(ctx, "error saving weather", slog.String("error", err.Error()))
		}
	}
	return nil
}

func NewPrediction(reader Reader, repo Repository, log *slog.Logger) *Prediction {
	return &Prediction{reader: reader, repo: repo, log: log}
}
