package form

import (
	"encoding/json"
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

type menu struct {
	title         string
	buttons       []Button
	content       string
	callback      func(p *player.Player, result int)
	closeCallback func(p *player.Player)
}

func NewMenu() Menu {
	return &menu{
		buttons: make([]Button, 0),
	}
}

func (m *menu) WithTitle(title string) Menu {
	m.title = title
	return m
}

func (m *menu) WithButton(text string, image ...string) Menu {
	btn := NewButton(text)
	if len(image) > 0 {
		btn.WithImage(image[0])
	}
	m.buttons = append(m.buttons, btn)
	return m
}

func (m *menu) WithContent(content string) Menu {
	m.content = content
	return m
}

func (m *menu) WithCallback(callback func(p *player.Player, result int)) Menu {
	m.callback = callback
	return m
}

func (m *menu) WithCloseCallback(callback func(p *player.Player)) Menu {
	m.closeCallback = callback
	return m
}

func (m *menu) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "form",
		"title":   m.title,
		"content": m.content,
	})
}

func (m *menu) SubmitJSON(data []byte, s form.Submitter, _ *world.Tx) error {
	p, ok := s.(*player.Player) // ONLY ACCEPT PLAYER
	if !ok {
		return nil
	}

	if data == nil {
		if m.closeCallback != nil {
			m.closeCallback(p)
		}
		return nil
	}

	var index int
	if err := json.Unmarshal(data, &index); err != nil {
		return fmt.Errorf("failed to parse button index as int: %w", err)
	}

	if index < 0 || index >= len(m.buttons) {
		return fmt.Errorf("button index out of bounds: %d", index)
	}

	if m.callback != nil {
		m.callback(p, index)
	}

	return nil
}

type Menu interface {
	form.Form

	WithTitle(title string) Menu
	WithButton(text string, image ...string) Menu
	WithContent(content string) Menu
	WithCallback(callback func(p *player.Player, result int)) Menu
	WithCloseCallback(callback func(p *player.Player)) Menu
}
