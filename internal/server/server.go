package server

import (
	"context"
	"errors"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/templating"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"time"

	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/config"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/handler"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/sender"
)

func GinSlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.Context().Value("start_time")
		if start == nil {
			start = time.Now()
			c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "start_time", start))
		}

		startTime := start.(time.Time)
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		logger.Info("HTTP request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", c.Request.URL.RawQuery,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"latency", latency.String(),
		)
	}
}

type Server struct {
	config *config.Config
	logger *slog.Logger
	engine *gin.Engine
}

func NewServer(cfg *config.Config, logger *slog.Logger, tm *templating.TemplateManager) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(GinSlogMiddleware(logger))

	var senders []sender.Sender

	if cfg.WardConfig.Enable {
		logger.Info("Ward sender is enabled", "config", cfg.WardConfig)
		wardSender := sender.WardSender{
			Config: cfg.WardConfig,
			Tm:     tm,
		}
		senders = append(senders, &wardSender)
	}

	if cfg.DoorayConfig.Enable {
		logger.Info("Dooray sender is enabled", "config", cfg.DoorayConfig)
		dooraySender := sender.DooraySender{
			Config: cfg.DoorayConfig,
			Tm:     tm,
		}
		senders = append(senders, &dooraySender)
	}

	if senders == nil || len(senders) == 0 {
		logger.Error("shoudd enable at least one sender, please check your config", "config", cfg)
		os.Exit(1)
	}

	s := &Server{
		config: cfg,
		logger: logger,
		engine: r,
	}

	s.initRoutes(senders)

	return s
}

func (s *Server) initRoutes(senders []sender.Sender) {
	s.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// AlertHandler: 템플릿 매니저(s.tm)와 senders를 주입
	alertHandler := handler.NewAlertHandler(s.logger, senders)
	s.engine.POST("/webhook", alertHandler.HandleWebhook)

	s.engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//// 추가로 /dooray, /ward 엔드포인트를 추가해서 각각 전송하게끔 해야겠다.
	// 아니면, webhook 으로 오면 sender에 있는 전체 전송 , webhook/dooray, webhook/ward 로 오면 각각 전송하게끔 코드 설계
	s.engine.POST("/webhook/dooray")

}

func (s *Server) Run() {
	s.logger.Info("Starting server", "listen_address", s.config.ListenAddress)
	srv := &http.Server{
		Addr:    s.config.ListenAddress,
		Handler: s.engine,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
