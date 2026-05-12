package run

import (
	"context"
	"fmt"

	"kefir/pkg/ui"
)

// SelectMounts - выбор монтируемых ui.MountOption-объектов
func SelectMounts(ctx context.Context, mountOptions []ui.MountOption) ([]ui.MountOption, error) {
	if len(mountOptions) == 0 {
		return mountOptions, nil
	}

	// Создаем handler для правой панели (форма редактирования)
	var handler ui.RightPanelHandler = ui.NewMountOptionsHandler(mountOptions)

	// Конфигурация окна
	config := ui.WindowConfig{
		Title:        SelectMountsTitleText,
		Instructions: SelectMountsInstructions,
		InitialFocus: ui.FocusLeft,
	}

	// Создаем и запускаем окно
	window := ui.NewTwoPanelWindow(config, mountOptions, handler)
	_, err := window.Run(ctx)
	if err != nil {
		if err.Error() == "selection cancelled" {
			// Для SelectMounts это нормальный выход
			return mountOptions, nil
		}
		return mountOptions, fmt.Errorf("%s: %w", ConfigureMountsErrorCancelled, err)
	}

	// Для SelectMounts result будет []ui.MountOption
	// Но поскольку мы работали с оригинальными объектами,
	// они уже модифицированы, просто возвращаем исходный слайс
	return mountOptions, nil
}
