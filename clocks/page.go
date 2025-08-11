package main

import (
	"context"
	"fmt"
	"io"

	"github.com/jfyne/live"
	"github.com/jfyne/live/page"
)

const (
	validateTZ = "validate-tz"
	addTime    = "add-time"
)

// PageState the state we are tracking for our page.
type PageState struct {
	Title           string
	ValidationError string
	Clocks          []*page.Component
}

// newPageState create a new page state.
func newPageState(title string) *PageState {
	return &PageState{
		Title:  title,
		Clocks: []*page.Component{},
	}
}

// pageRegister register the pages events.
func pageRegister(c *page.Component) error {
	// Handler for the timezone entry validation.
	c.HandleEvent(validateTZ, func(_ context.Context, p live.Params) (any, error) {
		// Get the current page component state.
		state, _ := c.State.(*PageState)

		// Get the tz coming from the form.
		tz := p.String("tz")

		// Try to make a new ClockState, this will return an error if the
		// timezone is not real.
		if _, err := NewClockState(tz); err != nil {
			state.ValidationError = fmt.Sprintf("Timezone %s does not exist", tz)
			return state, nil
		}

		// If there was no error loading the clock state reset the
		// validation error.
		state.ValidationError = ""

		return state, nil
	})

	// Handler for adding a timezone.
	c.HandleEvent(addTime, func(_ context.Context, p live.Params) (any, error) {
		// Get the current page component state.
		state, _ := c.State.(*PageState)

		// Get the timezone sent from the form input.
		tz := p.String("tz")
		if tz == "" {
			return state, nil
		}

		// Use the page.Init function to create a new clock, register it and mount it.
		clockComponent, err := page.Init(context.Background(), func() (*page.Component, error) {
			// Each clock requires its own unique stable ID. Events for each clock can then find
			// their own component.
			return NewClockComponent(fmt.Sprintf("clock-%d", len(state.Clocks)+1), c.Handler, c.Socket, tz)
		})
		if err != nil {
			return state, err
		}

		// Update the page state with the new clock.
		state.Clocks = append(state.Clocks, clockComponent)

		// Return the state to have it persisted.
		return state, nil
	})

	return nil
}

// pageMount initialise the page component.
func pageMount(title string) page.MountHandler {
	return func(_ context.Context, c *page.Component) error {
		// Create a new page state.
		c.State = newPageState(title)
		return nil
	}
}

// pageRender render the page component.
func pageRender(w io.Writer, cmp *page.Component) error {
	pageState, ok := cmp.State.(*PageState)
	if !ok {
		return fmt.Errorf("no page data")
	}
	return PageView(pageState, cmp).Render(context.Background(), w)
}

// NewPage create a new page component.
func NewPageComponent(ID string, h *live.Handler, s *live.Socket, title string) (*page.Component, error) {
	return page.NewComponent(ID, h, s,
		page.WithRegister(pageRegister),
		page.WithMount(pageMount(title)),
		page.WithRender(pageRender),
	)
}
