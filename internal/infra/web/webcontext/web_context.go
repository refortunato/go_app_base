package webcontext

import "context"

type WebContext interface {
	JSON(code int, obj any)
	BindJSON(obj any) error
	Param(key string) string
	Query(key string) string
	GetHeader(key string) string
	SetHeader(key, value string)
	GetContext() context.Context
}
