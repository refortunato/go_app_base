# Geração Automática de Documentação Swagger

## Visão Geral

Os scripts de criação de entidades (`create-entity.sh`) foram aprimorados para **gerar automaticamente documentação Swagger completa** nos controllers, eliminando a necessidade de digitação manual.

## Como Funciona

Quando você cria uma nova entidade usando:

```bash
./scripts/create-entity.sh
```

O controller gerado já virá com **todas as annotations Swagger** incluídas automaticamente:

### Exemplo de Controller Gerado

```go
// Create godoc
// @Summary      Create new product
// @Description  Creates a new product in the system
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        input  body      usecases.CreateProductInputDTO  true  "Product data"
// @Success      201    {object}  usecases.CreateProductOutputDTO
// @Failure      400    {object}  errors.ProblemDetails  "Invalid request data"
// @Failure      500    {object}  errors.ProblemDetails  "Internal server error"
// @Router       /products [post]
func (c *ProductController) Create(ctx context.WebContext) {
    // implementação...
}
```

## Funcionalidades Implementadas

### ✅ Annotations Completas

Cada método do controller recebe automaticamente:

- **@Summary** - Resumo da operação
- **@Description** - Descrição detalhada
- **@Tags** - Tag do módulo (agrupa endpoints)
- **@Accept / @Produce** - Tipos de conteúdo (JSON)
- **@Param** - Parâmetros (path, query, body)
- **@Success** - Resposta de sucesso com DTO
- **@Failure** - Respostas de erro (400, 404, 500)
- **@Router** - Rota e método HTTP

### ✅ Suporte para Ambas Arquiteturas

**DDD (Clean Architecture)**:
- Usa DTOs dos use cases: `usecases.CreateProductInputDTO`
- 5 operações: Create, Get, List, Update, Delete

**4-tier (Simplified)**:
- Usa models: `models.Product`
- 5 operações: Create, Get, List, Update, Delete

### ✅ Nomes Dinâmicos

As annotations se adaptam automaticamente:
- `${ENTITY_NAME_LOWER}` → `product`, `user`, `order`
- `${ENTITY_NAME_CAPITALIZED}` → `Product`, `User`, `Order`
- `${MODULE_NAME}` → `products`, `users`, `orders`

## Workflow Completo

### 1. Criar Nova Entidade

```bash
./scripts/create-entity.sh
```

### 2. Script Gera Controller com Swagger

O controller já virá documentado:

```go
package controllers

// Get godoc
// @Summary      Get product by ID
// @Description  Retrieves a specific product from the database
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Product ID (UUID format)"
// @Success      200  {object}  usecases.GetProductOutputDTO
// @Failure      404  {object}  errors.ProblemDetails  "Product not found"
// @Failure      500  {object}  errors.ProblemDetails  "Internal server error"
// @Router       /products/{id} [get]
func (c *ProductController) Get(ctx context.WebContext) {
    // ...
}
```

### 3. Gerar Documentação Swagger

```bash
make swagger
```

### 4. Acessar Swagger UI

```
http://localhost:8080/swagger/index.html
```

## Vantagens

### 🚀 **Zero Trabalho Manual**
- Não precisa digitar annotations
- Não precisa lembrar sintaxe
- Não precisa copiar/colar de outros controllers

### 📋 **Documentação Consistente**
- Todas as entidades seguem o mesmo padrão
- Mesma estrutura em todos os endpoints
- Menos erros e esquecimentos

### ⚡ **Produtividade**
- Cria entidade + documentação em 1 comando
- Foco na lógica de negócio, não em documentação
- Ideal para prototipagem rápida

### 🔄 **Atualização Fácil**
- Se padrão mudar, basta atualizar o script
- Todas as novas entidades seguirão o novo padrão
- Manutenção centralizada

## Personalização

### Modificar Templates

Caso queira alterar o padrão de documentação, edite o script:

**Arquivo**: `scripts/create-entity.sh`

**Localização**: Busque por `# Create godoc` ou `@Summary`

**Exemplo de customização**:

```bash
# Adicionar mais campos na descrição
# @Description  Creates a new ${ENTITY_NAME_LOWER} in the system with validation
```

### Adicionar Annotations Customizadas

Você pode adicionar annotations específicas para determinados tipos:

```bash
# Se campo for email, adicionar validação na doc
if [ "$field_name" = "email" ]; then
    echo "// @Param        email  body  string  true  \"Valid email address\" format(email)" >> "$FILE"
fi
```

## Limitações Conhecidas

### ⚠️ **Annotations Genéricas**
- Descrições são padrão (não customizadas por entidade)
- Exemplos não são incluídos automaticamente nos DTOs
- Validações específicas precisam ser adicionadas manualmente

**Solução**: Após gerar, edite o controller para adicionar detalhes específicos.

### ⚠️ **DTOs Sem Examples**
- Os DTOs gerados não incluem tag `example`
- Swagger UI não mostrará valores de exemplo automaticamente

**Solução futura**: Adicionar geração de examples baseado em tipos:
- `string` → `"Example Text"`
- `int` → `123`
- `float64` → `99.99`
- `time.Time` → `"2026-02-01T10:00:00Z"`

## Roadmap (Melhorias Futuras)

### 🔮 Fase 2: Examples Automáticos

Gerar examples em DTOs:

```go
type CreateProductInputDTO struct {
    Name  string  `json:"name" example:"iPhone 15"`
    Price float64 `json:"price" example:"999.99"`
}
```

### 🔮 Fase 3: Validações na Documentação

Adicionar constraints de validação:

```go
// @Param        price  body  float64  true  "Price in USD" minimum(0.01) maximum(999999.99)
```

### 🔮 Fase 4: Documentação de Erros Personalizados

Documentar erros específicos do domínio:

```go
// @Failure      422  {object}  errors.ProblemDetails  "Invalid email format"
```

## Conclusão

Com essa implementação, você tem **documentação Swagger automática** sem esforço adicional. Basta criar a entidade e a documentação está pronta!

**Fluxo ideal**:
1. `./scripts/create-entity.sh` → Cria entidade + controller + docs
2. `make swagger` → Gera Swagger
3. Implementar lógica de negócio
4. **Documentação já está pronta!** ✅
