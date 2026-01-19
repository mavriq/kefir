package ui

import "github.com/gdamore/tcell/v2"

// Цвета приложения
const (
	ColorBackground = tcell.ColorBlack
	ColorTitle      = tcell.ColorWhite
	ColorActive     = tcell.ColorYellow
	ColorInactive   = tcell.ColorGray
	ColorBorder     = tcell.ColorWhite
	ColorBorderDim  = tcell.ColorGray
	ColorText       = tcell.ColorWhite
	ColorHighlight  = tcell.ColorGreen
	ColorSelected   = tcell.ColorBlue
	ColorError      = tcell.ColorRed
	ColorSuccess    = tcell.ColorGreen
)

// Цвета для YAML подсветки
const (
	YAMLColorKey        = tcell.ColorGreen
	YAMLColorApiVersion = tcell.ColorYellow
	YAMLColorKind       = tcell.ColorYellow
	YAMLColorMetadata   = tcell.ColorAqua
	YAMLColorSpec       = tcell.ColorPurple
	YAMLColorStatus     = tcell.ColorRed
	YAMLColorString     = tcell.ColorYellow
	YAMLColorNumber     = tcell.ColorCyan
	YAMLColorBoolean    = tcell.ColorMagenta
)
