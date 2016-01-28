package main

import "github.com/nsf/termbox-go"

type screen interface {
	handleKeyEvent(event termbox.Event) int
	performLayout(stl style)
	drawScreen(stl style)
}

const (
	searchScreenIndex = iota
	aboutScreenIndex
	exitScreenIndex
)

func defaultScreens() []screen {
	searchScreen := searchScreen{}
	aboutScreen := aboutScreen{}
	screens := []screen{
		&searchScreen,
		&aboutScreen,
	}

	return screens
}

func drawBackground(bg termbox.Attribute) {
	termbox.Clear(0, bg)
}

func layoutAndDrawScreen(screen screen, stl style) {
	screen.performLayout(stl)
	drawBackground(stl.defaultBg)
	screen.drawScreen(stl)
	termbox.Flush()
}

type style struct {
	defaultBg termbox.Attribute
	defaultFg termbox.Attribute
	titleFg   termbox.Attribute
	titleBg   termbox.Attribute
	cursorFg  termbox.Attribute
	cursorBg  termbox.Attribute
}

func defaultStyle() style {
	var stl style
	stl.defaultBg = termbox.ColorBlack
	stl.defaultFg = termbox.ColorWhite
	stl.titleFg = termbox.ColorBlack
	stl.titleBg = termbox.ColorGreen
	stl.cursorFg = termbox.ColorBlack
	stl.cursorBg = termbox.ColorGreen

	return stl
}

type cursor struct {
	x int
	y int
}
