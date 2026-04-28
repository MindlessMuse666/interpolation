package interpolation

import (
	"math"
	"testing"
)

func TestLinearInterpolation(t *testing.T) {
	points := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 4},
	}

	// Тест внутри интервала
	val, curve, err := LinearInterpolation(points, 0.5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if math.Abs(val-0.5) > 1e-9 {
		t.Errorf("Expected 0.5, got %v", val)
	}
	if len(curve) != curvePointsCount {
		t.Errorf("Expected %d curve points, got %d", curvePointsCount, len(curve))
	}

	// Тест на границе
	val, _, _ = LinearInterpolation(points, 1.5)
	if math.Abs(val-2.5) > 1e-9 {
		t.Errorf("Expected 2.5, got %v", val)
	}

	// Тест ошибок
	_, _, err = LinearInterpolation([]Point{{X: 0, Y: 0}}, 0.5)
	if err != ErrNotEnoughPoints {
		t.Errorf("Expected ErrNotEnoughPoints, got %v", err)
	}

	_, _, err = LinearInterpolation([]Point{{X: 0, Y: 0}, {X: 0, Y: 1}}, 0.5)
	if err == nil || err.Error() == "" {
		t.Errorf("Expected duplicate X error")
	}
}

func TestLagrangeInterpolation(t *testing.T) {
	points := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 4},
	}

	// Для x^2 точки (0,0), (1,1), (2,4) дают точное совпадение
	val, curve, err := LagrangeInterpolation(points, 1.5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// 1.5^2 = 2.25
	if math.Abs(val-2.25) > 1e-9 {
		t.Errorf("Expected 2.25, got %v", val)
	}
	if len(curve) != curvePointsCount {
		t.Errorf("Expected %d curve points, got %d", curvePointsCount, len(curve))
	}

	// Тест на совпадение с узлом
	val, _, _ = LagrangeInterpolation(points, 2)
	if math.Abs(val-4) > 1e-9 {
		t.Errorf("Expected 4, got %v", val)
	}
}

func TestNewtonInterpolation(t *testing.T) {
	points := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 4},
	}

	val, curve, err := NewtonInterpolation(points, 1.5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if math.Abs(val-2.25) > 1e-9 {
		t.Errorf("Expected 2.25, got %v", val)
	}
	if len(curve) != curvePointsCount {
		t.Errorf("Expected %d curve points, got %d", curvePointsCount, len(curve))
	}

	// Тест на совпадение с узлом
	val, _, _ = NewtonInterpolation(points, 1)
	if math.Abs(val-1) > 1e-9 {
		t.Errorf("Expected 1, got %v", val)
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		points  []Point
		min     int
		wantErr error
	}{
		{"Empty", []Point{}, 2, ErrNotEnoughPoints},
		{"Too few", []Point{{0, 0}}, 2, ErrNotEnoughPoints},
		{"Duplicates", []Point{{0, 0}, {0, 1}}, 2, ErrDuplicateX},
		{"Valid", []Point{{0, 0}, {1, 1}}, 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePoints(tt.points, tt.min)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.wantErr)
				}
				// check if it wraps ErrDuplicateX
				if tt.wantErr == ErrDuplicateX {
					if !errors_Is(err, ErrDuplicateX) {
						t.Errorf("Expected error to wrap ErrDuplicateX, got %v", err)
					}
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper because I can't import errors inside the test easily without adding it to the file
func errors_Is(err, target error) bool {
	if err == target {
		return true
	}
	type wrapper interface {
		Unwrap() error
	}
	for {
		if w, ok := err.(wrapper); ok {
			err = w.Unwrap()
			if err == target {
				return true
			}
		} else {
			return false
		}
	}
}
