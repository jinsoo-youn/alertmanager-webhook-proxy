package templating

import (
	"fmt"
	"log"
	"path/filepath"
	"text/template"
)

// TemplateManager는 템플릿을 로드/캐싱하여 필요 시 반환하는 구조체
type TemplateManager struct {
	templates map[string]*template.Template
}

// NewTemplateManager: 주어진 디렉터리의 템플릿을 로딩해 TemplateManager를 생성
func NewTemplateManager(templatesDir string) (*TemplateManager, error) {
	log.Printf("[templating] NewTemplateManager called. Directory: %s", templatesDir)

	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	if err := tm.load(templatesDir); err != nil {
		return nil, err
	}
	return tm, nil
}

// load: 내부적으로 *.tmpl 파일을 파싱해 templates 맵에 저장
func (tm *TemplateManager) load(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.tmpl"))
	if err != nil {
		return fmt.Errorf("failed to glob templates: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no .tmpl files found in %s", dir)
	}

	log.Printf("[templating] Found %d template files in %s", len(files), dir)

	for _, f := range files {
		tmplName := filepath.Base(f)
		log.Printf("[templating] Parsing template: %s", tmplName)
		t, err := template.ParseFiles(f)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", f, err)
		}
		tm.templates[tmplName] = t
	}

	log.Printf("[templating] Successfully loaded all templates.")
	return nil
}

// Get: 템플릿 이름으로 템플릿을 얻어옴 (맵 조회 스타일)
func (tm *TemplateManager) Get(name string) (*template.Template, bool) {
	t, ok := tm.templates[name]
	return t, ok
}
