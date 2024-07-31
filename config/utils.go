package config

import (
	"net/http"
	"slices"
)

func ServiceCodeToHttpStatus(code ServiceCode) int {
	switch code {
	case CodeOK:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnprocessableEntity, CodeDatabaseError:
		return http.StatusUnprocessableEntity
	case CodeNotFound:
		return http.StatusNotFound
	case CodeDuplicateEvent:
		return http.StatusConflict
	case CodeOrderEnded:
		return http.StatusGone
	default:
		return http.StatusInternalServerError
	}
}

func IsFinalOrderStatus(status OrderStatus) bool {
	return slices.Contains(OrderFinalStatuses, status)
}
