// Package cui is a magopie client from the command line
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/devict/magopie"
	"github.com/nsf/termbox-go"
)

var client *magopie.Client
var sites map[string]string
var autoRefresh bool

var (
	httpAddr = flag.String("http", "", "Magopie Server Address")
	apiKey   = flag.String("key", "", "Shared API Key for the server")
)

func main() {
	flag.Parse()

	if *httpAddr == "" && *apiKey == "" {
		loadConfig()
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	stl := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)
	mainLoop(stl)
}

func createClient(addr, key string) bool {
	client = magopie.NewClient(addr, key)
	sc := client.ListSites()
	if sc.Length() == 0 {
		return false
	}
	sites = make(map[string]string)
	for i := 0; i < sc.Length(); i++ {
		sites[sc.Get(i).ID] = sc.Get(i).Name
	}
	return true
}

func mainLoop(stl style) {
	screens := defaultScreens()
	displayScreen := screens[searchScreenIndex]
	layoutAndDrawScreen(displayScreen, stl)
	eventChan := make(chan termbox.Event)
	go readUserInput(eventChan)
	go sendNoneEvent(eventChan)
	for {
		event := <-eventChan //termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeyCtrlZ {
				process, _ := os.FindProcess(os.Getpid())
				termbox.Close()
				process.Signal(syscall.SIGSTOP)
				termbox.Init()
			}
			newScreenIndex := displayScreen.handleKeyEvent(event)
			if newScreenIndex < len(screens) {
				displayScreen = screens[newScreenIndex]
				layoutAndDrawScreen(displayScreen, stl)
			} else {
				break
			}
		}
		if event.Type == termbox.EventResize || event.Type == termbox.EventNone {
			layoutAndDrawScreen(displayScreen, stl)
		}
	}
}

func readUserInput(e chan termbox.Event) {
	for {
		e <- termbox.PollEvent()
	}
}
func sendNoneEvent(e chan termbox.Event) {
	for {
		time.Sleep(time.Second)
		if autoRefresh {
			e <- termbox.Event{Type: termbox.EventNone}
		}
	}
}

func verifyOrCreateDirectory(path string) error {
	var tstDir *os.File
	var tstDirInfo os.FileInfo
	var err error
	if tstDir, err = os.Open(path); err != nil {
		if err = os.Mkdir(path, 0755); err != nil {
			return err
		}
		if tstDir, err = os.Open(path); err != nil {
			return err
		}
	}
	if tstDirInfo, err = tstDir.Stat(); err != nil {
		return err
	}
	if !tstDirInfo.IsDir() {
		return errors.New(path + " exists and is not a directory")
	}
	// We were able to open the path and it was a directory
	return nil
}

func loadConfig() {
	var cfgPath string
	cfgPath = os.Getenv("HOME")
	if cfgPath != "" {
		cfgPath = cfgPath + "/.config"
		err := verifyOrCreateDirectory(cfgPath)
		if err != nil {
			panic(err)
		}
		cfgPath = cfgPath + "/magopie-cui"
	}
	if cfgPath != "" {
		file, err := os.Open(cfgPath)
		if err != nil {
			// Couldn't load config even though one was specified
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tstString := scanner.Text()
			if strings.HasPrefix(tstString, "server=") {
				prts := strings.Split(tstString, "=")
				if len(prts) > 1 {
					flag.Set("http", prts[1])
				}
			} else if strings.HasPrefix(tstString, "key=") {
				prts := strings.Split(tstString, "=")
				if len(prts) > 1 {
					flag.Set("key", prts[1])
				}
			}
		}
	}
}

func saveConfig() {
	var cfgPath string
	var configLines []string
	configLines = append(configLines, "server="+client.ServerAddr)
	configLines = append(configLines, "key="+client.ServerKey)
	cfgPath = os.Getenv("HOME")
	if cfgPath != "" {
		cfgPath = cfgPath + "/.config"
		err := verifyOrCreateDirectory(cfgPath)
		if err != nil {
			panic(err)
		}
		cfgPath = cfgPath + "/magopie-cui"
	}
	if cfgPath != "" {
		file, err := os.Create(cfgPath)
		if err != nil {
			// Couldn't load config even though one was specified
			panic(err)
		}
		defer file.Close()

		w := bufio.NewWriter(file)
		for _, line := range configLines {
			fmt.Fprintln(w, line)
		}
		if err = w.Flush(); err != nil {
			panic(err)
		}
	}
}