package sender

import (
	"bytes"
	"errors"
	"fmt"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/config"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/internal/templating"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/pkg/models"
	"net/http"
	"strings"
	"time"
)

const (
	WardTemplateName = "ward.tmpl"
	WardName         = "ward"
)

type WardSender struct {
	Config *config.WardConfig
	Tm     *templating.TemplateManager
}

func (s *WardSender) Name() string {
	return WardName
}

func (s *WardSender) Send(data *models.AlertmanagerData) error {
	if len(data.Alerts) == 0 {
		return errors.New("no alerts found")
	}
	for i, alert := range data.Alerts {
		// 1) 먼저 Alert 필드를 정규화
		s.normalizeAlertFields(&alert)

		// 2) 템플릿 렌더링
		tmplData := struct {
			Actor        string
			Region       string
			Env          string
			GeneratorURL string
			OccurredAt   time.Time
			models.Alert
		}{
			Actor:        s.Config.Actor,
			Region:       s.Config.Region,
			Env:          s.Config.Env,
			GeneratorURL: alert.GeneratorURL, // or any rename if needed
			OccurredAt:   alert.StartsAt,     // or choose different logic
			//OccurredAt: time.Now(), // 테스트
			Alert: alert,
		}

		jsonBody, err := s.renderTemplate(tmplData)
		if err != nil {
			return fmt.Errorf("failed to render template for alert #%d: %w", i, err)
		}

		//fmt.Println("Ward JSON Body:", jsonBody) // Debugging line

		// (3) HTTP POST
		if err := s.postToWard(jsonBody); err != nil {
			//fmt.Println("Ward POST Error:", err) // Debugging line
			return fmt.Errorf("failed to post alert #%d to ward: %w", i, err)
		}
	}
	return nil
}

// renderTemplate: ward.tmpl을 실행해 최종 JSON 문자열 생성
func (s *WardSender) renderTemplate(data any) (string, error) {
	tmpl, ok := s.Tm.Get(WardTemplateName)
	if !ok {
		return "", fmt.Errorf("template not found: %s", WardTemplateName)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// postToWard: 완성된 JSON 문자열을 POST
func (s *WardSender) postToWard(body string) error {
	req, err := http.NewRequest("POST", s.Config.EventURL, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", s.Config.Token) // if needed

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Ward returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *WardSender) normalizeAlertFields(alert *models.Alert) {
	// (1) severity 매핑
	var mappedSeverity string
	switch strings.ToLower(alert.Labels["severity"]) {
	case "critical":
		mappedSeverity = "CRITICAL"
	case "error":
		mappedSeverity = "ERROR"
	case "warning", "warn":
		mappedSeverity = "WARNING"
	case "info", "informational":
		mappedSeverity = "INFO"
	case "ok", "resolved":
		mappedSeverity = "OK"
	case "", "none", "unknown":
		mappedSeverity = "NONE"
	default:
		mappedSeverity = "NONE"
	}
	// 실제 라벨에 바로 반영
	alert.Labels["severity"] = mappedSeverity

	// (2) 필요하다면 다른 필드도 수정 가능
	// e.g. alert.Labels["namespace"] = strings.ToLower(alert.Labels["namespace"])

	// (3) 그 외 라벨/어노테이션, 등등도 수정 가능
}
