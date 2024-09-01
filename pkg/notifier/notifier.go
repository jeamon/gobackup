package notifier

import (
	"context"
	"log"

	"github.com/jeamon/gobackup/pkg/events"
	"github.com/jeamon/gobackup/pkg/fstypes"
	"github.com/jeamon/gorsn"
)

// Notifier wraps third-party gorsn.ScanNotifier
// in order to implement the Monitoring interface.
type Notifier struct {
	notifier gorsn.ScanNotifier
}

// New provides an instance of Notifier.
func New(root string, opts *gorsn.Options) (*Notifier, error) {
	w, err := gorsn.New(root, opts)
	if err != nil {
		return nil, err
	}
	return &Notifier{notifier: w}, nil
}

// Start implements Monitor `Start` behavior. Each event received is wrapped
// as ChangeEvent and propagated to jobs queue for further processing.
func (n *Notifier) Start(ctx context.Context, quit <-chan struct{}, jobs events.Queue) error {
	go func() {
		for {
			select {
			case event := <-n.notifier.Queue():
				jobs <- &events.Change{
					Path:  event.Path,
					Type:  fstypes.Type(event.Type),
					Ops:   events.Event(event.Name),
					Error: event.Error,
				}
			case <-quit:
				n.notifier.Stop()
				log.Println("stopped files monitoring")
				return
			}
		}
	}()

	return n.notifier.Start(ctx)
}

func (w *Notifier) Stop() error {
	return w.notifier.Stop()
}
