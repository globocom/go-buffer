package flusher

type Options struct {
	bypassOnStartError bool // default false
	bypassOnEachError  bool // default false
	bypassOnEndError   bool // default false
	OnStart            func([]interface{}) error
	OnEach             func(interface{}) error
	OnEnd              func() error
	OnStartError       func([]interface{}, error)
	OnEachError        func(interface{}, error)
	OnEndError         func([]interface{}, error)
}

type Flusher struct {
	Options *Options
}

func (flusher *Flusher) flushStart(items []interface{}) error {
	err := flusher.Options.OnStart(items)
	if err != nil {
		flusher.Options.OnStartError(items, err)
		if !flusher.Options.bypassOnStartError {
			return err
		}
	}
	return nil
}

func (flusher *Flusher) flushItems(items []interface{}) error {
	for _, item := range items {
		err := flusher.Options.OnEach(item)
		if err != nil {
			flusher.Options.OnEachError(item, err)
			if !flusher.Options.bypassOnEndError {
				return err
			}
		}
	}
	return nil
}

func (flusher *Flusher) flushEnd(items []interface{}) error {
	err := flusher.Options.OnEnd()
	if err != nil {
		flusher.Options.OnEndError(items, err)
		if !flusher.Options.bypassOnEndError {
			return err
		}
	}
	return nil
}

func (flusher *Flusher) Flush(items []interface{}) error {
	err := flusher.flushStart(items)
	if err != nil {
		return err
	}

	err = flusher.flushItems(items)
	if err != nil {
		return err
	}

	err = flusher.flushEnd(items)
	if err != nil {
		return err
	}

	return nil
}

func NewFlusher(options *Options) (*Flusher, error) {
	flusherInstance := &Flusher{
		Options: options,
	}
	return flusherInstance, nil
}
