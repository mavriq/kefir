package run

import (
	"context"
	"fmt"

	"kefir/pkg/ui"
)

// SelectContainerWithInfo возвращает полную информацию о выбранном контейнере
// Возвращает ContainerSelection, где:
//   - Container == nil, IsInit == false, Index == -1: если выбрано "none"
//   - Container != nil: если выбран контейнер (с правильными IsInit и Index)
func SelectContainerWithInfo(ctx context.Context, allSelections []ContainerSelection) (*ContainerSelection, error) {
	totalCount := len(allSelections) + 1
	selectionsWithNone := make([]ContainerSelection, totalCount)

	// Первый элемент - пустой (опция "none")
	selectionsWithNone[0] = ContainerSelection{
		Container: nil,
		IsInit:    false,
		Index:     -1,
	}

	// Копируем остальные контейнеры
	copy(selectionsWithNone[1:], allSelections)

	// Подготавливаем данные для UI
	items := make([]interface{}, totalCount)
	names := make([]string, totalCount)

	for i, selection := range selectionsWithNone {
		if i == 0 {
			// Опция "none"
			items[i] = GetOrSelectContainerOptionNoneDesc
			names[i] = GetOrSelectContainerOptionNone
		} else {
			// Контейнеры
			items[i] = selection.Container
			names[i] = selection.String()
		}
	}

	// Создаем handler для правой панели
	handler := ui.NewYAMLPreviewHandler(items)

	// Конфигурация окна
	config := ui.WindowConfig{
		Title:          GetOrSelectContainerTitleText,
		Instructions:   GetOrSelectContainerInstructions,
		EnableCheckbox: false,
		InitialFocus:   ui.FocusLeft,
	}

	// Создаем и запускаем окно
	window := ui.NewTwoPanelWindow(config, names, handler)
	result, err := window.Run(ctx)
	if err != nil {
		if err.Error() == ui.ErrorSelectionCancelled {
			return nil, fmt.Errorf(GetOrSelectContainerErrorCancelled)
		}
		return nil, fmt.Errorf("failed to select container: %w", err)
	}

	// Приводим результат
	selection, ok := result.(ui.SelectionResult)
	if !ok {
		return nil, fmt.Errorf("invalid selection result type")
	}

	// Проверяем индекс
	idx := selection.Index
	if idx < 0 || idx >= len(selectionsWithNone) {
		return nil, fmt.Errorf("selection index out of range")
	}

	// Возвращаем выбранный элемент (может быть "none" или контейнер)
	selected := selectionsWithNone[idx]
	return &selected, nil
}
