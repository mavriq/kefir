package ui

import "github.com/rivo/tview"

// CreateTitle - создает заголовок окна
func CreateTitle(text string) *tview.TextView {
	view := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
	view.SetTextColor(ColorTitle)
	view.SetBackgroundColor(ColorBackground)
	return view
}

// CreateInstructions - создает подвал с инструкциями
func CreateInstructions(text string) *tview.TextView {
	view := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
	view.SetTextColor(ColorInactive)
	view.SetBackgroundColor(ColorBackground)
	return view
}

// CreateListPanel - создает панель со списком
func CreateListPanel(title string, enableCheckbox bool) (*tview.List, *tview.Frame) {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetMainTextColor(ColorText).
		SetSelectedBackgroundColor(ColorHighlight).
		SetSelectedTextColor(ColorBackground)

	frame := tview.NewFrame(list).
		SetBorders(1, 1, 0, 0, 2, 2)
	frame.SetBorder(true).
		SetTitle(" " + title + " ").
		SetTitleAlign(tview.AlignLeft)

	return list, frame
}

// CreateTextPanel - создает текстовую панель
func CreateTextPanel(title string, scrollable bool) (*tview.TextView, *tview.Frame) {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(false).
		SetWordWrap(false).
		SetScrollable(scrollable).
		SetTextColor(ColorText).
		SetBackgroundColor(ColorBackground)

	frame := tview.NewFrame(textView).
		SetBorders(1, 1, 0, 0, 2, 2)
	frame.SetBorder(true).
		SetTitle(" " + title + " ").
		SetTitleAlign(tview.AlignLeft)

	return textView, frame
}
