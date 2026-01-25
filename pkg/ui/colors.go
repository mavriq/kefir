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
	YAMLColorKey        = tcell.ColorGreen   // Зелёный - для ключей
	YAMLColorApiVersion = tcell.ColorYellow  // Жёлтый - для apiVersion/kind
	YAMLColorKind       = tcell.ColorYellow  // Жёлтый
	YAMLColorMetadata   = tcell.ColorAqua    // Бирюзовый - для metadata
	YAMLColorSpec       = tcell.ColorFuchsia // Фуксия - для spec
	YAMLColorStatus     = tcell.ColorRed     // Красный - для status
	YAMLColorString     = tcell.ColorYellow  // Жёлтый - для строк
	YAMLColorNumber     = tcell.ColorAqua    // Бирюзовый - для чисел
	YAMLColorBoolean    = tcell.ColorFuchsia // Фуксия - для true/false
)
