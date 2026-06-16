package ui

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TwoPanelWindow - универсальное двухпанельное окно
type TwoPanelWindow struct {
	app          *tview.Application
	config       WindowConfig
	leftList     *tview.List
	leftFrame    *tview.Frame
	rightPanel   tview.Primitive
	rightFrame   *tview.Frame
	rightHandler RightPanelHandler
	resultChan   chan interface{}
	errorChan    chan error

	// Для checkbox режима (ConfigureMounts)
	selectedIndices map[int]bool

	// Оригинальные тексты элементов
	// originalTexts []string
	items []RightPanelTitledElem
}

// NewTwoPanelWindow создает новое двухпанельное окно
func NewTwoPanelWindow[T RightPanelTitledElem](config WindowConfig, rawItems []T, handler RightPanelHandler) *TwoPanelWindow {
	// var items []RightPanelTitledElem = snippet.ConvertSlice[RightPanelTitledElem, RightPanelTitledElem](rawItems)

	items := make([]RightPanelTitledElem, len(rawItems))
	for i, e := range rawItems {
		items[i] = RightPanelTitledElem(e)
	}

	app := tview.NewApplication()

	// Создаем левую панель
	leftList, leftFrame := CreateListPanel(config.Title)

	// Создаем правую панель
	rightPanel := handler.CreatePanel()
	rightFrame := tview.NewFrame(rightPanel).
		SetBorders(1, 1, 0, 0, 2, 2)
	rightFrame.SetBorder(true).
		// SetTitle(" Details ").
		SetTitleAlign(tview.AlignLeft)

	window := &TwoPanelWindow{
		app:             app,
		config:          config,
		leftList:        leftList,
		leftFrame:       leftFrame,
		rightPanel:      rightPanel,
		rightFrame:      rightFrame,
		rightHandler:    handler,
		resultChan:      make(chan interface{}, 1),
		errorChan:       make(chan error, 1),
		selectedIndices: make(map[int]bool),
		items:           items,
	}

	// наполняем список в левой панели
	window.FillLeftList()

	// Настраиваем обработчики
	window.setupHandlers()

	return window
}

func (w *TwoPanelWindow) FillLeftList() {
	w.leftList.Clear()
	// Заполняем список
	for i, item := range w.items {
		w.leftList.AddItem(item.Title(), "", rune('1'+i), nil)
	}
}

func (w *TwoPanelWindow) setupHandlers() {
	// Обработчик изменения выбора в списке
	w.leftList.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		w.rightHandler.UpdateContent(index)
	})

	// Обработчик выбора (Enter)
	w.leftList.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		// if w.config.EnableCheckbox {
		// 	// В режиме checkbox - переключаем выбор
		// 	w.toggleSelection(index)
		// } else {
		//	// В обычном режиме - выбираем элемент

		w.selectItem(index)
		// }
	})

	// // Обработчик ввода в левой панели
	// w.leftList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if w.config.EnableCheckbox && event.Key() == tcell.KeyRune && event.Rune() == ' ' {
	// 		// Space для выбора в checkbox режиме
	// 		index := w.leftList.GetCurrentItem()
	// 		w.toggleSelection(index)
	// 		return nil
	// 	}
	// 	return event
	// })

	// Глобальная обработка клавиш
	w.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			w.switchFocus()
			return nil

		case tcell.KeyEsc:
			w.cancelSelection()
			return nil

		case tcell.KeyCtrlC:
			w.cancelSelection()
			return nil
		}

		// Передаем события в правую панель
		currentFocus := w.app.GetFocus()
		if currentFocus != w.leftList {
			return w.rightHandler.HandleInput(event, w.app)
		}

		return event
	})
}

func (w *TwoPanelWindow) switchFocus() {
	currentFocus := w.app.GetFocus()
	if currentFocus == w.leftList {
		// Переключаемся на правую панель
		w.app.SetFocus(w.rightPanel)
		w.leftFrame.SetBorderColor(ColorBorderDim)
		w.rightFrame.SetBorderColor(ColorActive)
	} else {
		// Переключаемся на левую панель
		w.app.SetFocus(w.leftList)
		w.leftFrame.SetBorderColor(ColorActive)
		w.rightFrame.SetBorderColor(ColorBorderDim)
	}
}

func (w *TwoPanelWindow) selectItem(index int) {
	result := SelectionResult{
		Index: index,
		Item:  w.rightHandler.GetResult(),
	}
	w.resultChan <- result
	w.app.Stop()
}

// func (w *TwoPanelWindow) toggleSelection(index int) {
// 	if _, selected := w.selectedIndices[index]; selected {
// 		delete(w.selectedIndices, index)
// 		w.leftList.SetItemText(index, "□ "+w.getOriginalItemText(index), "")
// 	} else {
// 		w.selectedIndices[index] = true
// 		w.leftList.SetItemText(index, "✓ "+w.getOriginalItemText(index), "")
// 	}

// 	// Обновляем правую панель
// 	w.rightHandler.UpdateContent(index)
// 	w.rightHandler.NotifyLeftPanelUpdate(index, w.leftList)
// }

func (w *TwoPanelWindow) getOriginalItemText(index int) string {
	if index < 0 || index >= len(w.items) {
		return ""
	}
	return w.items[index].Title()
}

func (w *TwoPanelWindow) cancelSelection() {
	// if w.config.EnableCheckbox {
	// 	// Для ConfigureMounts возвращаем модифицированные данные
	// 	w.resultChan <- w.rightHandler.GetResult()
	// } else {
	// Для окон выбора возвращаем ошибку отмены
	w.errorChan <- fmt.Errorf(ErrorSelectionCancelled)
	// }
	w.app.Stop()
}

// Run запускает окно и ожидает результат
func (w *TwoPanelWindow) Run(ctx context.Context) (interface{}, error) {

	rightPanelChangedChan := w.rightHandler.Watch(ctx)
	go func() {
		for range rightPanelChangedChan {
			w.app.QueueUpdateDraw(func() {
				w.FillLeftList()
			})
		}
	}()

	// Создаем layout
	panelsFlex := tview.NewFlex().
		AddItem(w.leftFrame, 0, 1, true).  // Левая панель, изначально в фокусе
		AddItem(w.rightFrame, 0, 2, false) // Правая панель

	// Создаем заголовок и инструкции
	titleView := CreateTitle(w.config.Title)
	instructionsView := CreateInstructions(w.config.Instructions)

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titleView, 1, 0, false).       // Заголовок
		AddItem(panelsFlex, 0, 1, true).       // Основное содержимое
		AddItem(instructionsView, 1, 0, false) // Подвал

	// Настраиваем цвета
	mainFlex.SetBackgroundColor(ColorBackground)
	w.leftFrame.SetBorderColor(ColorActive)     // Левая панель активна изначально
	w.rightFrame.SetBorderColor(ColorBorderDim) // Правая панель неактивна

	go func() {
		if err := w.app.SetRoot(mainFlex, true).EnableMouse(false).Run(); err != nil {
			w.errorChan <- fmt.Errorf("ui error: %w", err)
		}
	}()

	// Обновляем правую панель для первого элемента
	if w.leftList.GetItemCount() > 0 {
		w.rightHandler.UpdateContent(0)
	}

	// Ждем результат с учетом контекста
	select {
	case <-ctx.Done():
		w.app.Stop()
		return nil, fmt.Errorf("%s: %w", ErrorUICancelled, ctx.Err())

	case result := <-w.resultChan:
		return result, nil

	case err := <-w.errorChan:
		return nil, err
	}
}
