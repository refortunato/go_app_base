# Swagger Documentation Guide

## Visão Geral

Este projeto utiliza **swaggo/swag** para gerar documentação API automática a partir de comentários no código. A documentação é gerada via Docker, sem necessidade de instalação local.

## Acesso Rápido

Após iniciar o servidor:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **JSON**: http://localhost:8080/swagger/doc.json

## Comandos

### Gerar documentação
```bash
make swagger
```

Este comando:
1. Roda container Docker com Go Alpine
2. Instala `swag` CLI dentro do container
3. Gera documentação a partir dos comentários
4. Cria arquivos em `docs/` (docs.go, swagger.json, swagger.yaml)

### Após modificar endpoints
Sempre que adicionar/modificar comentários Swagger, regenere:
```bash
make swagger
```

## Como Documentar Controllers

### Estrutura Básica

Adicione comentários acima de cada método do controller seguindo o padrão:

```go
// NomeDoMetodo godoc
// @Summary      Resumo curto da operação
// @Description  Descrição detalhada do que o endpoint faz
// @Tags         nome-do-modulo
// @Accept       json
// @Produce      json
// @Param        nome   in   tipo   required   "Descrição"
// @Success      200    {object}   tipo.OutputDTO
// @Failure      404    {object}   errors.ProblemDetails
// @Router       /caminho/{param} [get]
func (c *Controller) NomeDoMetodo(ctx context.WebContext) {
    // implementação...
}
```

### Exemplo Real (módulo Example)

```go
// GetExample godoc
// @Summary      Get example by ID
// @Description  Retrieves a specific example entity from the database
// @Tags         examples
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Example ID (UUID format)"
// @Success      200  {object}  usecases.GetExampleOutputDTO
// @Failure      404  {object}  errors.ProblemDetails  "Example not found"
// @Failure      500  {object}  errors.ProblemDetails  "Internal server error"
// @Router       /examples/{id} [get]
func (controller *ExampleController) GetExample(c webcontext.WebContext) {
    // implementação existente...
}
```

## Tags de Documentação

### @Summary
Resumo curto (máx 120 caracteres) exibido na lista de endpoints.

### @Description
Descrição detalhada da funcionalidade. Pode ter múltiplas linhas.

### @Tags
Agrupa endpoints por módulo/contexto. Use o nome do módulo (examples, health, products).

### @Accept / @Produce
Tipos de conteúdo aceitos/retornados. Valores comuns:
- `json`
- `xml`
- `multipart/form-data`

### @Param
Define parâmetros da requisição.

**Sintaxe**: `@Param nome in tipo required "descrição"`

**Tipos de `in`**:
- `path` - Parâmetro na URL (`/users/{id}`)
- `query` - Query string (`/users?name=john`)
- `body` - Corpo da requisição
- `header` - Cabeçalho HTTP
- `formData` - Dados de formulário

**Tipos de dado**:
- `string`, `int`, `bool`, `number`
- `{object}` - Referência a struct/DTO

**Exemplos**:
```go
// Path parameter
// @Param id path string true "User ID"

// Query parameter (opcional)
// @Param page query int false "Page number" default(1)

// Body (JSON)
// @Param request body usecases.CreateUserInputDTO true "User data"

// Header
// @Param Authorization header string true "Bearer token"
```

### @Success
Define respostas de sucesso.

**Sintaxe**: `@Success código {tipo} modelo "descrição"`

```go
// @Success 200 {object} usecases.GetUserOutputDTO
// @Success 201 {object} usecases.CreateUserOutputDTO "User created"
// @Success 204 "No content"
```

### @Failure
Define respostas de erro.

```go
// @Failure 400 {object} errors.ProblemDetails "Bad request"
// @Failure 404 {object} errors.ProblemDetails "Not found"
// @Failure 500 {object} errors.ProblemDetails "Internal error"
```

### @Router
Define o caminho e método HTTP.

**Sintaxe**: `@Router caminho [metodo]`

```go
// @Router /users/{id} [get]
// @Router /users [post]
// @Router /users/{id} [put]
// @Router /users/{id} [delete]
// @Router /users/search [get]
```

## Documentando DTOs

Swagger automaticamente reconhece structs Go. Adicione tags JSON e comentários:

```go
type CreateUserInputDTO struct {
    // Nome completo do usuário
    Name string `json:"name" example:"João Silva"`
    
    // Email único do usuário
    Email string `json:"email" example:"joao@example.com"`
    
    // Idade (deve ser maior que 18)
    Age int `json:"age" example:"25"`
}

type UserOutputDTO struct {
    ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
    Name      string    `json:"name" example:"João Silva"`
    Email     string    `json:"email" example:"joao@example.com"`
    CreatedAt time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
}
```

**Tags úteis**:
- `json:"campo"` - Nome do campo no JSON
- `example:"valor"` - Exemplo exibido no Swagger UI
- `validate:"required"` - Indica campo obrigatório (validação)

## Configuração Global

As configurações globais estão no `cmd/server/main.go`:

```go
// @title           Go App Base API
// @version         1.0
// @description     Template base para aplicações Go seguindo Clean Architecture + DDD
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @schemes http https
```

**Campos importantes**:
- `@title` - Nome da API
- `@version` - Versão
- `@description` - Descrição geral
- `@host` - Host padrão (localhost:8080)
- `@BasePath` - Prefixo base das rotas (/ = sem prefixo)
- `@schemes` - Protocolos suportados

## Workflow Recomendado

### 1. Criar novo endpoint

```go
// internal/{module}/infra/web/controllers/user_controller.go

// CreateUser godoc
// @Summary      Create new user
// @Description  Creates a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body usecases.CreateUserInputDTO true "User data"
// @Success      201 {object} usecases.CreateUserOutputDTO
// @Failure      400 {object} errors.ProblemDetails "Invalid input"
// @Failure      500 {object} errors.ProblemDetails "Internal error"
// @Router       /users [post]
func (c *UserController) CreateUser(ctx context.WebContext) {
    // implementação...
}
```

### 2. Gerar documentação

```bash
make swagger
```

### 3. Testar no Swagger UI

1. Inicie o servidor: `make dev`
2. Acesse: http://localhost:8080/swagger/index.html
3. Teste o endpoint diretamente na UI

### 4. Commit

```bash
git add .
git commit -m "feat: add user creation endpoint with swagger docs"
```

**Importante**: Sempre commitar os arquivos gerados em `docs/` para que outros desenvolvedores tenham acesso sem precisar regenerar.

## Exemplos por Tipo de Endpoint

### GET com Path Parameter

```go
// GetUser godoc
// @Summary      Get user by ID
// @Description  Retrieves user details
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID (UUID)"
// @Success      200 {object} usecases.GetUserOutputDTO
// @Failure      404 {object} errors.ProblemDetails
// @Router       /users/{id} [get]
func (c *UserController) GetUser(ctx context.WebContext) { ... }
```

### GET com Query Parameters

```go
// ListUsers godoc
// @Summary      List users
// @Description  Returns paginated list of users
// @Tags         users
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        name query string false "Filter by name"
// @Success      200 {array} usecases.UserOutputDTO
// @Router       /users [get]
func (c *UserController) ListUsers(ctx context.WebContext) { ... }
```

### POST com Body

```go
// CreateUser godoc
// @Summary      Create user
// @Description  Creates a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body usecases.CreateUserInputDTO true "User data"
// @Success      201 {object} usecases.CreateUserOutputDTO
// @Failure      400 {object} errors.ProblemDetails
// @Router       /users [post]
func (c *UserController) CreateUser(ctx context.WebContext) { ... }
```

### PUT/PATCH

```go
// UpdateUser godoc
// @Summary      Update user
// @Description  Updates existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        request body usecases.UpdateUserInputDTO true "Updated data"
// @Success      200 {object} usecases.UserOutputDTO
// @Failure      404 {object} errors.ProblemDetails
// @Router       /users/{id} [put]
func (c *UserController) UpdateUser(ctx context.WebContext) { ... }
```

### DELETE

```go
// DeleteUser godoc
// @Summary      Delete user
// @Description  Removes user from system
// @Tags         users
// @Param        id path string true "User ID"
// @Success      204 "No content"
// @Failure      404 {object} errors.ProblemDetails
// @Router       /users/{id} [delete]
func (c *UserController) DeleteUser(ctx context.WebContext) { ... }
```

## Troubleshooting

### Erro: "no Go files in /app"
**Causa**: Warning do swag, pode ser ignorado. A documentação é gerada corretamente.

### Endpoint não aparece no Swagger
**Soluções**:
1. Verifique se os comentários estão corretos (sem espaços extras)
2. Rode `make swagger` novamente
3. Reinicie o servidor (`make down && make dev`)

### Tipo não reconhecido
**Soluções**:
1. Certifique-se que o tipo está exportado (primeira letra maiúscula)
2. Use caminho completo: `usecases.GetExampleOutputDTO`
3. Verifique se o arquivo está no mesmo módulo Go

### Swagger UI mostra erro 404
**Verificações**:
1. Servidor está rodando? (`make dev`)
2. Rota `/swagger/*any` está registrada? (verificar `internal/infra/web/register_routes.go`)
3. Import `_ "github.com/refortunato/go_app_base/docs"` existe no `main.go`?

## Arquivos Gerados

Após `make swagger`, os seguintes arquivos são criados em `docs/`:

- **docs.go** - Código Go com documentação embarcada
- **swagger.json** - Especificação OpenAPI 3.0 em JSON
- **swagger.yaml** - Especificação OpenAPI 3.0 em YAML

**Importante**: Sempre commite esses arquivos no git.

## Boas Práticas

### ✅ Faça
- Documente TODOS os endpoints
- Use exemplos nos DTOs (`example:"valor"`)
- Agrupe endpoints por Tags (módulo)
- Descreva possíveis erros com @Failure
- Especifique campos obrigatórios (`true/false`)
- Regenere docs após cada mudança

### ❌ Evite
- Comentários genéricos ("Gets data")
- Esquecer de rodar `make swagger`
- Usar tipos não exportados
- Misturar tags de diferentes módulos
- Documentação desatualizada no git

## Integração com Clean Architecture

### Controllers (Infra Layer)
**Onde documentar**: Nos controllers (`internal/{module}/infra/web/controllers/`)

```go
// internal/example/infra/web/controllers/example_controller.go
package controllers

// GetExample godoc
// @Summary Get example
// @Tags examples
// @Router /examples/{id} [get]
func (c *ExampleController) GetExample(ctx context.WebContext) { ... }
```

### Use Cases (Application Layer)
**DTOs devem ser documentados** com tags JSON:

```go
// internal/example/core/application/usecases/get_example.go
package usecases

type GetExampleOutputDTO struct {
    ID          string    `json:"id" example:"uuid-here"`
    Description string    `json:"description" example:"Sample description"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### Domain Entities
**Não documentar diretamente**. Use DTOs da camada de aplicação para expor na API.

## Checklist: Novo Endpoint

- [ ] Adicionar comentários Swagger no controller
- [ ] Documentar DTO de entrada (se houver)
- [ ] Documentar DTO de saída
- [ ] Especificar todos os @Param necessários
- [ ] Adicionar @Success e @Failure apropriados
- [ ] Rodar `make swagger`
- [ ] Testar no Swagger UI
- [ ] Commitar arquivos `docs/`

## Referências

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [Gin Swagger Integration](https://github.com/swaggo/gin-swagger)
