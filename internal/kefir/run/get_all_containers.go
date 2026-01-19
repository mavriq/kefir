package run

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// ContainerSelection представляет результат выбора контейнера
type ContainerSelection struct {
	Container *corev1.Container
	IsInit    bool // true для init-контейнеров
	Index     int  // Оригинальный индекс в соответствующем слайсе
}

func (selection ContainerSelection) String() string {
	if selection.Container == nil {
		return GetOrSelectContainerOptionNone
	}

	i := ""
	if selection.IsInit {
		i = GetOrSelectContainerOptionInitPrefix
	}

	return fmt.Sprintf("%s#%i %s", i, selection.Index, selection.Container.Name)
}

// GetAllContainers возвращает все контейнеры в виде ContainerSelection
func GetAllContainers(pod *corev1.Pod) []ContainerSelection {
	initCount := len(pod.Spec.InitContainers)
	regularCount := len(pod.Spec.Containers)
	totalCount := initCount + regularCount

	result := make([]ContainerSelection, totalCount)

	// Init контейнеры
	for i := 0; i < initCount; i++ {
		result[i] = ContainerSelection{
			Container: &pod.Spec.InitContainers[i],
			IsInit:    true,
			Index:     i,
		}
	}

	// Обычные контейнеры
	for i := 0; i < regularCount; i++ {
		result[initCount+i] = ContainerSelection{
			Container: &pod.Spec.Containers[i],
			IsInit:    false,
			Index:     i,
		}
	}

	return result
}
