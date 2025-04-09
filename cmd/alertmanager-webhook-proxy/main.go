package main

import (
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/config"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/logging"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/server"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/templating"
	"os"
)

func main() {
	// 1. Config, Logger 설정
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := logging.NewLogger(cfg.LogLevel)

	// 2. TemplateManager 생성
	tm, err := templating.NewTemplateManager(cfg.TemplateDir)
	if err != nil {
		logger.Error("failed to load templates", "error", err)
		os.Exit(1)
	}

	// 3. Server 객체 생성 시 TemplateManager를 인자로 넘겨줌
	s := server.NewServer(cfg, logger, tm)
	s.Run()
}
