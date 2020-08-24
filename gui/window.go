package gui

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"log"
	"os"
	"tcpProxy/proxy"
)

type RowInput struct {
	SourcePort *widget.Entry
	DestIp     *widget.Entry
	DestPort   *widget.Entry
}
type RowInputData struct {
	SourcePort string
	DestIp     string
	DestPort   string
}

var RowInputMap = make(map[string]*RowInput)
var DataMap = make(map[string]*RowInputData, 0)
var isAdd = false

var addBtn = widget.NewButton(" + ", func() {})
var startBtn = widget.NewButton(" start ", func() {})

func RunWindow() {
	loadFromFile()

	a := app.New()
	w := a.NewWindow("TcpProxy")

	inputs := widget.NewVBox()

	addBtn.OnTapped = func() {
		addCoupleInputs(inputs)
	}
	startBtn.OnTapped = func() {
		startTcpProxy()
	}

	funcBtn := widget.NewHBox(
		layout.NewSpacer(),
		addBtn,
		startBtn,
		layout.NewSpacer(),
	)

	container := &widget.ScrollContainer{
		Content: widget.NewVBox(
			funcBtn,
			widget.NewGroup("proxys", ),
			inputs,
		),
		Direction: widget.ScrollVerticalOnly,
	}
	container.SetMinSize(fyne.NewSize(600, 450))

	loadInputRow(inputs)

	w.SetContent(container)
	w.ShowAndRun()
}

func loadFromFile() {
	b, err := ioutil.ReadFile("config.json") // just pass the file name
	if err != nil {
		fmt.Print(err)
		return
	}
	json.Unmarshal(b, &DataMap)
}

func saveToFile() {
	f, err := os.OpenFile("config.json", os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	os.Truncate("config.json",0)
	b, _ := json.Marshal(DataMap)
	_, err = f.Write(b)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func loadInputRow(box *widget.Box) {
	for k, _ := range DataMap {
		loadCoupleInputs(box, k)
	}
}

func loadCoupleInputs(box *widget.Box, key string) {
	sourcePort_Entry := widget.NewEntry()
	destIp_Entry := widget.NewEntry()
	destPort_Entry := widget.NewEntry()

	sourcePort_Entry.Text = DataMap[key].SourcePort
	destIp_Entry.Text = DataMap[key].DestIp
	destPort_Entry.Text = DataMap[key].DestPort

	delBtn := widget.NewButton("Remove", func() {
		deleteInputRow(box, key)
		saveToFile()
		refreshWindow(box)
	})
	RowInputMap[key] = &RowInput{
		SourcePort: sourcePort_Entry,
		DestIp:     destIp_Entry,
		DestPort:   destPort_Entry,
	}
	obs := []fyne.CanvasObject{
		layout.NewSpacer(),
		//widget.NewLabel(key),
		layout.NewSpacer(),
		sourcePort_Entry,
		widget.NewLabel("==>"),
		destIp_Entry,
		widget.NewLabel(":"),
		destPort_Entry,
		layout.NewSpacer(),
		delBtn,
		layout.NewSpacer(),
	}
	inputRow := widget.NewHBox(
		obs...
	)
	box.Children = append(box.Children, inputRow)
	box.Refresh()
}

func addCoupleInputs(box *widget.Box) {
	if isAdd {
		return
	}
	sourcePort_Entry := widget.NewEntry()
	destIp_Entry := widget.NewEntry()
	destPort_Entry := widget.NewEntry()
	okBtn := widget.NewButton("Confirm", func() {
		key := sourcePort_Entry.Text
		if DataMap[key] != nil {
			return
		}
		DataMap[key] = &RowInputData{
				SourcePort: key,
				DestIp:     destIp_Entry.Text,
				DestPort:   destPort_Entry.Text,
			}
		saveToFile()
		refreshWindow(box)

	})
	delBtn := widget.NewButton("Cancel", func() {
		refreshWindow(box)
	})
	obs := []fyne.CanvasObject{
		layout.NewSpacer(),
		layout.NewSpacer(),
		sourcePort_Entry,
		widget.NewLabel("==>"),
		destIp_Entry,
		widget.NewLabel(":"),
		destPort_Entry,
		layout.NewSpacer(),
		okBtn,
		delBtn,
		layout.NewSpacer(),
	}
	inputRow := widget.NewHBox(
		obs...
	)
	box.Children = append(box.Children, inputRow)
	box.Refresh()
	isAdd = true
}

func deleteInputRow(box *widget.Box, key string) {
	delete(DataMap,key)
	refreshWindow(box)
}

func refreshWindow(box *widget.Box) {
	box.Children = nil
	box.Refresh()
	loadInputRow(box)
	isAdd = false
}

func startTcpProxy() {
	if !proxy.IsStart {
		proxys := []string{}
		for _, p := range DataMap {
			s := fmt.Sprintf("%v,%v:%v", p.SourcePort, p.DestIp, p.DestPort)
			proxys = append(proxys, s)
		}
		proxy.SetProxys(proxys)
		proxy.Start()
		startBtn.Text = "stop"
	} else {
		proxy.Stop()
		startBtn.Text = "start"
	}
}
