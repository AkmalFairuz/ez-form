package form

import (
	"encoding/json"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/df-mc/dragonfly/server/world"
)

type modal struct {
	title   string
	content string
	button1 string
	button2 string

	callback      func(p *player.Player, button1 bool)
	closeCallback func(p *player.Player)
}

func NewModal(title string) Modal {
	return &modal{title: title}
}

func (m *modal) WithContent(content string) Modal {
	m.content = content
	return m
}

func (m *modal) WithButton1(text string) Modal {
	m.button1 = text
	return m
}

func (m *modal) WithButton2(text string) Modal {
	m.button2 = text
	return m
}

func (m *modal) WithCallback(callback func(p *player.Player, button1 bool)) Modal {
	m.callback = callback
	return m
}

func (m *modal) WithCloseCallback(callback func(p *player.Player)) Modal {
	m.closeCallback = callback
	return m
}

func (m *modal) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":    "modal",
		"title":   m.title,
		"content": m.content,
		"button1": m.button1,
		"button2": m.button2,
	})
}

func (m *modal) SubmitJSON(data []byte, s form.Submitter, _ *world.Tx) error {
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

	var value bool
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if m.callback != nil {
		m.callback(p, value)
	}
	return nil
}

type Modal interface {
	form.Form

	WithContent(content string) Modal
	WithButton1(text string) Modal
	WithButton2(text string) Modal
	WithCallback(callback func(p *player.Player, button1 bool)) Modal
	WithCloseCallback(callback func(p *player.Player)) Modal
}
