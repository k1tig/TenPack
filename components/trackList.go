package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	selectTrack key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		selectTrack: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Select Track"),
		),
	}
}

var items = []list.Item{
	item{title: "Free Race", desc: "Split all gates"},
	item{title: "FMV - LazE Classic", desc: "By: LazyDay & MrE"},
	item{title: "FMV - Multirotor Vermont July 25", desc: "By: Fastodon"},
	item{title: "MultiGP Canadian Nationals 2025 sub250", desc: "By: Solace & Nailed"},
}

type track struct {
	name   string
	author string
	gates  []int
}

var tracks = []track{
	{name: "Free Race", author: "Split all gates", gates: []int{}},
	{name: "FMV - LazE Classic", author: "LazyDay & MrE", gates: []int{3, 6, 9}},
	{name: "FMV - Multirotor Vermont July 25", author: "Fastodon", gates: []int{4, 2, 0}},
	{name: "FMultiGP Canadian Nationals 2025 sub250", author: "Solace & Nailed", gates: []int{6, 9}},
}

func initTrackList() model {
	m := model{gateOptions: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.gateOptions.Title = "Select Track"
	listKeys := newListKeyMap()
	m.keys = listKeys
	return m
}
