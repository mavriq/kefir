package ui

import (
	"fmt"
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

func (h *YAMLPreviewHandler) NotifyLeftPanelUpdate(index int, list *tview.List) {
	// Не требуется для YAML preview
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
}

func NewMountOptionsHandler(options []MountOption) *MountOptionsHandler {
	return &MountOptionsHandler{
		options:    options,
		currentIdx: -1,
	}
}

func (h *MountOptionsHandler) CreatePanel() tview.Primitive {
	// Создаем форму
	h.mountPathField = tview.NewInputField().
		SetLabel("Mount Path: ").
		SetFieldWidth(40).
		SetChangedFunc(func(text string) {
			if h.currentIdx >= 0 && h.currentIdx < len(h.options) {
				h.options[h.currentIdx].SetMountPath(text)
				h.updateDescription()
				h.notifyLeftPanelUpdate()
			}
		})

	h.readOnlyCheckbox = tview.NewCheckbox().
		SetLabel("Read Only: ").
		SetChangedFunc(func(checked bool) {
			if h.currentIdx >= 0 && h.currentIdx < len(h.options) {
				h.options[h.currentIdx].SetReadOnly(checked)
				h.updateDescription()
				h.notifyLeftPanelUpdate()
			}
		})

	h.descriptionView = tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true)

	// Создаем layout для формы
	h.formFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(h.mountPathField, 1, 0, true).
		AddItem(h.readOnlyCheckbox, 1, 0, false).
		AddItem(h.descriptionView, 0, 1, false)

	// Обертываем в Frame
	frame := tview.NewFrame(h.formFlex).
		SetBorders(1, 1, 0, 0, 2, 2)
	frame.SetBorder(true).
		SetTitle(" Mount Configuration ").
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
	switch event.Key() {
	case tcell.KeyTab:
		// Переключение между полями формы
		currentFocus := app.GetFocus()
		if currentFocus == h.mountPathField {
			app.SetFocus(h.readOnlyCheckbox)
		} else if currentFocus == h.readOnlyCheckbox {
			app.SetFocus(h.mountPathField)
		}
		return nil
	case tcell.KeyEnter:
		// Подтверждение ввода в поле
		return nil
	}
	return event
}

func (h *MountOptionsHandler) GetResult() interface{} {
	return h.options // Возвращаем модифицированные options
}

func (h *MountOptionsHandler) NotifyLeftPanelUpdate(index int, list *tview.List) {
	h.list = list
}

// Вспомогательные методы
func (h *MountOptionsHandler) updateDescription() {
	if h.currentIdx >= 0 && h.currentIdx < len(h.options) {
		desc := h.options[h.currentIdx].Description()
		h.descriptionView.SetText(desc)
	}
}

func (h *MountOptionsHandler) notifyLeftPanelUpdate() {
	if h.list != nil && h.currentIdx >= 0 && h.currentIdx < len(h.options) {
		// Обновляем текст в списке
		mainText := h.options[h.currentIdx].Title()
		h.list.SetItemText(h.currentIdx, mainText, "")
	}
}

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
