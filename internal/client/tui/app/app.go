// Package app содержит основные компоненты и логику для текстового пользовательского интерфейса (TUI) приложения.
package app

import (
	"errors"
	"github.com/charmbracelet/bubbletea"
	"github.com/romanp1989/gophkeeper/internal/client/config"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/tui/top"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TuiApplication структура, представляющая приложение с текстовым интерфейсом.
type TuiApplication struct {
	client   grpc.ClientGRPCInterface
	config   *config.Config
	logger   *zap.Logger
	notify   chan error
	program  *tea.Program
	topModel *top.Model
}

// NewTuiApplication создаёт новый экземпляр TuiApplication с заданной конфигурацией и логгером.
// Эта функция также инициализирует модель интерфейса, обрабатывает ошибки инициализации.
func NewTuiApplication(grpcClient grpc.ClientGRPCInterface, config *config.Config, logger *zap.Logger) *TuiApplication {
	model, err := top.NewModel(config, grpcClient)
	if err != nil {
		logger.Fatal("Error initializing model", zap.Error(err))
	}
	return &TuiApplication{
		client:   grpcClient,
		config:   config,
		logger:   logger,
		notify:   make(chan error, 1),
		program:  tea.NewProgram(model, tea.WithAltScreen()),
		topModel: model,
	}
}

// Start запускает TUI-приложение и его компоненты. Этот метод запускает уведомления клиента,
// программирует и обрабатывает системные сигналы для грациозного завершения работы.
func (a *TuiApplication) Start() {
	ErrExitCmd := errors.New("exit command")

	go func() {
		_, err := a.program.Run()
		if err != nil {
			log.Fatal("failed to run tui app: ", err)
		}
		a.notify <- ErrExitCmd
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-quit:
		a.logger.Info("interrupt: signal " + sig.String())
	case err := <-a.Notify():
		if errors.Is(err, ErrExitCmd) {
			a.shutdown()
		} else {
			a.logger.Error("TUI app error", zap.Error(err))
		}
	}
}

// Notify возвращает канал для отправки уведомлений об ошибках или командах выхода.
func (a *TuiApplication) Notify() chan error {
	return a.notify
}

// Останавливает приложение и выполняет необходимые действия для грациозного завершения работы.
func (a *TuiApplication) shutdown() {
	a.logger.Info("Shutting down application...")
	a.logger.Info("Application shutdown complete.")
}
