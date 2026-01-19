package run

import (
	"context"
	"fmt"

	"kefir/pkg/ui"
)

// SelectMounts - выбор монтируемых ui.MountOption-объектов
func SelectMounts(ctx context.Context, mountOptions []*MountOptionImpl) ([]*MountOptionImpl, error) {
	if len(mountOptions) == 0 {
		return mountOptions, nil
	}

	// Конвертируем в интерфейс ui.MountOption
	options := make([]ui.MountOption, len(mountOptions))
	for i, opt := range mountOptions {
		options[i] = opt
	}

	// Создаем handler для правой панели (форма редактирования)
	handler := ui.NewMountOptionsHandler(options)

	// Создаем имена для левой панели
	names := make([]string, len(options))
	for i, opt := range options {
		names[i] = opt.Title()
	}

	// Конфигурация окна
	config := ui.WindowConfig{
		Title:          SelectMountsTitleText,
		Instructions:   SelectMountsInstructions,
		EnableCheckbox: true, // Включаем checkbox режим
		InitialFocus:   ui.FocusLeft,
	}

	// Создаем и запускаем окно
	window := ui.NewTwoPanelWindow(config, names, handler)
	result, err := window.Run(ctx)
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
