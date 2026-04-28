package interpolation

import (
	"errors"
	"fmt"
	"sort"
)

// Point представляет точку в двумерном пространстве.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var (
	ErrNotEnoughPoints = errors.New("недостаточно точек для интерполяции")
	ErrDuplicateX      = errors.New("обнаружены дубликаты координат X")
)

const curvePointsCount = 100

// validatePoints проверяет входные данные на наличие ошибок.
func validatePoints(points []Point, minRequired int) error {
	if len(points) < minRequired {
		return ErrNotEnoughPoints
	}

	seenX := make(map[float64]bool)
	for _, p := range points {
		if seenX[p.X] {
			return fmt.Errorf("%w: %v", ErrDuplicateX, p.X)
		}
		seenX[p.X] = true
	}
	return nil
}

// generateCurvePoints создает массив точек для отрисовки графика.
func generateCurvePoints(points []Point, interpolate func(float64) float64) []Point {
	if len(points) == 0 {
		return nil
	}

	minX, maxX := points[0].X, points[0].X
	for _, p := range points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
	}

	// Добавляем небольшой отступ для красоты графика
	rangeX := maxX - minX
	if rangeX == 0 {
		rangeX = 1
	}
	start := minX - rangeX*0.1
	end := maxX + rangeX*0.1

	curve := make([]Point, curvePointsCount)
	step := (end - start) / float64(curvePointsCount-1)

	for i := 0; i < curvePointsCount; i++ {
		x := start + float64(i)*step
		curve[i] = Point{X: x, Y: interpolate(x)}
	}

	return curve
}

// LinearInterpolation реализует линейную интерполяцию.
// Возвращает интерполированное значение y, точки для графика и ошибку.
func LinearInterpolation(points []Point, x float64) (float64, []Point, error) {
	if err := validatePoints(points, 2); err != nil {
		return 0, nil, err
	}

	// Сортируем точки по X для линейной интерполяции
	sorted := make([]Point, len(points))
	copy(sorted, points)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].X < sorted[j].X
	})

	interpolate := func(targetX float64) float64 {
		// Находим интервал
		if targetX <= sorted[0].X {
			p0, p1 := sorted[0], sorted[1]
			return p0.Y + (targetX-p0.X)*(p1.Y-p0.Y)/(p1.X-p0.X)
		}
		if targetX >= sorted[len(sorted)-1].X {
			p0, p1 := sorted[len(sorted)-2], sorted[len(sorted)-1]
			return p0.Y + (targetX-p0.X)*(p1.Y-p0.Y)/(p1.X-p0.X)
		}

		for i := 0; i < len(sorted)-1; i++ {
			if targetX >= sorted[i].X && targetX <= sorted[i+1].X {
				p0, p1 := sorted[i], sorted[i+1]
				return p0.Y + (targetX-p0.X)*(p1.Y-p0.Y)/(p1.X-p0.X)
			}
		}
		return 0
	}

	resY := interpolate(x)
	curve := generateCurvePoints(sorted, interpolate)

	return resY, curve, nil
}

// LagrangeInterpolation реализует интерполяцию методом Лагранжа.
func LagrangeInterpolation(points []Point, x float64) (float64, []Point, error) {
	if err := validatePoints(points, 2); err != nil {
		return 0, nil, err
	}

	interpolate := func(targetX float64) float64 {
		var result float64
		n := len(points)
		for i := 0; i < n; i++ {
			term := points[i].Y
			for j := 0; j < n; j++ {
				if i != j {
					term *= (targetX - points[j].X) / (points[i].X - points[j].X)
				}
			}
			result += term
		}
		return result
	}

	resY := interpolate(x)
	curve := generateCurvePoints(points, interpolate)

	return resY, curve, nil
}

// NewtonInterpolation реализует интерполяцию методом Ньютона.
func NewtonInterpolation(points []Point, x float64) (float64, []Point, error) {
	if err := validatePoints(points, 2); err != nil {
		return 0, nil, err
	}

	n := len(points)
	// Вычисляем разделенные разности
	coefs := make([]float64, n)
	for i := 0; i < n; i++ {
		coefs[i] = points[i].Y
	}

	for j := 1; j < n; j++ {
		for i := n - 1; i >= j; i-- {
			coefs[i] = (coefs[i] - coefs[i-1]) / (points[i].X - points[i-j].X)
		}
	}

	interpolate := func(targetX float64) float64 {
		result := coefs[n-1]
		for i := n - 2; i >= 0; i-- {
			result = coefs[i] + (targetX-points[i].X)*result
		}
		return result
	}

	resY := interpolate(x)
	curve := generateCurvePoints(points, interpolate)

	return resY, curve, nil
}
