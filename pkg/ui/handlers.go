package ui

import (
	"context"
	"fmt"
	"kefir/pkg/snippet"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sigs.k8s.io/yaml"
)

// ==================== YAMLPreviewHandler ====================

// YAMLPreviewHandler - для SelectPod и SelectContainer
type YAMLPreviewHandler struct {
	data     []interface{}
	textView *tview.TextView
}

func NewYAMLPreviewHandler(data []interface{}) *YAMLPreviewHandler {
	return &YAMLPreviewHandler{data: data}
}

func (h *YAMLPreviewHandler) CreatePanel() tview.Primitive {
	textView, _ := CreateTextPanel("Details", true)
	h.textView = textView
	return textView
}

func (h *YAMLPreviewHandler) UpdateContent(index int) {
	if index < 0 || index >= len(h.data) {
		return
	}

	yamlBytes, err := yaml.Marshal(h.data[index])
	if err != nil {
		h.textView.SetText(fmt.Sprintf("Error: %v", err))
		return
	}

	yamlText := string(yamlBytes)
	yamlText = highlightYAML(yamlText)
	h.textView.SetText(yamlText)
	h.textView.ScrollToBeginning()
}

func (h *YAMLPreviewHandler) HandleInput(event *tcell.EventKey, app *tview.Application) *tcell.EventKey {
	// Пропускаем навигационные клавиши
	switch event.Key() {
	case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight,
		tcell.KeyPgUp, tcell.KeyPgDn, tcell.KeyHome, tcell.KeyEnd:
		return event
	}
	return event
}

func (h *YAMLPreviewHandler) GetResult() interface{} {
	return nil // Результат получается через SelectionResult
}

// func (h *YAMLPreviewHandler) NotifyLeftPanelUpdate(index int, list *tview.List) {
// 	// Не требуется для YAML preview
// }

func (h *YAMLPreviewHandler) Watch(ctx context.Context) <-chan struct{} {
	return neverWatch(ctx)
}

// ==================== MountOptionsHandler ====================

// MountOptionsHandler - для ConfigureMounts
type MountOptionsHandler struct {
	options    []MountOption
	currentIdx int
	list       *tview.List

	// Компоненты формы
	mountPathField   *tview.InputField
	readOnlyCheckbox *tview.Checkbox
	descriptionView  *tview.TextView

	// Контейнер для формы
	formFlex *tview.Flex

	// нужен для переключения фокуса между полями ввода в правой панели
	app *tview.Application
	// onChanged func(idx int, updatedTitle string) // колбэк для уведомления левой панели

	watchBroadcaster snippet.Broadcaster[struct{}]
}

func NewMountOptionsHandler(options []MountOption) *MountOptionsHandler {
	return &MountOptionsHandler{
		options:          options,
		currentIdx:       -1,
		watchBroadcaster: snippet.NewBroadcaster[struct{}](),
	}
}

func (h *MountOptionsHandler) CreatePanel() tview.Primitive {
	// Создаем форму
	h.mountPathField = tview.NewInputField().
		SetLabel("Mount Path: ").
		SetFieldWidth(40).
		SetChangedFunc(func(text string) {
			if h.currentIdx >= 0 && h.currentIdx < len(h.options) {
				// обновляем опцию
				h.options[h.currentIdx].SetMountPath(text)

				// Обновляем только описание и заголовок, не трогая фокус
				h.updateDescription()
				// h.notifyLeftPanelUpdate()
				h.watchBroadcaster.Publish(struct{}{})
			}
		})

	h.mountPathField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyEnter:
			h.app.SetFocus(h.readOnlyCheckbox) // Переходим на чекбокс
			return nil                         // Событие обработано
		}
		return event
	})

	h.readOnlyCheckbox = tview.NewCheckbox().
		SetLabel("Read Only: ").
		SetChangedFunc(func(checked bool) {
			if h.currentIdx >= 0 && h.currentIdx < len(h.options) {
				h.options[h.currentIdx].SetReadOnly(checked)
				h.updateDescription()
				// h.notifyLeftPanelUpdate()
				h.watchBroadcaster.Publish(struct{}{})
			}
		})

	h.readOnlyCheckbox.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			h.app.SetFocus(h.mountPathField) // Возвращаемся вверх
			return nil
		case tcell.KeyDown:
			h.app.SetFocus(h.descriptionView) // Идем к описанию
			return nil
		}
		return event
	})

	h.descriptionView = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)
	h.descriptionView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyUp {
			row, _ := h.descriptionView.GetScrollOffset()
			if row == 0 {
				h.app.SetFocus(h.readOnlyCheckbox)
				return nil
			}
		}
		return event
	})

	// Создаем layout для формы
	h.formFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(h.mountPathField, 1, 0, true).
		AddItem(h.readOnlyCheckbox, 1, 0, false).
		AddItem(h.descriptionView, 0, 1, false)

	// Обертываем в Frame
	frame := tview.NewFrame(h.formFlex).
		SetBorders(1, 1, 0, 0, 2, 2)
	frame.SetBorder(false).
		SetTitleAlign(tview.AlignLeft)

	return frame
}

func (h *MountOptionsHandler) UpdateContent(index int) {
	h.currentIdx = index

	if index < 0 || index >= len(h.options) {
		h.mountPathField.SetText("")
		h.readOnlyCheckbox.SetChecked(false)
		h.descriptionView.SetText("")
		return
	}

	option := h.options[index]
	h.mountPathField.SetText(option.GetMountPath())
	h.readOnlyCheckbox.SetChecked(option.GetReadOnly())
	h.updateDescription()

	// Устанавливаем фокус на первое поле
	h.mountPathField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
}

func (h *MountOptionsHandler) HandleInput(event *tcell.EventKey, app *tview.Application) *tcell.EventKey {
	h.app = app

	return event
}

func (h *MountOptionsHandler) GetResult() interface{} {
	return h.options // Возвращаем модифицированные options
}

// func (h *MountOptionsHandler) NotifyLeftPanelUpdate(index int, list *tview.List) {
// 	h.list = list
// }

func (h *MountOptionsHandler) Watch(ctx context.Context) <-chan struct{} {
	return h.watchBroadcaster.Subscribe(ctx, 1)
}

// Вспомогательные методы
func (h *MountOptionsHandler) updateDescription() {
	if h.currentIdx >= 0 && h.currentIdx < len(h.options) {
		desc := h.options[h.currentIdx].Description()
		h.descriptionView.SetText(desc)
	}
}

// func (h *MountOptionsHandler) notifyLeftPanelUpdate() {
// 	if h.list != nil && h.currentIdx >= 0 && h.currentIdx < len(h.options) {
// 		// Обновляем текст в списке
// 		mainText := h.options[h.currentIdx].Title()
// 		h.list.SetItemText(h.currentIdx, mainText, "")
// 	}
// }

// ==================== Вспомогательные функции ====================

func highlightYAML(yamlText string) string {
	lines := strings.Split(yamlText, "\n")
	var highlighted []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(trimmed, "apiVersion:"):
			line = strings.Replace(line, "apiVersion:", "[yellow]apiVersion:[white]", 1)
		case strings.HasPrefix(trimmed, "kind:"):
			line = strings.Replace(line, "kind:", "[yellow]kind:[white]", 1)
		case strings.HasPrefix(trimmed, "metadata:"):
			line = strings.Replace(line, "metadata:", "[aqua]metadata:[white]", 1)
		case strings.HasPrefix(trimmed, "spec:"):
			line = strings.Replace(line, "spec:", "[purple]spec:[white]", 1)
		case strings.HasPrefix(trimmed, "status:"):
			line = strings.Replace(line, "status:", "[red]status:[white]", 1)
		case strings.Contains(line, ":"):
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				line = fmt.Sprintf("[green]%s[white]:%s", parts[0], parts[1])
			}
		}

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

func neverWatch(ctx context.Context) <-chan struct{} {
	ret := make(chan struct{}, 1)

	go func() {
		defer close(ret)
		<-ctx.Done()
	}()

	return ret
}
