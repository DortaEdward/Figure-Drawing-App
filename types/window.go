package types

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Screen int

var timeRemaining time.Duration

const (
	ScreenStart Screen = iota
	ScreenMain
)

const (
	SCREEN_WIDTH  = 1080
	SCREEN_HEIGHT = 720
	IMAGES_DIR    = "./images/"
)

type Window struct {
	Width             int
	Height            int
	Title             string
	Images            []string
	Texture           rl.Texture2D
	CurrentPosition   int
	UsedImages        []int
	Buttons           []Button
	StartTime         time.Time
	EndTime           time.Time
	CurrentTime       time.Time
	TimeLeft          time.Duration
	CountdownDuration int
	Screen            Screen
}

func CreateNewWindow() *Window {
	return &Window{
		Width:  SCREEN_WIDTH,
		Height: SCREEN_HEIGHT,
		Title:  "References",
		Screen: ScreenStart,
	}
}

func (w *Window) ProgramInit() {
	w.GetImages()
	w.GetNewPosition()
	w.CreateTexture()
	w.InitButtons()

	// start timer
	w.CountdownDuration = 3 * 180
	w.ResetTimer()
}

func (w *Window) Run() {
	for !rl.WindowShouldClose() {
		w.HandleInput()
		rl.BeginDrawing()
		w.Draw()
		rl.EndDrawing()
	}
}

func (w *Window) InitButtons() {
	startButton := CreateButton(
		"StartButton",
		"Start",
		float32((w.Width/2)-(100)),
		float32((w.Height/2)-(25)),
		200,
		50,
		32,
		rl.Green,
		rl.White,
		func() {
			w.Screen = ScreenMain
		}, ScreenStart)

	nextButton := CreateButton(
		"NextButton",
		"Next",
		float32((w.Width-200)-50),
		float32(w.Height-100),
		200,
		50,
		32,
		rl.Green,
		// rl.Color{R: rl.Green.R, G: rl.Green.G, B: rl.Green.B, A: 128},
		rl.White,
		func() {
			w.NextImage()
		}, ScreenMain)

	w.Buttons = append(w.Buttons, *startButton, *nextButton)
}

func (w *Window) DrawStartScreen() {
	for _, button := range w.Buttons {
		if button.Screen == ScreenStart {
			button.Draw()
		}
	}
}

func (w *Window) DrawMainScreen() {
	currentTime := time.Now()
	timeRemaining = w.EndTime.Sub(currentTime)
	originalAspectRatio := float32(w.Texture.Width) / float32(w.Texture.Height)
	newHeight := float32(w.Height)
	newWidth := newHeight * originalAspectRatio

	x := (float32(w.Width) - newWidth) / 2
	y := (float32(w.Height) - newHeight) / 2

	rl.DrawTexturePro(
		w.Texture,
		rl.NewRectangle(0, 0, float32(w.Texture.Width), float32(w.Texture.Height)),
		rl.NewRectangle(x, y, newWidth, newHeight),
		rl.NewVector2(0, 0),
		0,
		rl.White,
	)

	for _, button := range w.Buttons {
		if button.Screen == ScreenMain {
			button.Draw()
		}
	}

	if timeRemaining <= 0 {
		w.NextImage()
	}
	timerText := fmt.Sprintf("Time Left: %02d:%02d", int(timeRemaining.Minutes()), int(timeRemaining.Seconds())%60)
	rl.DrawText(timerText, 10, 10, 30, rl.DarkGray)
}

func (w *Window) Draw() {
	rl.ClearBackground(rl.RayWhite)
	switch w.Screen {
	case ScreenStart:
		w.DrawStartScreen()
	case ScreenMain:
		w.DrawMainScreen()
	}
}

func (w *Window) GetImages() {
	dir, err := os.ReadDir(IMAGES_DIR)
	if err != nil {
		panic("ERROR: Error Reading dir")
	}

	for _, entry := range dir {
		w.Images = append(w.Images, entry.Name())
	}
}

func (w *Window) CreateTexture() {
	currentImage := rl.LoadImage(fmt.Sprintf("%s/%s", IMAGES_DIR, w.Images[w.CurrentPosition]))
	w.Texture = rl.LoadTextureFromImage(currentImage)
	rl.SetTextureFilter(w.Texture, rl.FilterBilinear)
	rl.UnloadImage(currentImage)
}

func (w *Window) HasImageBeenUsed(arrPos int) bool {
	used := false
	for _, pos := range w.UsedImages {
		if pos == arrPos {
			return true
		}
	}
	return used
}

func (w *Window) Check() int {
	if len(w.UsedImages) >= len(w.Images) {
		w.UsedImages = []int{}
	}
	used := false
	arrPos := rand.Intn(len(w.Images))
	for _, num := range w.UsedImages {
		if num == arrPos {
			used = true
			break
		}
	}
	if used {
		return w.Check()
	}
	w.UsedImages = append(w.UsedImages, arrPos)
	return arrPos
}

func (w *Window) GetNewPosition() {
	pos := w.Check()

	w.UsedImages = append(w.UsedImages, pos)
	w.CurrentPosition = pos
}

func (w *Window) HandleInput() {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mousePos := rl.GetMousePosition()

		switch w.Screen {
		case ScreenStart:
			for _, button := range w.Buttons {
				if button.Screen == ScreenStart {
					if button.CheckCollision(mousePos) {
						button.OnClicked()
					}
				}
			}
		case ScreenMain:
			for _, button := range w.Buttons {
				if button.Screen == ScreenMain {
					if button.CheckCollision(mousePos) {
						button.OnClicked()
					}
				}
			}
		}

	}
}

func (w *Window) ResetTimer() {
	w.StartTime = time.Now()
	w.EndTime = w.StartTime.Add(time.Second * time.Duration(w.CountdownDuration))
	w.TimeLeft = w.EndTime.Sub(w.CurrentTime)
}

func (w *Window) NextImage() {
	w.GetNewPosition()
	w.CreateTexture()
	w.ResetTimer()
}
