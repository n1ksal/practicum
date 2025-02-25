package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделяю строку на слайс строк
	parts := strings.Split(data, ",")
	// Проверяю что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// Преобразую первый элемент в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	// Преобразую второй элемент в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	// Если всё прошло без ошибок, возвращаю количество шагов и продолжительность
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получил данные о количестве шагов
	// и продолжительности прогулки с помощью функции
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	// Вычислил дистанцию в метрах.
	distance := float64(steps) * StepLength / 1000
	// Вычислил количество калорий, потраченных на прогулке.
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	//Сформировать строку, которую буду возвращать
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distance, calories)
}
