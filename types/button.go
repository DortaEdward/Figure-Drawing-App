package types

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	Name      string
	Text      string
	Bounds    rl.Rectangle
	Color     rl.Color
	TextColor rl.Color
	Pressed   bool
	FontSize  int
	Screen    Screen
	OnClicked func()
}

func CreateButton(name, text string, x, y float32, width, height, fontSize int, color, textColor rl.Color, onClicked func(), screen Screen) *Button {
	return &Button{
		Name:      name,
		Text:      text,
		Bounds:    rl.NewRectangle(x, y, float32(width), float32(height)),
		Color:     color,
		TextColor: textColor,
		FontSize:  fontSize,
		OnClicked: onClicked,
		Screen:    screen,
	}
}

func (b *Button) CheckCollision(mousePosition rl.Vector2) bool {
	return rl.CheckCollisionPointRec(mousePosition, b.Bounds)
}

func (b *Button) Draw() {
	textWidth := rl.MeasureText(b.Text, int32(b.FontSize))
	rl.DrawRectangleRec(b.Bounds, b.Color)
	rl.DrawText(b.Text,
		(int32(b.Bounds.X+(b.Bounds.Width/2)) - (textWidth / 2)),
		int32(b.Bounds.Y+(b.Bounds.Height/4)),
		int32(b.FontSize),
		b.TextColor)
}
