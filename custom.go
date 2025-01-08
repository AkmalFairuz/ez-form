package form

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

type custom struct {
	title              string
	elements           []Element
	elementsKeyToIndex map[string]int
	submitCallback     func(p *player.Player, response CustomResponse)
	closeCallback      func(p *player.Player)
}

func NewCustom() Custom {
	return &custom{
		elements:           make([]Element, 0),
		elementsKeyToIndex: make(map[string]int),
	}
}

func (c *custom) WithTitle(title string) Custom {
	c.title = title
	return c
}

func (c *custom) WithElement(key string, element Element) Custom {
	c.elements = append(c.elements, element)
	c.elementsKeyToIndex[key] = len(c.elements) - 1
	return c
}

func (c *custom) WithCallback(callback func(p *player.Player, response CustomResponse)) Custom {
	c.submitCallback = callback
	return c
}

func (c *custom) WithCloseCallback(callback func(p *player.Player)) Custom {
	c.closeCallback = callback
	return c
}

func (c *custom) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "custom_form",
		"title":   c.title,
		"content": c.elements,
	})
}

func (c *custom) SubmitJSON(data []byte, s form.Submitter, _ *world.Tx) error {
	p, ok := s.(*player.Player) // ONLY ACCEPT PLAYER
	if !ok {
		return nil
	}

	if data == nil {
		if c.closeCallback != nil {
			c.closeCallback(p)
		}
		return nil
	}

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	var inputData []any
	if err := dec.Decode(&inputData); err != nil {
		return fmt.Errorf("failed to decode JSON data to []any: %w", err)
	}
	if len(inputData) != len(c.elements) {
		return fmt.Errorf("expected %d elements, got %d", len(c.elements), len(inputData))
	}

	resp := customResponse{response: map[string]any{}}
	for key, index := range c.elementsKeyToIndex {
		if err := c.elements[index].submit(func(v any) {
			resp.response[key] = v
		}, inputData[index]); err != nil {
			return fmt.Errorf("failed to parse element %s: %w", key, err)
		}
	}

	c.submitCallback(p, resp)
	return nil
}

type Custom interface {
	form.Form

	WithTitle(title string) Custom
	WithElement(key string, element Element) Custom
	WithCallback(callback func(p *player.Player, response CustomResponse)) Custom
	WithCloseCallback(callback func(p *player.Player)) Custom
}
