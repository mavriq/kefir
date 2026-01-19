package run

import "errors"

// Константы для SelectPod
const (
	GetOrSelectPodTitleText    = "📦 Select Kubernetes Pod"
	GetOrSelectPodInstructions = "↑↓: Navigate  Enter: Select  Tab: Switch panels  Esc: Cancel"
)

var (
	ErrorPodsAvailable              = errors.New("no pods available")
	ErrorPodSelectionCancelled      = errors.New("pod selection cancelled")
	ErrorInvalidSelectionResultType = errors.New("invalid selection result type")
	ErrorSelectionIndexOutOfRange   = errors.New("selection index out of range")
)

// Константы для выбора контейнера
const (
	GetOrSelectContainerTitleText         = "🐳 Select Container"
	GetOrSelectContainerInstructions      = "↑↓: Navigate  Enter: Select  Tab: Switch panels  Esc: Cancel"
	GetOrSelectContainerErrorCancelled    = "container selection cancelled"
	GetOrSelectContainerErrorNoContainers = "no containers available in pod"

	// Специальная опция "ничего не выбирать"
	GetOrSelectContainerOptionNone       = "- none -"
	GetOrSelectContainerOptionInitPrefix = "[init] "
)

// Константы для SelectMounts
const (
	SelectMountsTitleText      = "⚙️  Select Mount Options"
	SelectMountsInstructions   = "Space: Select  Tab: Switch  ↑↓: Navigate  Esc: Save & Exit"
	SelectMountsErrorCancelled = "mount configuration cancelled"
)
