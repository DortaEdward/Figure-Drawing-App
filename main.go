package main

import (
	"github.com/dortaedward/image_viewer/types"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	window := types.CreateNewWindow()
	rl.InitWindow(int32(window.Width), int32(window.Height), window.Title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	window.ProgramInit()
	window.Run()
}
