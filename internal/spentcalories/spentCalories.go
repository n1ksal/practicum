package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {

	// Проверил, чтобы длина слайса была равна 3
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных")
	}

	// Преобразовал первый элемент слайса в тип int.
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	activity := parts[1]

	// Преобразовал третий элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	//вернул количество шагов, вид активности, продолжительность и nil
	return steps, activity, duration, nil
}

func distance(steps int) float64 {
	// Для вычисления дистанции умножил шаги
	// на длину шага lenStep и разделил на mInKm
	return float64(steps) * lenStep / mInKm
}

func meanSpeed(steps int, duration time.Duration) float64 {
	// Проверил, что продолжительность duration больше 0.
	if duration <= 0 {
		return 0
	}

	// Вычислил дистанцию с помощью distance().
	dist := distance(steps)
	// Вычислил и вернул среднюю скорость.
	hours := duration.Hours()
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) string {
	// Получил значения из строки данных с помощью функции parseTraining(),
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка парсинга данных: %v", err)
	}

	// Проверил, какой вид тренировки был передан в строке, которую парсил.
	// Для каждого из видов тренировки рассчитал дистанцию, среднюю скорость и калории.
	var distance, speed, calories float64

	switch activity {
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
		distance = float64(steps) * lenStep / 1000
		speed = meanSpeed(steps, duration)

	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
		distance = float64(steps) * lenStep / 1000
		speed = meanSpeed(steps, duration)

	default:
		return "неизвестный тип тренировки"
	}

	durationInHours := duration.Seconds() / 3600

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, durationInHours, distance, speed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// Рассчитал среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, duration)
	// Рассчитал и вернул количество калорий.
	return ((runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight)
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// Рассчитал среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, duration)
	// Продолжительность duration перевел в часы
	hours := duration.Seconds() / 3600
	// Рассчитал и вернул количество калорий.
	return ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * hours * minInH
}
