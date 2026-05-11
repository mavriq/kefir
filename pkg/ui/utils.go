package ui

import "fmt"

// Вспомогательные утилиты

// StringSliceToInterface конвертирует []string в []interface{}
func StringSliceToInterface(strings []string) []interface{} {
	result := make([]interface{}, len(strings))
	for i, s := range strings {
		result[i] = s
	}
	return result
}

// MapToTitleStrings создает []string из объектов с методом Title()
func MapToTitleStrings(items []interface{}) []string {
	result := make([]string, len(items))
	for i, item := range items {
		if withTitle, ok := item.(interface{ Title() string }); ok {
			result[i] = withTitle.Title()
		} else {
			result[i] = fmt.Sprintf("Item %d", i)
		}
	}
	return result
}

type permanentRightPanelTitledElem struct {
	title string
}

func (e *permanentRightPanelTitledElem) Title() string {
	return e.title
}

func PermanentRightPanelTitledElem(title string) RightPanelTitledElem {
	return &permanentRightPanelTitledElem{title: title}
}
