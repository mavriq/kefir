package run

import (
	"context"

	"kefir/pkg/ui"

	corev1 "k8s.io/api/core/v1"
)

func SelectPod(ctx context.Context, pods []corev1.Pod) (*corev1.Pod, error) {
	if len(pods) == 0 {
		return nil, ErrorPodsAvailable
	}

	// Конвертируем Pod'ы в интерфейсы для handler
	items := make([]interface{}, len(pods))
	for i := range pods {
		items[i] = &pods[i]
	}

	// Создаем handler для правой панели (YAML preview)
	handler := ui.NewYAMLPreviewHandler(items)

	// Создаем имена для левой панели
	names := make([]string, len(pods))
	for i, pod := range pods {
		names[i] = pod.Name
	}

	// Конфигурация окна
	config := ui.WindowConfig{
		Title:          GetOrSelectPodTitleText,
		Instructions:   GetOrSelectPodInstructions,
		EnableCheckbox: false, // Обычный выбор, не checkbox
		InitialFocus:   ui.FocusLeft,
	}

	// Создаем и запускаем окно
	window := ui.NewTwoPanelWindow(config, names, handler)
	result, err := window.Run(ctx)
	if err != nil {
		// Проверяем, это отмена пользователем или ошибка
		if err.Error() == "selection cancelled" {
			return nil, ErrorPodSelectionCancelled
		}
		return nil, err
	}

	// Приводим результат к SelectionResult
	selection, ok := result.(ui.SelectionResult)
	if !ok {
		return nil, ErrorInvalidSelectionResultType
	}

	// Возвращаем выбранный Pod
	if selection.Index < 0 || selection.Index >= len(pods) {
		return nil, ErrorSelectionIndexOutOfRange
	}

	return &pods[selection.Index], nil
}
