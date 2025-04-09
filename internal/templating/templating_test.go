// internal/templating/templating_test.go
package templating

import (
	"bytes"
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/pkg/models"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewTemplateManager_Success(t *testing.T) {
	log.Println("[TestNewTemplateManager_Success] Creating temporary directory")
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	// 가짜 템플릿 파일 생성
	tmplFile1 := filepath.Join(tmpDir, "sample1.tmpl")
	tmplContent1 := `Hello, {{.Name}}`
	if err := os.WriteFile(tmplFile1, []byte(tmplContent1), 0644); err != nil {
		t.Fatalf("failed to write sample1.tmpl: %v", err)
	}

	// TemplateManager 초기화
	tm, err := NewTemplateManager(tmpDir)
	if err != nil {
		t.Fatalf("NewTemplateManager returned error: %v", err)
	}

	tmpl, ok := tm.Get("sample1.tmpl")
	if !ok {
		t.Errorf("expected to find sample1.tmpl in templates map")
	}
	// 템플릿 렌더링 테스트
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]string{
		"Name": "Jinsoo",
	})
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	rendered := buf.String()
	log.Printf("Rendered template output: %s", rendered)

	expected := "Hello, Jinsu"
	if rendered != expected {
		t.Errorf("unexpected render result. got: %s, want: %s", rendered, expected)
	}
}

func TestNewTemplateManager_EmptyDir(t *testing.T) {
	log.Println("[TestNewTemplateManager_EmptyDir] Creating empty directory")
	tmpDir := os.TempDir()
	defer os.RemoveAll(tmpDir)

	// TemplateManager 초기화 시 에러 기대
	_, err := NewTemplateManager(tmpDir)
	if err == nil {
		t.Fatalf("expected an error when no templates exist, got nil")
	}
	t.Logf("Got expected error: %v", err)
}

func TestWardTemplate_Parsing(t *testing.T) {
	// 임시 TemplateManager 초기화 (또는 직접 template.ParseFiles 호출)
	tm, err := NewTemplateManager("/Users/nhn/Documents/NHN/alertmanager-webhook-proxy/templates")
	if err != nil {
		t.Fatalf("failed to load templates: %v", err)
	}

	tmpl, ok := tm.Get("ward.tmpl")
	if !ok {
		t.Fatalf("template not found: ward.tmpl")
	}

	// 간단히 Execute 해보기 (Mock Alert 데이터)
	var buf bytes.Buffer
	mockAlert := models.Alert{
		Labels:      map[string]string{"alertname": "TestAlert", "namespace": "test-ns"},
		Annotations: map[string]string{"summary": "TestSummary", "description": "TestDesc"},
	}
	data := struct {
		Actor         string
		Region        string
		Env           string
		CGeneratorURL string
		OccurredAt    time.Time
		models.Alert
	}{
		Actor:  "tester",
		Region: "test-region",
		Env:    "dev",
		Alert:  mockAlert,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("failed to execute ward template: %v", err)
	}

	t.Logf("Rendered template: %s", buf.String())
	// 이 때 buf.String()이 올바른 JSON 구조인지(중괄호·쉼표 위치 등) 확인 가능
}
