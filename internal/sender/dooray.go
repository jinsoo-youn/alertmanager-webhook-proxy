package sender

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/config"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/templating"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/pkg/models"
)

// 상수 선언
const (
	DoorayTemplateName = "dooray.tmpl"
	DoorayName         = "dooray"
)

type DooraySender struct {
	Config *config.DoorayConfig
	Tm     *templating.TemplateManager
}

// dooray.go에서 사용할 데이터 구조체 예: 템플릿 주입용
type doorayTmplData struct {
	Region   string
	Stage    string
	Actor    string
	StartsAt time.Time
	EndsAt   time.Time
	Color    string
	models.Alert
}

func (s *DooraySender) Name() string {
	return DoorayName
}

// Send: AlertmanagerData를 받아, alerts를 순회하며 dooray.tmpl을 렌더링 → POST
func (s *DooraySender) Send(data *models.AlertmanagerData) error {
	if len(data.Alerts) == 0 {
		return errors.New("no alerts found")
	}

	for i, alert := range data.Alerts {
		// (1) 필드 정규화 (severity → 색상 등)
		s.normalizeAlertFields(&alert)

		// (2) 템플릿에 주입할 데이터 구성
		tmplData := doorayTmplData{
			Region:   strings.ToUpper(s.Config.Region),
			Stage:    strings.ToUpper(s.Config.Stage),
			Actor:    "Alertmanager-Webhook-Proxy", // 필요 시 s.Config.Actor 등 사용
			StartsAt: alert.StartsAt,
			EndsAt:   alert.EndsAt,
			Color:    severityToColor(alert.Labels["severity"], alert.Status),
			Alert:    alert,
		}

		// (3) 템플릿 렌더링
		jsonBody, err := s.renderTemplate(tmplData)
		if err != nil {
			return fmt.Errorf("failed to render dooray template (alert #%d): %w", i, err)
		}

		// (4) HTTP POST
		if err := s.postToDooray(jsonBody); err != nil {
			return fmt.Errorf("failed to post dooray alert #%d: %w", i, err)
		}
	}
	return nil
}

// renderTemplate: dooray.tmpl을 읽어 최종 문자열(JSON or text) 생성
func (s *DooraySender) renderTemplate(data any) (string, error) {
	tmpl, ok := s.Tm.Get(DoorayTemplateName)
	if !ok {
		return "", fmt.Errorf("template not found: %s", DoorayTemplateName)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// postToDooray: 렌더링된 결과를 Dooray Webhook에 전송
func (s *DooraySender) postToDooray(body string) error {
	req, err := http.NewRequest("POST", s.Config.WebhookURL, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("dooray returned status %d", resp.StatusCode)
	}
	return nil
}

// normalizeAlertFields: Alert 라벨 혹은 상태를 사전에 변경 (Ward와 비슷한 개념)
func (s *DooraySender) normalizeAlertFields(alert *models.Alert) {
	// 예: severity 라벨 대문자화
	raw := alert.Labels["severity"]
	alert.Labels["severity"] = strings.ToUpper(raw)
	// 다른 라벨도 수정 가능
	//alert.Status = strings.ToUpper(alert.Status)
}

// severityToColor: 문서에 맞게 severity & status → Dooray 색상
func severityToColor(severity, status string) string {
	if strings.EqualFold(status, "resolved") {
		return "green"
	}
	switch strings.ToUpper(severity) {
	case "CRITICAL", "ERROR":
		return "red"
	case "WARNING", "WARN":
		return "orange"
	case "INFO", "OK":
		return "green"
	default:
		return "green"
	}
}
