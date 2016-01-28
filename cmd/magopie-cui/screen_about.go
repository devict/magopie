package main

import (
	"fmt"

	"github.com/br0xen/termbox-util"
	"github.com/nsf/termbox-go"
)

type command struct {
	key         string
	description string
}

type aboutScreen struct {
	initialized  bool
	col1Commands []command
	col2Commands []command
	titleArt     *termboxUtil.ASCIIArt
}

func drawCommandsAtPoint(commands []command, x int, y int, stl style) {
	xPos, yPos := x, y
	for index, cmd := range commands {
		termboxUtil.DrawStringAtPoint(fmt.Sprintf("%6s", cmd.key), xPos, yPos, stl.defaultFg, stl.defaultBg)
		termboxUtil.DrawStringAtPoint(cmd.description, xPos+8, yPos, stl.defaultFg, stl.defaultBg)
		yPos++
		if index > 2 && index%2 == 1 {
			yPos++
		}
	}
}

func (screen *aboutScreen) handleKeyEvent(event termbox.Event) int {
	return searchScreenIndex
}

func (screen *aboutScreen) performLayout(stl style) {
	if !screen.initialized {
		width, _ := termbox.Size()
		var template []string
		template = append(template, "                                   _      ")
		template = append(template, "                                  (_)     ")
		template = append(template, " _ __ ___   __ _  __ _  ___  _ __  _  ___ ")
		template = append(template, "| '_ ` _ \\ / _` |/ _` |/ _ \\| '_ \\| |/ _ \\")
		template = append(template, "| | | | | | (_| | (_| | (_) | |_) | |  __/")
		template = append(template, "|_| |_| |_|\\__,_|\\__, |\\___/| .__/|_|\\___|")
		template = append(template, "                  __/ |     | |           ")
		template = append(template, "                 |___/      |_|           ")
		startX := (width - len(template[0])) / 2
		startY := 1

		screen.titleArt = termboxUtil.CreateASCIIArt(template, startX, startY, stl.defaultFg, stl.defaultBg)

		screen.col1Commands = append(screen.col1Commands, command{"j,↓", "down"})
		screen.col1Commands = append(screen.col1Commands, command{"k,↑", "up"})
		screen.col1Commands = append(screen.col1Commands, command{"", ""})
		screen.col1Commands = append(screen.col1Commands, command{"<enter>", "trigger download"})

		screen.col2Commands = append(screen.col2Commands, command{"ctrl+s", "save config"})
		screen.col2Commands = append(screen.col2Commands, command{"ctrl+q", "jump up"})
		screen.col2Commands = append(screen.col2Commands, command{"ctrl+h", "this screen"})
	}
	screen.initialized = true
}

func (screen *aboutScreen) drawScreen(stl style) {
	width, height := termbox.Size()
	defaultFg := stl.defaultFg
	defaultBg := stl.defaultBg

	screen.titleArt.Draw()

	xPos := screen.titleArt.GetX()
	yPos := screen.titleArt.GetY() + screen.titleArt.GetHeight() + 1

	tabTxt := "<tab>  switch between search and results"
	termboxUtil.DrawStringAtPoint(tabTxt, (width-len(tabTxt))/2, yPos, defaultFg, defaultBg)
	yPos++
	drawCommandsAtPoint(screen.col1Commands, xPos-5, yPos+1, stl)
	drawCommandsAtPoint(screen.col2Commands, xPos+30, yPos+1, stl)
	exitTxt := "Press any key to return to search"
	termboxUtil.DrawStringAtPoint(exitTxt, (width-len(exitTxt))/2, height-1, stl.titleFg, stl.titleBg)
}
