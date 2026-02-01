#!/bin/bash

# Script to create a new module in the Go application
# Supports both 4-tier (simplified) and DDD architectures

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored messages
print_info() {
    echo -e "${BLUE}ℹ ${NC}$1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Get the project root directory (where the script is located)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Get the module path from go.mod
MODULE_PATH=$(grep -E "^module " "$PROJECT_ROOT/go.mod" | awk '{print $2}')

if [ -z "$MODULE_PATH" ]; then
    print_error "Could not find module path in go.mod"
    exit 1
fi

print_info "Project module path: $MODULE_PATH"
echo ""

# Ask for module name
print_info "Enter the module name (e.g., user, order, payment):"
read -r MODULE_NAME

if [ -z "$MODULE_NAME" ]; then
    print_error "Module name cannot be empty"
    exit 1
fi

# Convert module name to lowercase and replace spaces with underscores
MODULE_NAME=$(echo "$MODULE_NAME" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')

print_info "Module name: $MODULE_NAME"
echo ""

# Ask for architecture type
print_info "Select architecture type:"
echo "1) 4-tier (simplified) - For simple CRUD operations"
echo "2) DDD (Clean Architecture) - For complex business logic"
read -p "Enter your choice (1 or 2): " ARCH_TYPE

if [ "$ARCH_TYPE" != "1" ] && [ "$ARCH_TYPE" != "2" ]; then
    print_error "Invalid choice. Please enter 1 or 2"
    exit 1
fi

MODULE_DIR="$PROJECT_ROOT/internal/$MODULE_NAME"

# Check if module already exists
if [ -d "$MODULE_DIR" ]; then
    print_error "Module '$MODULE_NAME' already exists at $MODULE_DIR"
    exit 1
fi

print_info "Creating module '$MODULE_NAME'..."
echo ""

# Create module structure based on architecture type
if [ "$ARCH_TYPE" = "1" ]; then
    print_info "Creating 4-tier (simplified) architecture..."
    
    # Create directories
    mkdir -p "$MODULE_DIR/models"
    mkdir -p "$MODULE_DIR/repositories"
    mkdir -p "$MODULE_DIR/services"
    mkdir -p "$MODULE_DIR/controllers"
    
    print_success "Created directory structure"
    
    # Create module.go
    cat > "$MODULE_DIR/module.go" <<EOF
package ${MODULE_NAME}

import (
	"database/sql"

	"${MODULE_PATH}/internal/${MODULE_NAME}/controllers"
	"${MODULE_PATH}/internal/${MODULE_NAME}/repositories"
	"${MODULE_PATH}/internal/${MODULE_NAME}/services"
)

// ${MODULE_NAME^}Module holds all initialized dependencies for the ${MODULE_NAME} module (4-tier architecture)
type ${MODULE_NAME^}Module struct {
	// TODO: Add your controllers here
	// Example: ProductController *controllers.ProductController
}

// New${MODULE_NAME^}Module creates and wires all dependencies for the ${MODULE_NAME} module
func New${MODULE_NAME^}Module(db *sql.DB) *${MODULE_NAME^}Module {
	// TODO: Wire your dependencies here
	// Step 1: Initialize repositories
	// Step 2: Initialize services (inject repositories)
	// Step 3: Initialize controllers (inject services)
	// Step 4: Return module with all dependencies wired
	
	return &${MODULE_NAME^}Module{
		// TODO: Initialize your dependencies
	}
}
EOF
    
    print_success "Created module.go"
    
    # Create routes.go
    cat > "$MODULE_DIR/routes.go" <<EOF
package ${MODULE_NAME}

import (
	"github.com/gin-gonic/gin"
	"${MODULE_PATH}/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the ${MODULE_NAME} module (4-tier architecture)
func RegisterRoutes(router *gin.Engine, module *${MODULE_NAME^}Module) {
	// TODO: Add your routes here
	// Example:
	// router.GET("/${MODULE_NAME}/:id", func(ctx *gin.Context) {
	//     module.YourController.GetMethod(context.NewGinContextAdapter(ctx))
	// })
}
EOF
    
    print_success "Created routes.go"
    
else
    print_info "Creating DDD (Clean Architecture) structure..."
    
    # Create directories
    mkdir -p "$MODULE_DIR/core/application/repositories"
    mkdir -p "$MODULE_DIR/core/application/usecases"
    mkdir -p "$MODULE_DIR/core/domain/entities"
    mkdir -p "$MODULE_DIR/infra/repositories"
    mkdir -p "$MODULE_DIR/infra/web/controllers"
    
    print_success "Created directory structure"
    
    # Create module.go
    cat > "$MODULE_DIR/infra/module.go" <<EOF
package infra

import (
	"database/sql"

	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/usecases"
	"${MODULE_PATH}/internal/${MODULE_NAME}/infra/repositories"
	"${MODULE_PATH}/internal/${MODULE_NAME}/infra/web/controllers"
)

// ${MODULE_NAME^}Module encapsulates all dependencies for the ${MODULE_NAME} module
type ${MODULE_NAME^}Module struct {
	// TODO: Add your controllers and use cases here
	// Example: GetExampleUseCase *usecases.GetExampleUseCase
	// Example: ExampleController *controllers.ExampleController
}

// New${MODULE_NAME^}Module creates and wires all dependencies for the ${MODULE_NAME} module
func New${MODULE_NAME^}Module(db *sql.DB) *${MODULE_NAME^}Module {
	// TODO: Wire your dependencies here
	// Step 1: Initialize repositories
	// Step 2: Initialize use cases (inject repositories)
	// Step 3: Initialize controllers (inject use cases)
	// Step 4: Return module with all dependencies wired
	
	return &${MODULE_NAME^}Module{
		// TODO: Initialize your dependencies
	}
}
EOF
    
    print_success "Created infra/module.go"
    
    # Create routes.go
    cat > "$MODULE_DIR/infra/web/routes.go" <<EOF
package web

import (
	"github.com/gin-gonic/gin"
	"${MODULE_PATH}/internal/${MODULE_NAME}/infra"
	"${MODULE_PATH}/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the ${MODULE_NAME} module
func RegisterRoutes(router *gin.Engine, module *infra.${MODULE_NAME^}Module) {
	// TODO: Add your routes here
	// Example:
	// router.GET("/${MODULE_NAME}/:id", func(ctx *gin.Context) {
	//     module.YourController.GetMethod(context.NewGinContextAdapter(ctx))
	// })
}
EOF
    
    print_success "Created infra/web/routes.go"
fi

echo ""
print_info "Updating container.go..."

# Update container.go
CONTAINER_FILE="$PROJECT_ROOT/cmd/server/container/container.go"

# Determine the import path and struct name based on architecture
if [ "$ARCH_TYPE" = "1" ]; then
    IMPORT_ALIAS="${MODULE_NAME}"
    IMPORT_PATH="${MODULE_PATH}/internal/${MODULE_NAME}"
    STRUCT_NAME="${MODULE_NAME^}Module"
    STRUCT_TYPE="*${MODULE_NAME}.${MODULE_NAME^}Module"
    INIT_CALL="${MODULE_NAME}.New${MODULE_NAME^}Module(db)"
else
    IMPORT_ALIAS="${MODULE_NAME}Infra"
    IMPORT_PATH="${MODULE_PATH}/internal/${MODULE_NAME}/infra"
    STRUCT_NAME="${MODULE_NAME^}Module"
    STRUCT_TYPE="*${MODULE_NAME}Infra.${MODULE_NAME^}Module"
    INIT_CALL="${MODULE_NAME}Infra.New${MODULE_NAME^}Module(db)"
fi

# Add import
if ! grep -q "\"$IMPORT_PATH\"" "$CONTAINER_FILE"; then
    # Find the last import line and add the new import after it
    sed -i.bak "/^import (/a\\
	${IMPORT_ALIAS} \"${IMPORT_PATH}\"
" "$CONTAINER_FILE"
    print_success "Added import to container.go"
else
    print_warning "Import already exists in container.go"
fi

# Add field to Container struct
if ! grep -q "${STRUCT_NAME}" "$CONTAINER_FILE"; then
    # Find the Container struct and add the new field
    sed -i.bak "/type Container struct {/a\\
	${STRUCT_NAME} ${STRUCT_TYPE}
" "$CONTAINER_FILE"
    print_success "Added field to Container struct"
else
    print_warning "Field already exists in Container struct"
fi

# Add initialization in New function
if ! grep -q "${INIT_CALL}" "$CONTAINER_FILE"; then
    # Find the line with "Initialize modules" comment and add the initialization
    sed -i.bak "/Initialize modules (each module wires its own dependencies)/a\\
	${MODULE_NAME}Module := ${INIT_CALL}
" "$CONTAINER_FILE"
    print_success "Added module initialization to New function"
else
    print_warning "Module initialization already exists in New function"
fi

# Add to return statement
if ! grep -q "${STRUCT_NAME}:" "$CONTAINER_FILE"; then
    # Find the return statement and add the new field
    sed -i.bak "/return &Container{/a\\
		${STRUCT_NAME}: ${MODULE_NAME}Module,
" "$CONTAINER_FILE"
    print_success "Added field to Container return statement"
else
    print_warning "Field already exists in Container return statement"
fi

# Remove backup file
rm -f "${CONTAINER_FILE}.bak"

echo ""
print_info "Updating register_routes.go..."

# Update register_routes.go
ROUTES_FILE="$PROJECT_ROOT/internal/infra/web/register_routes.go"

# Determine the import path and function call based on architecture
if [ "$ARCH_TYPE" = "1" ]; then
    ROUTE_IMPORT_ALIAS="${MODULE_NAME}"
    ROUTE_IMPORT_PATH="${MODULE_PATH}/internal/${MODULE_NAME}"
    ROUTE_CALL="${MODULE_NAME}.RegisterRoutes(router, c.${STRUCT_NAME})"
else
    ROUTE_IMPORT_ALIAS="${MODULE_NAME}Web"
    ROUTE_IMPORT_PATH="${MODULE_PATH}/internal/${MODULE_NAME}/infra/web"
    ROUTE_CALL="${MODULE_NAME}Web.RegisterRoutes(router, c.${STRUCT_NAME})"
fi

# Add import
if ! grep -q "\"$ROUTE_IMPORT_PATH\"" "$ROUTES_FILE"; then
    # Find the last import line and add the new import after it
    sed -i.bak "/^import (/a\\
	${ROUTE_IMPORT_ALIAS} \"${ROUTE_IMPORT_PATH}\"
" "$ROUTES_FILE"
    print_success "Added import to register_routes.go"
else
    print_warning "Import already exists in register_routes.go"
fi

# Add route registration
if ! grep -q "${ROUTE_CALL}" "$ROUTES_FILE"; then
    # Find the line with "Register routes for each module" comment and add the registration
    sed -i.bak "/Register routes for each module/a\\
		${ROUTE_CALL}
" "$ROUTES_FILE"
    print_success "Added route registration to register_routes.go"
else
    print_warning "Route registration already exists in register_routes.go"
fi

# Remove backup file
rm -f "${ROUTES_FILE}.bak"

echo ""
print_success "Module '$MODULE_NAME' created successfully!"
echo ""
print_info "Next steps:"
echo "  1. Implement your domain entities, repositories, services, and controllers"
echo "  2. Define your routes in the routes.go file"
echo "  3. Wire your dependencies in the module.go file"
echo "  4. Test your module by running the application"
echo ""

if [ "$ARCH_TYPE" = "1" ]; then
    print_info "Module location: internal/${MODULE_NAME}/"
    print_info "Architecture: 4-tier (simplified)"
    echo "  - models/       - Data structures"
    echo "  - repositories/ - Database access"
    echo "  - services/     - Business logic"
    echo "  - controllers/  - HTTP handlers"
else
    print_info "Module location: internal/${MODULE_NAME}/"
    print_info "Architecture: DDD (Clean Architecture)"
    echo "  - core/application/repositories/ - Repository interfaces"
    echo "  - core/application/usecases/     - Use cases"
    echo "  - core/domain/entities/          - Domain entities"
    echo "  - infra/repositories/            - Repository implementations"
    echo "  - infra/web/controllers/         - HTTP controllers"
fi

echo ""
