package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/br0xen/termbox-util"
	"github.com/devict/magopie"
	"github.com/nsf/termbox-go"
)

const (
	tabSearch = iota
	tabResults
	tabError
)

type searchScreen struct {
	initialized bool
	tabIdx      int
	searchBox   *termboxUtil.InputField
	resMenu     *termboxUtil.Menu

	showClientAddrModal   bool
	showClientKeyModal    bool
	clientAddrM           *termboxUtil.InputModal
	clientKeyM            *termboxUtil.InputModal
	showClientNotSetModal bool
	clientNotSet          *termboxUtil.AlertModal

	showSearchingModal bool
	searchingModal     *termboxUtil.AlertModal

	message        string
	messageTimeout time.Duration
	messageTime    time.Time

	results      []magopie.Torrent
	resetResults bool

	defaultFg termbox.Attribute
	defaultBg termbox.Attribute
}

func (screen *searchScreen) handleKeyEvent(event termbox.Event) int {
	if event.Key == termbox.KeyCtrlH {
		// About
		return aboutScreenIndex
	} else if event.Key == termbox.KeyCtrlC {
		// Quit
		return exitScreenIndex
	} else if event.Key == termbox.KeyCtrlS {
		saveConfig()
		return searchScreenIndex
	} else if screen.showClientNotSetModal {
		screen.clientNotSet.HandleKeyPress(event)
		if screen.clientNotSet.IsDone() {
			screen.showClientNotSetModal = false
			screen.showClientAddrModal = true
		}
		return searchScreenIndex
	} else if screen.showClientAddrModal {
		screen.clientAddrM.HandleKeyPress(event)
		if screen.clientAddrM.IsDone() {
			val := screen.clientAddrM.GetValue()
			if !strings.HasPrefix(val, "http") {
				screen.clientAddrM.SetValue("http://" + val)
			}
			screen.showClientAddrModal = false
			screen.showClientKeyModal = true
		}
		return searchScreenIndex
	} else if screen.showClientKeyModal {
		screen.clientKeyM.HandleKeyPress(event)
		if screen.clientKeyM.IsDone() {
			screen.showClientKeyModal = false
			if !createClient(screen.clientAddrM.GetValue(), screen.clientKeyM.GetValue()) {
				screen.showClientNotSetModal = true
			}
		}
		return searchScreenIndex
	} else if event.Key == termbox.KeyTab {
		if screen.tabIdx == tabSearch {
			screen.tabIdx = tabResults
		} else {
			screen.tabIdx = tabSearch
		}
		return searchScreenIndex
	}

	switch screen.tabIdx {
	case tabSearch:
		if event.Key == termbox.KeyEnter {
			screen.performSearch()
		} else {
			screen.searchBox.HandleKeyPress(event)
		}
	case tabResults:
		screen.resMenu.HandleKeyPress(event)
		if screen.resMenu.IsDone() {
			i := screen.resMenu.GetSelectedIndex()
			if len(screen.results) > i {
				client.Download(&screen.results[i])
				screen.resMenu.SetDone(false)
			}
		}
	}
	return searchScreenIndex
}

func (screen *searchScreen) performLayout(stl style) {
	screen.defaultFg, screen.defaultBg = stl.defaultFg, stl.defaultBg
	if screen.messageTimeout > 0 && time.Since(screen.messageTime) > screen.messageTimeout {
		screen.clearMessage()
	}
	if !screen.initialized {
		width, height := termbox.Size()
		screen.searchBox = termboxUtil.CreateInputField(0, 1, width-2, 2, stl.defaultFg, stl.defaultBg)
		screen.searchBox.SetBordered(true)
		screen.resMenu = termboxUtil.CreateMenu("", []string{}, 0, 3, width-2, height-6, stl.defaultFg, stl.defaultBg)
		screen.resMenu.SetBordered(true)
		screen.resMenu.EnableVimMode()
		screen.setMessageWithTimeout("Ctrl+H for help, Ctrl+C to quit", -1)

		screen.clientAddrM = termboxUtil.CreateInputModal("Server Address", (width/2)-15, (height/2)-5, 30, 6, stl.defaultFg, stl.defaultBg)
		screen.clientAddrM.ShowHelp(false)
		screen.clientKeyM = termboxUtil.CreateInputModal("Server API Key", (width/2)-15, (height/2)-5, 30, 6, stl.defaultFg, stl.defaultBg)
		screen.clientKeyM.ShowHelp(false)
		screen.clientNotSet = termboxUtil.CreateAlertModal("Client is not set up!", (width/2)-15, (height/2)-2, 30, 4, stl.defaultFg, stl.defaultBg)
		screen.clientNotSet.SetText("Press Enter to Setup")
		screen.clientNotSet.ShowHelp(false)
		if !createClient(*httpAddr, *apiKey) {
			screen.showClientNotSetModal = true
			screen.clientAddrM.SetDone(false)
			screen.clientKeyM.SetDone(false)
		}

		screen.searchingModal = termboxUtil.CreateAlertModal("Searching...", (width/2)-15, (height/2)-2, 30, 4, stl.defaultFg, stl.defaultBg)
		screen.searchingModal.SetText("Please Wait...")
		screen.searchingModal.ShowHelp(false)
	}
	screen.initialized = true

	if screen.resetResults {
		var resStr []string
		for i := range screen.results {
			res := screen.results[i]
			resStr = append(resStr, "("+prettyFileSize(res.Size)+")\t↑:"+strconv.Itoa(res.Seeders)+"\t↓:"+strconv.Itoa(res.Leechers)+"\t- "+res.Title+" ("+sites[res.SiteID]+")")
		}
		screen.resMenu.SetOptionsFromStrings(resStr)
		screen.resetResults = false
		autoRefresh = false
	}
}

func (screen *searchScreen) drawScreen(stl style) {
	if screen.message == "" {
		screen.setMessageWithTimeout("Ctrl+H for help, Ctrl+C to quit", -1)
	}

	screen.searchBox.Draw()
	screen.resMenu.Draw()

	screen.drawHeader(stl)
	screen.drawFooter(stl)

	if screen.showClientNotSetModal {
		screen.clientNotSet.Draw()
	} else if screen.showClientAddrModal {
		screen.clientAddrM.Draw()
	} else if screen.showClientKeyModal {
		screen.clientKeyM.Draw()
	}

	if screen.showSearchingModal {
		screen.searchingModal.Draw()
	}
}

func (screen *searchScreen) drawHeader(stl style) {
	width, _ := termbox.Size()
	spaces := strings.Repeat(" ", (width / 2))
	termboxUtil.DrawStringAtPoint(fmt.Sprintf("%smagopie%s", spaces, spaces), 0, 0, stl.titleFg, stl.titleBg)
}
func (screen *searchScreen) drawFooter(stl style) {
	if screen.messageTimeout > 0 && time.Since(screen.messageTime) > screen.messageTimeout {
		screen.clearMessage()
	}
	_, height := termbox.Size()
	termboxUtil.DrawStringAtPoint(screen.message, 0, height-1, stl.defaultFg, stl.defaultBg)
}

func (screen *searchScreen) performSearch() {
	screen.showSearchingModal = true
	srchVal := screen.searchBox.GetValue()
	autoRefresh = true
	go func() {
		tc := client.Search(srchVal)
		// Clear the previous results
		screen.results = screen.results[:0]
		for i := 0; i < tc.Length(); i++ {
			screen.results = append(screen.results, *tc.Get(i))
		}
		screen.tabIdx = tabResults
		screen.resetResults = true
		screen.showSearchingModal = false
	}()
}

func (screen *searchScreen) setMessage(msg string) {
	screen.message = msg
	screen.messageTime = time.Now()
	screen.messageTimeout = time.Second * 2
}

// setMessageWithTimeout lets you specify the timeout for the message
// setting it to -1 means it won't timeout
func (screen *searchScreen) setMessageWithTimeout(msg string, timeout time.Duration) {
	screen.message = msg
	screen.messageTime = time.Now()
	screen.messageTimeout = timeout
}

func (screen *searchScreen) clearMessage() {
	screen.message = ""
	screen.messageTimeout = -1
}

func prettyFileSize(i int) string {
	iters := 0
	for i >= 1000 {
		i = i / 1000
		iters++
	}
	unit := "b"
	switch iters {
	case 1:
		unit = "K"
	case 2:
		unit = "M"
	case 3:
		unit = "G"
	case 4:
		unit = "T"
	}
	return strconv.Itoa(i) + unit
}
