package macropad

import (
	"errors"

	"github.com/pkg/term"
)

type EventHandler interface {
	OnEncoderChange(position int)
	OnEncoderPress()
	OnKeyPress(keynumber int)
}

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type Pixel struct {
	color      Color
	brightness uint8
}

func (p Pixel) Color(c Color) {
	panic("not implemented")
}

func (p Pixel) Brightness(b uint8) {
	panic("not implemented")
}

type Pixels []Pixel

func (p Pixels) Brightness(b uint8) {
}

type Macropad struct {
	vertical bool
	events   EventHandler
	Pixels   Pixels
	term     *term.Term
}

func New(device string, handler EventHandler) (Macropad, error) {
	term, err := term.Open(device, term.Speed(115200))
	if err != nil {
		return Macropad{}, err
	}

	return Macropad{
		term:   term,
		events: handler,
	}, nil
}

func (m *Macropad) Close() error {
	return errors.New("not implemented")
}

func (m *Macropad) DisplayTitle(title string) error {
	return errors.New("not implemented")
}

type TextFormat struct {
	Title      string
	Lines      []string
	TitleScale int
	TextScale  int
	Font       string
}

func (m *Macropad) DisplayText(tf TextFormat) error {
	return errors.New("not implemented")
}

func (m *Macropad) SetPixel() error {
	return errors.New("not implemented")
}
