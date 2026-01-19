package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FocusSide int

const (
	FocusLeft FocusSide = iota
	FocusRight
)

// Константы для стандартных сообщений об ошибках
const (
	ErrorSelectionCancelled = "selection cancelled"
	ErrorUICancelled        = "ui operation cancelled"
)

// WindowConfig - конфигурация окна
type WindowConfig struct {
	Title          string    // Заголовок окна
	Instructions   string    // Инструкции в подвале
	EnableCheckbox bool      // Включить checkbox режим (для множественного выбора)
	InitialFocus   FocusSide // Какая панель в фокусе при старте
}

// RightPanelHandler - интерфейс для поведения правой панели
type RightPanelHandler interface {
	// Создание компонента правой панели
	CreatePanel() tview.Primitive

	// Обновление контента при выборе элемента
	UpdateContent(index int)

	// Обработка ввода в правой панели
	HandleInput(event *tcell.EventKey, app *tview.Application) *tcell.EventKey

	// Получение результата (для окон с выбором)
	GetResult() interface{}

	// Обновление левой панели при изменениях (для ConfigureMounts)
	NotifyLeftPanelUpdate(index int, list *tview.List)
}

// MountOption - интерфейс для ConfigureMounts
type MountOption interface {
	Title() string             // Название, коротко описывающее Volume
	Description() string       // Детальное описание Volume
	GetMountPath() string      // Вернуть точку монтирования для этого Volume в создаваемом эфемерном контейнере
	SetMountPath(path string)  // Задать точку монтирования для этого Volume в создаваемом эфемерном контейнере
	GetReadOnly() bool         // Вернуть режим монтирования RO/RW
	SetReadOnly(readOnly bool) // Задать режим монтирования RO/RW
	GetName() string           // имя Volume (настраивается при создании объекта. Не может быть перенастроено)
}

// SelectionResult - результат выбора для окон 1 и 2
type SelectionResult struct {
	Index int
	Item  interface{}
}
