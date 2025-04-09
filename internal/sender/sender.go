package sender

import (
	"git.toastoven.net/ToastCloudOps/alertmanager-webhook-proxy/pkg/models"
)

type Sender interface {
	Name() string
	Send(data *models.AlertmanagerData) error
}

//func renderTemplate(alert models.Alert) (string, error) {
//	const tmpl = `요약 : {{.Annotations.summary}}
//내용 : {{.Annotations.description}}
//
//라벨 :
//{{- range $k, $v := .Labels }}
//  - {{$k}}: {{$v}}
//{{- end }}
//
//발생시간 : {{.StartsAt.Format "2006-01-02 15:04:05"}}
//{{- if eq .Status "resolved" }}
//종료시간 : {{.EndsAt.Format "2006-01-02 15:04:05"}}
//{{- end }}`
//
//	funcMap := template.FuncMap{
//		"ToUpper": strings.ToUpper,
//	}
//
//	t := template.Must(template.New("alert").Funcs(funcMap).Parse(tmpl))
//	var buf bytes.Buffer
//	err := t.Execute(&buf, alert)
//	return buf.String(), err
//}
