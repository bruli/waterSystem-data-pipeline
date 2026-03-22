package forecast

import "time"

const timeLocation = "Europe/Madrid"

type Slot struct {
	from, to time.Time
}

func (s Slot) From() time.Time {
	return s.from
}

func (s Slot) To() time.Time {
	return s.to
}

func Tomorrow() *Slot {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return &Slot{
		from: tomorrow,
		to:   tomorrow,
	}
}
