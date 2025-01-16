package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ivanov-nikolay/game/internal/entity"
	"github.com/ivanov-nikolay/game/internal/service"
)

// Decorator добавляет middleware к обработчикам
type Decorator func(http.Handler) http.Handler

// LifeStates объект для хранения состояния игры
type LifeStates struct {
	service.LifeService
	Height int
	Width  int
	Fill   int
}

// New создает новый объект
func New(ctx context.Context, height, width int, lifeService service.LifeService) (http.Handler, error) {
	serveMux := http.NewServeMux()

	lifeState := LifeStates{
		LifeService: lifeService,
		Height:      height,
		Width:       width,
		Fill:        40,
	}

	serveMux.HandleFunc("/nextstate", lifeState.nextState)
	serveMux.HandleFunc("/setstate", lifeState.setState)
	serveMux.HandleFunc("/reset", lifeState.reset)

	return serveMux, nil
}

// Decorate функция добавления middleware
func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	decorated := next
	for d := len(ds) - 1; d >= 0; d-- {
		decorated = ds[d](decorated)
	}

	return decorated
}

func (ls *LifeStates) nextState(w http.ResponseWriter, r *http.Request) {
	worldState := ls.LifeService.NewState() // Получаем очередное состояние игры

	err := json.NewEncoder(w).Encode(worldState.Cells)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ls *LifeStates) setState(w http.ResponseWriter, r *http.Request) {
	var req entity.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Fill < 0 || req.Fill > 100 {
		http.Error(w, "fill out of range", http.StatusBadRequest)
		return
	}

	ls.Fill = req.Fill

	file, err := os.OpenFile("state.cfg", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	percentage := strconv.Itoa(req.Fill)
	if _, err = file.WriteString(percentage + "%"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Пересоздаем LifeService c новым состоянием
	newLifeService, err := service.New(ls.Height, ls.Width, ls.Fill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ls.LifeService = *newLifeService
	worldState := ls.LifeService.NewState()
	err = json.NewEncoder(w).Encode(worldState.Cells)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (ls *LifeStates) reset(w http.ResponseWriter, r *http.Request) {
	var resp entity.Response

	file, err := os.Open("state.cfg")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fill := strings.TrimSpace(strings.ReplaceAll(string(content), "%", ""))

	resp.Fill, err = strconv.Atoi(fill)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
