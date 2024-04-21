package Display

import (
	"errors"
	"math"
	"sync"
	"time"

	rws "github.com/PlayerR9/MyGoLib/CustomData/Safe/RWSafe"
	ers "github.com/PlayerR9/MyGoLibUnits/Errors"
	"github.com/gdamore/tcell"
)

// Display represents a display that can be used to draw to the screen.
type Display struct {
	// frameRate represents the frame rate of the display.
	frameRate time.Duration

	// screen represents the screen that the display will draw to.
	screen tcell.Screen

	// width and height represent the width and height of the screen, respectively.
	width, height int

	// majorTable represents the major table that the display will draw.
	majorTable *DrawTable

	// closeSignal represents a signal that will be used to close the display.
	closeSignal *rws.RWSafe[bool]

	// eventCh represents the channel that will be used to send events to the display.
	eventCh chan tcell.Event

	// errCh represents the channel that will be used to send errors to the display.
	errCh chan error

	// wg represents the WaitGroup that will be used to wait for goroutines to finish.
	wg sync.WaitGroup
}

// NewDisplay creates a new display with the given frame rate.
//
// If the frame rate is less than or equal to 0, an error of type *ers.ErrInvalidParameter
// will be returned.
//
// Parameters:
//
//   - frameRate: The frame rate of the display.
//
// Returns:
//
//   - *Display: A pointer to the new display.
//   - error: An error if the display could not be created.
func NewDisplay(frameRate float64) (*Display, error) {
	if frameRate <= 0 {
		return nil, ers.NewErrInvalidParameter(
			"frameRate",
			errors.New("value must be greater than 0"),
		)
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

	d.majorTable, err = NewDrawTable(d.width, d.height, tcell.StyleDefault) // FIXME: Use a default style
	if err != nil {
		return nil, err
	}

	d.closeSignal = rws.NewRWSafe(false)
	d.eventCh = make(chan tcell.Event)
	d.errCh = make(chan error)

	return d, nil
}

// GetErrChannel returns the error channel of the display.
//
// Returns:
//
//   - <-chan error: The error channel of the display.
func (d *Display) GetErrChannel() <-chan error {
	return d.errCh
}

// GetTable returns the table that the display will draw to.
//
// Returns:
//
//   - WriteOnlyDTer: The table that the display will draw to.
func (d *Display) GetTable() WriteOnlyDTer {
	return d.majorTable
}

// Close closes the display.
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

// Start starts the display.
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

// draw is a helper function that draws the major table to the screen.
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

// Clear clears the display.
func (d *Display) Clear() {
	d.screen.Clear()
}

// Wait waits for the display to finish.
func (d *Display) Wait() {
	d.wg.Wait()
}
