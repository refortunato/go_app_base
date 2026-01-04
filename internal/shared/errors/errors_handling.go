package errors

import (
	"encoding/json"
)

// ProblemDetails segue RFC7807 e inclui campos extras.
type ProblemDetails struct {
	Type         string `json:"type,omitempty"`     // URI identificando o tipo do erro
	Title        string `json:"title"`              // Título curto do erro
	Status       int    `json:"status"`             // Código HTTP
	Detail       string `json:"detail,omitempty"`   // Descrição detalhada
	Instance     string `json:"instance,omitempty"` // URI da ocorrência do erro
	Code         string `json:"code"`               // Código específico do erro
	ErrorContext string `json:"error_context"`      // business ou infra
}

// Função para criar um novo erro RFC7807
func NewProblemDetails(status int, title, detail, code, errorContext string) *ProblemDetails {
	return &ProblemDetails{
		Type:         "about:blank",
		Title:        title,
		Status:       status,
		Detail:       detail,
		Code:         code,
		ErrorContext: errorContext,
	}
}

// Implementa a interface error
func (pd *ProblemDetails) Error() string {
	b, err := json.Marshal(pd)
	if err != nil {
		return pd.Title + " (" + pd.Code + ")"
	}
	return string(b)
}
