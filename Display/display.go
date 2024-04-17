package Display

import (
	"errors"
	"math"
	"sync"
	"time"

	rws "github.com/PlayerR9/MyGoLib/CustomData/Safe/RWSafe"
	ers "github.com/PlayerR9/MyGoLib/Utility/Errors"
	"github.com/gdamore/tcell"
)

type Display struct {
	frameRate time.Duration
	screen    tcell.Screen

	width, height int
	majorTable    *DrawTable

	closeSignal *rws.RWSafe[bool]
	eventCh     chan tcell.Event
	errCh       chan error

	wg sync.WaitGroup // WaitGroup for goroutines
}

func NewDisplay(frameRate float64) (*Display, error) {
	if frameRate <= 0 {
		return nil, ers.NewErrInvalidParameter("frameRate").
			Wrap(errors.New("value must be greater than 0"))
	}

	d := &Display{
		frameRate: time.Duration(math.Round(1000/frameRate)) * time.Millisecond,
	}

	var err error

	d.screen, err = tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := d.screen.Init(); err != nil {
		return nil, err
	}

	d.screen.Clear()
	d.screen.EnableMouse()

	d.width, d.height = d.screen.Size()

	d.majorTable = NewDrawTable(d.width, d.height, tcell.StyleDefault) // FIXME: Use a default style

	d.closeSignal = rws.NewRWSafe(false)
	d.eventCh = make(chan tcell.Event)
	d.errCh = make(chan error)

	return d, nil
}

func (d *Display) GetErrChannel() <-chan error {
	return d.errCh
}

func (d *Display) GetTable() WriteOnlyDTer {
	return d.majorTable
}

func (d *Display) Close() {
	d.closeSignal.Set(true)

	if d.errCh != nil {
		close(d.errCh)
		d.errCh = nil
	}

	if d.eventCh != nil {
		close(d.eventCh)
		d.eventCh = nil
	}

	d.wg.Wait()

	if d.screen != nil {
		d.screen.Fini()
		d.screen = nil
	}

	d.majorTable = nil
	d.closeSignal = nil
}

func (d *Display) Start() {
	// Handle the events in a separate goroutine
	go func() {
		for {
			d.eventCh <- d.screen.PollEvent()
		}
	}()

	d.wg.Add(2)

	// Handle the events
	go func() {
		defer d.wg.Done()
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				if d.closeSignal.Get() {
					return
				}
			case event := <-d.eventCh:
				switch event := event.(type) {
				case *tcell.EventResize:
					d.width, d.height = event.Size()

					err := d.majorTable.ResizeWidth(d.width)
					if err != nil {
						d.errCh <- err
					}

					err = d.majorTable.ResizeHeight(d.height)
					if err != nil {
						d.errCh <- err
					}
				case *tcell.EventKey:
					if event.Key() == tcell.KeyCtrlC {
						d.errCh <- NewErrESCPressed()
					}
				}
			}
		}
	}()

	go func() {
		defer d.wg.Done()

		for !d.closeSignal.Get() {
			d.screen.Clear()

			d.draw()

			d.screen.Show()

			time.Sleep(d.frameRate)
		}
	}()
}

func (d *Display) draw() {
	d.majorTable.mu.RLock()
	defer d.majorTable.mu.RUnlock()

	for y := 0; y < d.height; y++ {
		for x := 0; x < d.width; x++ {
			cell := d.majorTable.cells[y][x]

			d.screen.SetContent(x, y, cell.Char, nil, cell.Style)
		}
	}
}

func (d *Display) Clear() {
	d.screen.Clear()
}

func (d *Display) Wait() {
	d.wg.Wait()
}
