package handler

import (
	"fmt"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/sender"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/templating"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/pkg/models"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"sync"
)

type AlertHandler struct {
	logger  *slog.Logger
	senders []sender.Sender
	tm      *templating.TemplateManager // 템플릿 매니저 추가
}

func NewAlertHandler(logger *slog.Logger, senders []sender.Sender) *AlertHandler {
	return &AlertHandler{
		logger:  logger,
		senders: senders,
	}
}

func (h *AlertHandler) HandleWebhook(c *gin.Context) {
	var data models.AlertmanagerData

	if err := c.ShouldBindJSON(&data); err != nil {
		h.logger.Error("failed to bind JSON", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	h.logger.Info("received request", slog.String("receiver", data.Receiver))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var sendErrors []string
	for _, s := range h.senders {
		wg.Add(1)
		go func(sender sender.Sender) {
			defer wg.Done()
			if err := sender.Send(&data); err != nil {
				mu.Lock()
				sendErrors = append(sendErrors, fmt.Sprintf("[%s] send failed: %v", sender.Name(), err))
				mu.Unlock()
				h.logger.Error("send failed", slog.String("sender", sender.Name()), slog.Any("error", err))
			} else {
				h.logger.Info("send succeeded", slog.String("sender", sender.Name()))
			}
		}(s)
	}
	wg.Wait()

	if len(sendErrors) > 0 {
		for cnt, err := range sendErrors {
			h.logger.Error("send error", slog.Int("index", cnt+1), slog.Int("total", len(sendErrors)), slog.String("message", err))
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "send failed", "details": sendErrors})
	} else {
		h.logger.Info("all sends succeeded")
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	}
}
