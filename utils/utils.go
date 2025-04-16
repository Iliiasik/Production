package utils

import (
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Возвращает путь к файлу и номер строки, где была вызвана функция (удобно для логов)
func CallerLocation(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "неизвестно"
	}
	return file + ":" + strconv.Itoa(line)
}

// Достаёт только читаемую часть ошибки (после ": ")
// Также округляем значение до 2 знаков
func ParseSQLErrorMessage(fullMsg string) string {
	if idx := strings.Index(fullMsg, " (SQLSTATE"); idx != -1 {
		fullMsg = fullMsg[:idx]
	}
	if idx := strings.Index(fullMsg, ": "); idx != -1 {
		fullMsg = fullMsg[idx+2:]
	}
	fullMsg = strings.TrimSpace(fullMsg)

	// Округляем длинные числа после точки (например, 9759.06792267... -> 9759.07)
	re := regexp.MustCompile(`(\d+\.\d{2})\d+`)
	fullMsg = re.ReplaceAllString(fullMsg, "$1")

	return fullMsg
}
