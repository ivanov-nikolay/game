package life

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// World тип данных для представления сетки.
// Каждая клетка должна быть представлена двумя координатами,
// поэтому используется двумерный слайс или слайс из слайсов
type World struct {
	Height int      // Высота сетки
	Width  int      // Ширина сетки
	Cells  [][]bool // Значений у клетки может быть только два — живая или мёртвая, поэтому тип данных в клетках — bool.
}

// NewWorld выделение памяти под сетку
func NewWorld(height, width int) *World {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return &World{
		Height: height,
		Width:  width,
		Cells:  cells,
	}
}

// Next вычисляет следующее состояние клетки на основе текущего и количества живых соседей
func (w *World) Next(x, y int) bool {
	n := w.Neighbors(x, y)       // Получим количество живых соседей
	alive := w.Cells[y][x]       // Текущее состояние клетки
	if n < 4 && n > 1 && alive { // Если соседей двое или трое, а клетка жива,
		return true // то следующее её состояние — жива
	}
	if !alive && n == 3 { // Если клетка мертва, но у неё трое соседей,
		return true // клетка оживает
	}

	return false // В любых других случаях — клетка мертва
}

// NextState определяет полное состояние сетки
func NextState(oldWorld, newWorld *World) {
	// Переберём все клетки, чтобы понять, в каком они состоянии
	for i := 0; i < oldWorld.Height; i++ {
		for j := 0; j < oldWorld.Width; j++ {
			// Для каждой клетки получим новое состояние
			newWorld.Cells[i][j] = oldWorld.Next(j, i)
		}
	}
}

// Итак, мы подготовили функцию, которая изменяет состояние клеток. Но для того чтобы игра началась,
// ей нужно какое-то исходное состояние. Мы можем задать его вручную или написать отдельный метод Seed.

// Seed заполняет сетку живыми клетками в случайном порядке:
func (w *World) Seed() {
	// Снова переберём все клетки
	for _, row := range w.Cells {
		for i := range row {
			//rand.Intn(10) возвращает случайное число из диапазона	от 0 до 9
			if rand.Intn(10) == 1 {
				row[i] = true
			}
		}
	}
}

// Neighbors возвращает количество живых соседей клетки (x, y) на торроидальном поле
func (w *World) Neighbors(x, y int) int {
	count := 0
	deltas := []struct{ dx, dy int }{
		{-1, -1}, {0, -1}, {1, -1}, // Верхние соседи
		{-1, 0}, {1, 0}, // Слева и справа
		{-1, 1}, {0, 1}, {1, 1}, // Нижние соседи
	}

	for _, d := range deltas {
		nx, ny := (x+d.dx+w.Width)%w.Width, (y+d.dy+w.Height)%w.Height // Учтем границы сетки
		if w.Cells[ny][nx] {                                           // Проверим, жива ли клетка-сосед
			count++
		}
	}

	return count
}

// Neighbors возвращает количество живых соседей клетки (x, y) на плоском поле
//func (w *World) Neighbors(x, y int) int {
//	count := 0
//	deltas := []struct{ dx, dy int }{
//		{-1, -1}, {0, -1}, {1, -1}, // Верхние соседи
//		{-1, 0}, {1, 0}, // Соседи слева и справа
//		{-1, 1}, {0, 1}, {1, 1}, // Нижние соседи
//	}
//
//	for _, d := range deltas {
//		nx, ny := x+d.dx, y+d.dy
//
//		// Проверяем, чтобы сосед был внутри границ сетки и был живой
//		if nx >= 0 && nx < w.Width && ny >= 0 && ny < w.Height && w.Cells[ny][nx] {
//			count++
//		}
//	}
//
//	return count
//}

// RandInit заполняет сетку случайными состояниями в соответствии с заданным процентом живых клеток
func (w *World) RandInit(percentage int) {
	numAlive := percentage * w.Width * w.Height / 100
	w.fillAlive(numAlive)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < numAlive; i++ {
		randRowLeft := r.Intn(w.Width)
		randColLeft := r.Intn(w.Height)
		randRowRight := r.Intn(w.Width)
		randColRight := r.Intn(w.Width)

		w.Cells[randRowLeft][randColLeft] = w.Cells[randRowRight][randColRight]
	}
}

func (w *World) fillAlive(num int) {
	aliveCount := 0
	for j, row := range w.Cells {
		for k := range row {
			w.Cells[j][k] = true
			aliveCount++
			if aliveCount == num {
				return
			}
		}
	}
}

func (w *World) SaveState(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("could not open world state file: %v", err)
	}
	defer file.Close()

	for i, row := range w.Cells {
		rowString := ""
		for _, cell := range row {
			if cell {
				rowString += "1"
			} else {
				rowString += "0"
			}
		}

		// Добавляем перенос строки только если это не последняя строка
		if i < len(w.Cells)-1 {
			rowString += "\n"
		}

		if _, err := file.WriteString(rowString); err != nil {
			return fmt.Errorf("could not write world state file: %v", err)
		}
	}

	return nil
}

// LoadState считывает исходное состояние из файла и устанавливает размер сетки
func (w *World) LoadState(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("could not open world state file: %v", err)
	}
	defer file.Close()

	cells := make([][]bool, 0)
	width := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if width == 0 {
			width = len(line)
		} else if width != len(line) {
			return fmt.Errorf("world state file contains more than one line of length")
		}

		row := make([]bool, len(line))

		for i, cell := range line {
			switch cell {
			case '1':
				row[i] = true
			case '0':
				row[i] = false
			default:
				return fmt.Errorf("world state file contains invalid characters")
			}
		}

		cells = append(cells, row)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not read world state file: %v", err)
	}

	w.Height = len(cells)
	w.Width = width
	w.Cells = cells

	return nil
}

// String визуализирует процесс игры (Вариант №1)
func (w *World) String() string {
	var sb strings.Builder

	// Символы для отображения живых и мёртвых клеток
	brownSquare := "\xF0\x9F\x9F\xAB" // Коричневый квадрат
	greenSquare := "\xF0\x9F\x9F\xA9" // Зелёный квадрат

	for _, row := range w.Cells {
		for _, cell := range row {
			if cell {
				sb.WriteString(greenSquare) // Живая клетка
			} else {
				sb.WriteString(brownSquare) // Мёртвая клетка
			}
		}
		sb.WriteString("\n") // Переход на следующую строку
	}

	return sb.String()
}

//func (w *World) String() string { // (Вариант №2)
//	result := ""
//	for _, row := range w.Cells {
//		for _, cell := range row {
//			if cell {
//				result += "⬛" // Живая клетка
//			} else {
//				result += "⬜" // Мёртвая клетка
//			}
//		}
//		result += "\n"
//	}
//	return result
//}
