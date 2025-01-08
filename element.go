package form

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Element interface {
	json.Marshaler
	submit(setter func(v any), value any) error
}

type label struct {
	text string
}

func (l label) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": "label",
		"text": l.text,
	})
}

func (l label) submit(setter func(v any), value any) error {
	return nil
}

func NewLabel(text string) Element {
	return label{text: text}
}

type input struct {
	text         string
	placeholder  string
	defaultValue string
}

func (i input) WithPlaceholder(placeholder string) Input {
	i.placeholder = placeholder
	return i
}

func (i input) WithDefault(value string) Input {
	i.defaultValue = value
	return i
}

func (i input) submit(setter func(v any), value any) error {
	text, ok := value.(string)
	if !ok {
		return fmt.Errorf("value %v is not allowed in input element", value)
	}
	if !utf8.ValidString(text) {
		return fmt.Errorf("invalid utf8 string: %s", text)
	}
	setter(text)
	return nil
}

func NewInput(text string) Input {
	return input{text: text}
}

func (i input) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":        "input",
		"text":        i.text,
		"placeholder": i.placeholder,
		"default":     i.defaultValue,
	})
}

type Input interface {
	Element

	WithPlaceholder(placeholder string) Input
	WithDefault(value string) Input
}

type toggle struct {
	text    string
	enabled bool
}

func (t toggle) submit(setter func(v any), value any) error {
	enabled, ok := value.(bool)
	if !ok {
		return fmt.Errorf("value %v is not allowed in toggle element", value)
	}
	setter(enabled)
	return nil
}

func NewToggle(text string, enabled bool) Element {
	return toggle{text: text, enabled: enabled}
}

func (t toggle) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "toggle",
		"text":    t.text,
		"default": t.enabled,
	})
}

type dropdown struct {
	text         string
	options      []string
	defaultIndex int
}

func (d dropdown) WithOptions(options ...string) Dropdown {
	d.options = options
	return d
}

func (d dropdown) WithDefaultIndex(index int) Dropdown {
	d.defaultIndex = index
	return d
}

func (d dropdown) submit(setter func(v any), value any) error {
	index, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("value %v is not allowed in dropdown element", value)
	}
	i, err := index.Int64()
	if err != nil {
		return fmt.Errorf("failed to parse index as int: %w", err)
	}
	if i < 0 || i >= int64(len(d.options)) {
		return fmt.Errorf("index %d is out of bounds", i)
	}
	setter(int(i))
	return nil
}

func NewDropdown(text string) Dropdown {
	return dropdown{text: text}
}

func (d dropdown) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "dropdown",
		"text":    d.text,
		"options": d.options,
		"default": d.defaultIndex,
	})
}

type Dropdown interface {
	Element

	WithOptions(options ...string) Dropdown
	WithDefaultIndex(index int) Dropdown
}

type slider struct {
	text         string
	min          float64
	max          float64
	stepSize     float64
	defaultValue float64
}

func (s slider) WithStepSize(stepSize float64) Slider {
	s.stepSize = stepSize
	return s
}

func (s slider) WithDefault(value float64) Slider {
	s.defaultValue = value
	return s
}

func (s slider) submit(setter func(v any), value any) error {
	number, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("value %v is not allowed in slider element", value)
	}
	f, err := number.Float64()
	if err != nil {
		return fmt.Errorf("failed to parse value as float64: %w", err)
	}
	if f < s.min || f > s.max {
		return fmt.Errorf("value %f is out of bounds", f)
	}
	setter(f)
	return nil
}

func NewSlider(text string, min, max float64) Slider {
	return slider{text: text, min: min, max: max}
}

func (s slider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "slider",
		"text":    s.text,
		"min":     s.min,
		"max":     s.max,
		"step":    s.stepSize,
		"default": s.defaultValue,
	})
}

type Slider interface {
	Element

	WithStepSize(stepSize float64) Slider
	WithDefault(value float64) Slider
}

type stepSlider struct {
	dropdown
}

func NewStepSlider(text string) Dropdown {
	return stepSlider{dropdown{text: text}}
}

func (s stepSlider) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":    "step_slider",
		"text":    s.text,
		"options": s.options,
		"default": s.defaultIndex,
	})
}

func (s stepSlider) submit(setter func(v any), value any) error {
	index, ok := value.(json.Number)
	if !ok {
		return fmt.Errorf("value %v is not allowed in step slider element", value)
	}
	i, err := index.Int64()
	if err != nil {
		return fmt.Errorf("failed to parse index as int: %w", err)
	}
	if i < 0 || i >= int64(len(s.options)) {
		return fmt.Errorf("index %d is out of bounds", i)
	}
	setter(int(i))
	return nil
}

type button struct {
	text  string
	image string
}

func (btn button) WithImage(image string) Button {
	btn.image = image
	return btn
}

func (btn button) MarshalJSON() ([]byte, error) {
	m := map[string]any{"text": btn.text}
	if btn.image != "" {
		imageType := "path"
		if strings.HasPrefix(btn.image, "http://") || strings.HasPrefix(btn.image, "https://") {
			imageType = "url"
		}
		m["image"] = map[string]any{
			"type": imageType,
			"data": btn.image,
		}
	}
	return json.Marshal(m)
}

func NewButton(text string) Button {
	return button{
		text: text,
	}
}

type Button interface {
	WithImage(image string) Button
	MarshalJSON() ([]byte, error)
}
