# Create Entity Script Guide

The `create-entity.sh` script automates the creation of complete CRUD entities within existing modules, supporting both **DDD (Clean Architecture)** and **4-tier (simplified)** architectures.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Usage](#usage)
- [Entity Creation Process](#entity-creation-process)
- [What Gets Created](#what-gets-created)
- [Step-by-Step Example](#step-by-step-example)
- [Field Types](#field-types)
- [Generated Code](#generated-code)
- [Best Practices](#best-practices)

## Overview

This script helps you:
- Create complete CRUD operations for an entity
- Generate domain entities with proper encapsulation (DDD) or models (4-tier)
- Create repository interfaces and implementations
- Generate use cases (DDD) or services (4-tier)
- Create controllers with proper error handling
- Automatically wire dependencies
- Register routes
- Generate SQL table creation script

## Prerequisites

- Module already created (use `create-module.sh` first)
- Go 1.25+ installed
- Project initialized with `go.mod`
- Write permissions in the project directory

## Usage

From the project root, run:

```bash
./scripts/create-entity.sh
```

The script will guide you through:

1. **Select Module**: Choose from existing modules
2. **Architecture Detection**: Automatically detects if module is DDD or 4-tier
3. **Entity Name**: Enter entity name (e.g., `User`, `Order`, `Product`)
4. **Define Fields**: Add fields with their MySQL types
5. **Automatic Generation**: Creates all necessary files

## Entity Creation Process

### Step 1: Select Module

```
ℹ Available modules:
  - example
  - health
  - payment
  - user

ℹ Enter the module name where the entity will be created:
payment
```

### Step 2: Architecture Detection

```
ℹ Detected architecture: DDD (Clean Architecture)
```

or

```
ℹ Detected architecture: 4-tier (simplified)
```

### Step 3: Entity Name

```
ℹ Enter the entity name (e.g., User, Order, Payment):
CreditCard
```

The script converts to:
- **Capitalized**: `CreditCard` (for types)
- **Lowercase**: `creditcard` (for files and tables)

### Step 4: Define Fields

```
ℹ Enter the fields for the entity (format: field_name:TYPE)
ℹ Examples: name:VARCHAR(255), age:INT, price:DECIMAL(10,2), created_at:TIMESTAMP
ℹ Type 'done' when finished

Field (or 'done'): card_number:VARCHAR(16)
✓ Added field: CardNumber (string) -> DB: card_number

Field (or 'done'): holder_name:VARCHAR(255)
✓ Added field: HolderName (string) -> DB: holder_name

Field (or 'done'): expiry_date:DATE
✓ Added field: ExpiryDate (time.Time) -> DB: expiry_date

Field (or 'done'): cvv:VARCHAR(4)
✓ Added field: Cvv (string) -> DB: cvv

Field (or 'done'): is_active:BOOLEAN
✓ Added field: IsActive (bool) -> DB: is_active

Field (or 'done'): done
```

### Step 5: Automatic Generation

The script generates all files and provides the SQL to create the table.

## What Gets Created

### For DDD Architecture:

1. **Domain Entity** (`core/domain/entities/{entity}.go`)
   - Private fields with public getters/setters
   - Constructor `New{Entity}()` with validation
   - Restore function for DB reconstitution
   - Validation method

2. **Repository Interface** (`core/application/repositories/{entity}_repository.go`)
   - Save, FindById, FindAll, Count, Update, Delete methods

3. **Repository Implementation** (`infra/repositories/{entity}_mysql_repository.go`)
   - MySQL implementation of the interface
   - DB model struct
   - Mapping functions (mapToDomain)

4. **Use Cases** (`core/application/usecases/`)
   - `create_{entity}.go` - Create entity
   - `get_{entity}.go` - Get by ID
   - `list_{entity}.go` - List with pagination
   - `update_{entity}.go` - Update entity
   - `delete_{entity}.go` - Delete entity

5. **Controller** (`infra/web/controllers/{entity}_controller.go`)
   - HTTP handlers for all operations
   - Proper error handling

6. **Updates**
   - `infra/module.go` - Wires all dependencies
   - `infra/web/routes.go` - Registers all routes

### For 4-Tier Architecture:

1. **Model** (`models/{entity}.go`)
   - Public fields
   - Basic structure

2. **Repository** (`repositories/{entity}_repository.go`)
   - Database access methods

3. **Service** (`services/{entity}_service.go`)
   - Business logic and validation

4. **Controller** (`controllers/{entity}_controller.go`)
   - HTTP handlers

5. **Updates**
   - `module.go` - Wires all dependencies
   - `routes.go` - Registers all routes

## Field Types

### Supported MySQL Types:

| MySQL Type | Go Type | Examples |
|------------|---------|----------|
| VARCHAR, CHAR, TEXT | string | `name:VARCHAR(255)` |
| INT, BIGINT, SMALLINT | int | `age:INT` |
| FLOAT, DOUBLE, DECIMAL | float64 | `price:DECIMAL(10,2)` |
| BOOL, BOOLEAN | bool | `is_active:BOOLEAN` |
| DATE, DATETIME, TIMESTAMP | time.Time | `created_at:TIMESTAMP` |

### Field Naming Conventions:

Input can be any format. The script converts to:
- **snake_case** for database columns
- **camelCase** for private Go fields (DDD)
- **PascalCase** for public methods and struct fields

Examples:
```
Input: card_number  → DB: card_number, Private: cardNumber, Public: CardNumber
Input: CardNumber   → DB: card_number, Private: cardNumber, Public: CardNumber
Input: card-number  → DB: card_number, Private: cardNumber, Public: CardNumber
```

## Step-by-Step Example

### Example 1: Creating a Product Entity (4-tier)

```bash
$ ./scripts/create-entity.sh

ℹ Available modules:
  - user

ℹ Enter the module name where the entity will be created:
user

ℹ Detected architecture: 4-tier (simplified)

ℹ Enter the entity name (e.g., User, Order, Payment):
Product

ℹ Entity name: Product (lowercase: product)

ℹ Enter the fields for the entity (format: field_name:TYPE)
Field (or 'done'): name:VARCHAR(255)
✓ Added field: Name (string) -> DB: name

Field (or 'done'): description:TEXT
✓ Added field: Description (string) -> DB: description

Field (or 'done'): price:DECIMAL(10,2)
✓ Added field: Price (float64) -> DB: price

Field (or 'done'): stock:INT
✓ Added field: Stock (int) -> DB: stock

Field (or 'done'): done

ℹ Creating entity 'Product' with 4 field(s)...

ℹ Generating 4-tier structure...
✓ Created model: internal/user/models/product.go
✓ Created repository: internal/user/repositories/product_repository.go
✓ Created service: internal/user/services/product_service.go
✓ Created controller: internal/user/controllers/product_controller.go
✓ Updated module.go
✓ Updated routes.go

✓ Entity 'Product' created successfully!

ℹ SQL to create the table:

CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    price DECIMAL(10,2),
    stock INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

ℹ API Endpoints created:
  POST   /products          - Create new product
  GET    /products/:id      - Get product by ID
  GET    /products          - List all products (with pagination)
  PUT    /products/:id      - Update product
  DELETE /products/:id      - Delete product
```

### Example 2: Creating a Boleto Entity (DDD)

```bash
$ ./scripts/create-entity.sh

ℹ Available modules:
  - payment

ℹ Enter the module name where the entity will be created:
payment

ℹ Detected architecture: DDD (Clean Architecture)

ℹ Enter the entity name (e.g., User, Order, Payment):
Boleto

Field (or 'done'): description:VARCHAR(255)
✓ Added field: Description (string) -> DB: description

Field (or 'done'): destination_id:VARCHAR(36)
✓ Added field: DestinationId (string) -> DB: destination_id

Field (or 'done'): expiry_date:DATE
✓ Added field: ExpiryDate (time.Time) -> DB: expiry_date

Field (or 'done'): status:VARCHAR(20)
✓ Added field: Status (string) -> DB: status

Field (or 'done'): value:DECIMAL(10,2)
✓ Added field: Value (float64) -> DB: value

Field (or 'done'): done

ℹ Generating DDD structure...
✓ Created entity: internal/payment/core/domain/entities/boleto.go
✓ Created repository interface: internal/payment/core/application/repositories/boleto_repository.go
✓ Created repository implementation: internal/payment/infra/repositories/boleto_mysql_repository.go
✓ Created use case: internal/payment/core/application/usecases/create_boleto.go
✓ Created use case: internal/payment/core/application/usecases/get_boleto.go
✓ Created use case: internal/payment/core/application/usecases/list_boleto.go
✓ Created use case: internal/payment/core/application/usecases/update_boleto.go
✓ Created use case: internal/payment/core/application/usecases/delete_boleto.go
✓ Created controller: internal/payment/infra/web/controllers/boleto_controller.go
✓ Updated module.go
✓ Updated routes.go

✓ Entity 'Boleto' created successfully!
```

## Generated Code

### DDD Architecture Example

**Domain Entity:**
```go
type Boleto struct {
    id            string
    description   string
    destinationId string
    expiryDate    time.Time
    status        string
    value         float64
    createdAt     time.Time
    updatedAt     time.Time
}

func NewBoleto(
    description string,
    destinationId string,
    expiryDate time.Time,
    status string,
    value float64,
) (*Boleto, error) {
    entity := &Boleto{
        id:            shared.GenerateId(),
        description:   description,
        destinationId: destinationId,
        expiryDate:    expiryDate,
        status:        status,
        value:         value,
        createdAt:     time.Now().UTC(),
        updatedAt:     time.Now().UTC(),
    }
    if err := entity.Validate(); err != nil {
        return nil, err
    }
    return entity, nil
}

func (e *Boleto) GetDescription() string {
    return e.description
}

func (e *Boleto) SetDescription(description string) {
    e.description = description
    e.updatedAt = time.Now().UTC()
}
```

**Use Case:**
```go
type CreateBoletoInputDTO struct {
    Description   string
    DestinationId string
    ExpiryDate    time.Time
    Status        string
    Value         float64
}

type CreateBoletoOutputDTO struct {
    Id string `json:"id"`
}

func (uc *CreateBoletoUseCase) Execute(input CreateBoletoInputDTO) (*CreateBoletoOutputDTO, error) {
    entity, err := entities.NewBoleto(
        input.Description,
        input.DestinationId,
        input.ExpiryDate,
        input.Status,
        input.Value,
    )
    if err != nil {
        return nil, err
    }

    if err := uc.repository.Save(entity); err != nil {
        return nil, err
    }

    return &CreateBoletoOutputDTO{Id: entity.GetId()}, nil
}
```

### 4-Tier Architecture Example

**Model:**
```go
type Product struct {
    ID          string
    Name        string
    Description string
    Price       float64
    Stock       int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Service:**
```go
func (s *ProductService) CreateProduct(
    name string,
    description string,
    price float64,
    stock int,
) (*models.Product, error) {
    // TODO: Add validation logic here
    
    now := time.Now().UTC()
    entity := &models.Product{
        ID:          shared.GenerateId(),
        Name:        name,
        Description: description,
        Price:       price,
        Stock:       stock,
        CreatedAt:   now,
        UpdatedAt:   now,
    }

    if err := s.repository.Save(entity); err != nil {
        return nil, fmt.Errorf("failed to create product: %w", err)
    }

    return entity, nil
}
```

## Best Practices

### 1. Choose Meaningful Names
```bash
# Good
Entity: User, CreditCard, OrderItem
Fields: first_name, email_address, is_active

# Avoid
Entity: Data, Info, Record
Fields: field1, val, tmp
```

### 2. Use Appropriate Types
```bash
# Strings
email:VARCHAR(255)
description:TEXT

# Numbers
age:INT
price:DECIMAL(10,2)

# Dates
birth_date:DATE
created_at:TIMESTAMP

# Booleans
is_active:BOOLEAN
has_permissions:BOOLEAN
```

### 3. Validation
Add business validation in:
- **DDD**: `entity.Validate()` method
- **4-tier**: Service methods

### 4. Error Handling
The generated code includes:
- Repository errors
- Validation errors
- Not found errors (with `ReturnNotFoundError`)

### 5. Pagination
All List operations include pagination:
```bash
GET /products?page=1&limit=10
```

Response:
```json
{
  "items": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 100,
    "total_pages": 10
  }
}
```

## API Endpoints

All entities get full CRUD endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/{entities}` | Create new entity |
| GET | `/{entities}/:id` | Get entity by ID |
| GET | `/{entities}` | List all entities (paginated) |
| PUT | `/{entities}/:id` | Update entity |
| DELETE | `/{entities}/:id` | Delete entity |

Example requests:

```bash
# Create
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-end laptop",
    "price": 1299.99,
    "stock": 10
  }'

# Get
curl http://localhost:8080/products/123e4567-e89b-12d3-a456-426614174000

# List (with pagination)
curl http://localhost:8080/products?page=1&limit=20

# Update
curl -X PUT http://localhost:8080/products/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Laptop",
    "description": "Updated description",
    "price": 1199.99,
    "stock": 15
  }'

# Delete
curl -X DELETE http://localhost:8080/products/123e4567-e89b-12d3-a456-426614174000
```

## Database Table

After running the script, execute the provided SQL:

```sql
CREATE TABLE creditcards (
    id VARCHAR(36) PRIMARY KEY,
    card_number VARCHAR(16),
    holder_name VARCHAR(255),
    expiry_date DATE,
    cvv VARCHAR(4),
    is_active BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Troubleshooting

### Module Not Found
```
✗ Module 'payment' does not exist
ℹ Please create the module first using: ./scripts/create-module.sh
```
**Solution**: Create the module first with `create-module.sh`

### Architecture Not Detected
```
✗ Could not detect module architecture
```
**Solution**: Module must have either `core/domain` (DDD) or `models` (4-tier) directory

### Invalid Field Format
```
Field must be in format: field_name:TYPE
```
**Solution**: Use format like `name:VARCHAR(255)` or `age:INT`

### Compilation Errors
After generation, if you get errors:
1. Run `go mod tidy`
2. Check imports in generated files
3. Verify field names don't conflict with Go keywords

## See Also

- [Create Module Guide](./create-module-guide.md) - Create new modules first
- [Architecture Decision](../implementation/routes-management.md) - Understanding the routing approach
- [Dependency Management](../implementation/dependency-management.md) - How DI works in this project
