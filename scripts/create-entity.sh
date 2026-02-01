#!/bin/bash

# Script to create a new entity with CRUD operations in an existing module
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

# Function to convert string to snake_case
to_snake_case() {
    local str="$1"
    # Convert camelCase/PascalCase to snake_case
    echo "$str" | sed -E 's/([a-z0-9])([A-Z])/\1_\2/g' | tr '[:upper:]' '[:lower:]'
}

# Function to convert snake_case to camelCase
to_camel_case() {
    local str="$1"
    # If already in camelCase (no underscores), return as is
    if [[ ! "$str" =~ _ ]]; then
        echo "$str"
        return
    fi
    # Convert snake_case to camelCase
    echo "$str" | awk -F_ '{printf "%s", $1; for(i=2; i<=NF; i++) printf "%s", toupper(substr($i,1,1)) substr($i,2); print ""}'
}

# Function to capitalize first letter (PascalCase)
capitalize() {
    local str="$1"
    echo "$(echo "${str:0:1}" | tr '[:lower:]' '[:upper:]')${str:1}"
}

# Get the project root directory
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

# List available modules
print_info "Available modules:"
for dir in "$PROJECT_ROOT/internal"/*; do
    if [ -d "$dir" ] && [ "$(basename "$dir")" != "shared" ] && [ "$(basename "$dir")" != "infra" ]; then
        echo "  - $(basename "$dir")"
    fi
done
echo ""

# Ask for module name
print_info "Enter the module name where the entity will be created:"
read -r MODULE_NAME

if [ -z "$MODULE_NAME" ]; then
    print_error "Module name cannot be empty"
    exit 1
fi

MODULE_DIR="$PROJECT_ROOT/internal/$MODULE_NAME"

# Check if module exists
if [ ! -d "$MODULE_DIR" ]; then
    print_error "Module '$MODULE_NAME' does not exist at $MODULE_DIR"
    print_info "Please create the module first using: ./scripts/create-module.sh"
    exit 1
fi

# Detect architecture type
if [ -d "$MODULE_DIR/core/domain" ]; then
    ARCH_TYPE="ddd"
    print_info "Detected architecture: DDD (Clean Architecture)"
elif [ -d "$MODULE_DIR/models" ]; then
    ARCH_TYPE="4tier"
    print_info "Detected architecture: 4-tier (simplified)"
else
    print_error "Could not detect module architecture"
    print_info "Module must have either 'core/domain' (DDD) or 'models' (4-tier) directory"
    exit 1
fi
echo ""

# Ask for entity name
print_info "Enter the entity name (e.g., User, Order, Payment):"
read -r ENTITY_NAME

if [ -z "$ENTITY_NAME" ]; then
    print_error "Entity name cannot be empty"
    exit 1
fi

# Create capitalized and lowercase versions
ENTITY_NAME_LOWER=$(echo "$ENTITY_NAME" | tr '[:upper:]' '[:lower:]')
ENTITY_NAME_CAPITALIZED="$(echo "${ENTITY_NAME:0:1}" | tr '[:lower:]' '[:upper:]')${ENTITY_NAME:1}"

print_info "Entity name: $ENTITY_NAME_CAPITALIZED (lowercase: $ENTITY_NAME_LOWER)"
echo ""

# Function to map MySQL types to Go types
map_mysql_to_go() {
    local mysql_type="$1"
    case "$mysql_type" in
        VARCHAR*|CHAR*|TEXT*|MEDIUMTEXT*|LONGTEXT*|TINYTEXT*)
            echo "string"
            ;;
        INT*|TINYINT*|SMALLINT*|MEDIUMINT*|BIGINT*)
            echo "int"
            ;;
        FLOAT*|DOUBLE*|DECIMAL*)
            echo "float64"
            ;;
        BOOL*|BOOLEAN*)
            echo "bool"
            ;;
        DATE*|DATETIME*|TIMESTAMP*)
            echo "time.Time"
            ;;
        *)
            echo "string"
            ;;
    esac
}

# Ask for fields
print_info "Enter the fields for the entity (format: field_name:TYPE)"
print_info "Examples: name:VARCHAR(255), age:INT, price:DECIMAL(10,2), created_at:TIMESTAMP"
print_info "Type 'done' when finished"
echo ""

declare -a FIELDS
declare -a FIELD_NAMES
declare -a FIELD_TYPES_MYSQL
declare -a FIELD_TYPES_GO

while true; do
    read -p "Field (or 'done'): " field_input
    
    if [ "$field_input" = "done" ]; then
        break
    fi
    
    if [ -z "$field_input" ]; then
        continue
    fi
    
    # Parse field_name:TYPE
    field_name_input=$(echo "$field_input" | cut -d':' -f1)
    field_type_mysql=$(echo "$field_input" | cut -d':' -f2 | tr '[:lower:]' '[:upper:]')
    field_type_go=$(map_mysql_to_go "$field_type_mysql")
    
    # Generate different name formats
    field_name_snake=$(to_snake_case "$field_name_input")      # snake_case for MySQL/JSON
    field_name_camel=$(to_camel_case "$field_name_snake")      # camelCase for Go private fields
    field_name_pascal=$(capitalize "$field_name_camel")         # PascalCase for getters/setters
    
    # Store all formats: snake:mysql:go:camel:pascal
    FIELDS+=("$field_name_snake:$field_type_mysql:$field_type_go:$field_name_camel:$field_name_pascal")
    FIELD_NAMES+=("$field_name_snake")
    FIELD_TYPES_MYSQL+=("$field_type_mysql")
    FIELD_TYPES_GO+=("$field_type_go")
    
    print_success "Added field: $field_name_pascal ($field_type_go) -> DB: $field_name_snake"
done

if [ ${#FIELDS[@]} -eq 0 ]; then
    print_error "At least one field is required"
    exit 1
fi

echo ""
print_info "Creating entity '$ENTITY_NAME_CAPITALIZED' with ${#FIELDS[@]} field(s)..."
echo ""

# Check if any field uses time.Time
HAS_TIME_FIELD=false
for field in "${FIELDS[@]}"; do
    field_type_go=$(echo "$field" | cut -d':' -f3)
    if [ "$field_type_go" = "time.Time" ]; then
        HAS_TIME_FIELD=true
        break
    fi
done

# Create module name capitalized
MODULE_NAME_CAPITALIZED="$(echo "${MODULE_NAME:0:1}" | tr '[:lower:]' '[:upper:]')${MODULE_NAME:1}"

# Generate code based on architecture type
if [ "$ARCH_TYPE" = "ddd" ]; then
    print_info "Generating DDD structure..."
    
    # ==================== 1. Create Entity ====================
    print_info "Creating domain entity..."
    
    ENTITY_FILE="$MODULE_DIR/core/domain/entities/${ENTITY_NAME_LOWER}.go"
    
    cat > "$ENTITY_FILE" <<EOF
package entities

import (
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$ENTITY_FILE"
        echo '' >> "$ENTITY_FILE"
    fi
    
    cat >> "$ENTITY_FILE" <<EOF
	"${MODULE_PATH}/internal/shared"
)

type ${ENTITY_NAME_CAPITALIZED} struct {
	id        string
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_camel}  ${field_type_go}" >> "$ENTITY_FILE"
    done
    
    cat >> "$ENTITY_FILE" <<EOF
	createdAt time.Time
	updatedAt time.Time
}

func New${ENTITY_NAME_CAPITALIZED}(
EOF
    
    # Constructor parameters
    first=true
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        if [ "$first" = true ]; then
            echo "	${field_name_camel} ${field_type_go}," >> "$ENTITY_FILE"
            first=false
        else
            echo "	${field_name_camel} ${field_type_go}," >> "$ENTITY_FILE"
        fi
    done
    
    cat >> "$ENTITY_FILE" <<EOF
) (*${ENTITY_NAME_CAPITALIZED}, error) {
	entity := &${ENTITY_NAME_CAPITALIZED}{
		id:        shared.GenerateId(),
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        echo "		${field_name_camel}: ${field_name_camel}," >> "$ENTITY_FILE"
    done
    
    cat >> "$ENTITY_FILE" <<EOF
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

func Restore${ENTITY_NAME_CAPITALIZED}(
	id string,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_camel} ${field_type_go}," >> "$ENTITY_FILE"
    done
    
    cat >> "$ENTITY_FILE" <<EOF
	createdAt time.Time,
	updatedAt time.Time,
) (*${ENTITY_NAME_CAPITALIZED}, error) {
	return &${ENTITY_NAME_CAPITALIZED}{
		id:        id,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        echo "		${field_name_camel}: ${field_name_camel}," >> "$ENTITY_FILE"
    done
    
    cat >> "$ENTITY_FILE" <<EOF
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (e *${ENTITY_NAME_CAPITALIZED}) Validate() error {
	// TODO: Add validation logic here
	return nil
}

// Getters

func (e *${ENTITY_NAME_CAPITALIZED}) GetId() string {
	return e.id
}

EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        cat >> "$ENTITY_FILE" <<EOF
func (e *${ENTITY_NAME_CAPITALIZED}) Get${field_name_pascal}() ${field_type_go} {
	return e.${field_name_camel}
}

EOF
    done
    
    cat >> "$ENTITY_FILE" <<EOF
func (e *${ENTITY_NAME_CAPITALIZED}) GetCreatedAt() time.Time {
	return e.createdAt
}

func (e *${ENTITY_NAME_CAPITALIZED}) GetUpdatedAt() time.Time {
	return e.updatedAt
}

// Setters

EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        cat >> "$ENTITY_FILE" <<EOF
func (e *${ENTITY_NAME_CAPITALIZED}) Set${field_name_pascal}(${field_name_camel} ${field_type_go}) {
	e.${field_name_camel} = ${field_name_camel}
	e.updatedAt = time.Now().UTC()
}

EOF
    done
    
    print_success "Created entity: $ENTITY_FILE"
    
    # ==================== 2. Create Repository Interface ====================
    print_info "Creating repository interface..."
    
    REPO_INTERFACE_FILE="$MODULE_DIR/core/application/repositories/${ENTITY_NAME_LOWER}_repository.go"
    
    cat > "$REPO_INTERFACE_FILE" <<EOF
package repositories

import (
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/domain/entities"
)

type ${ENTITY_NAME_CAPITALIZED}Repository interface {
	Save(entity *entities.${ENTITY_NAME_CAPITALIZED}) error
	FindById(id string) (*entities.${ENTITY_NAME_CAPITALIZED}, error)
	FindAll(limit, offset int) ([]*entities.${ENTITY_NAME_CAPITALIZED}, error)
	Count() (int, error)
	Update(entity *entities.${ENTITY_NAME_CAPITALIZED}) error
	Delete(id string) error
}
EOF
    
    print_success "Created repository interface: $REPO_INTERFACE_FILE"
    
    # ==================== 3. Create Repository Implementation ====================
    print_info "Creating repository implementation..."
    
    REPO_IMPL_FILE="$MODULE_DIR/infra/repositories/${ENTITY_NAME_LOWER}_mysql_repository.go"
    
    # Start file with package and imports
    {
        echo "package repositories"
        echo ""
        echo "import ("
        echo '	"database/sql"'
        if [ "$HAS_TIME_FIELD" = true ]; then
            echo '	"time"'
        fi
        echo ""
        echo "	\"${MODULE_PATH}/internal/${MODULE_NAME}/core/domain/entities\""
        echo ")"
        echo ""
        echo "type ${ENTITY_NAME_LOWER}Entity struct {"
        echo "	Id        string    \`db:\"id\"\`"
    } > "$REPO_IMPL_FILE"
    
    # Add field definitions
    for field in "${FIELDS[@]}"; do
        field_name_snake=$(echo "$field" | cut -d':' -f1)
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_pascal} ${field_type_go} \`db:\"${field_name_snake}\"\`" >> "$REPO_IMPL_FILE"
    done
    
    # Add CreatedAt and UpdatedAt fields
    {
        echo "	CreatedAt time.Time \`db:\"created_at\"\`"
        echo "	UpdatedAt time.Time \`db:\"updated_at\"\`"
        echo "}"
        echo ""
        echo "type ${ENTITY_NAME_CAPITALIZED}MySQLRepository struct {"
        echo "	db *sql.DB"
        echo "}"
        echo ""
        echo "func New${ENTITY_NAME_CAPITALIZED}MySQLRepository(db *sql.DB) *${ENTITY_NAME_CAPITALIZED}MySQLRepository {"
        echo "	return &${ENTITY_NAME_CAPITALIZED}MySQLRepository{db: db}"
        echo "}"
        echo ""
        echo "func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) Save(entity *entities.${ENTITY_NAME_CAPITALIZED}) error {"
        echo "	query := \`INSERT INTO ${ENTITY_NAME_LOWER}s (id,"
    } >> "$REPO_IMPL_FILE"
    
    # Build field list for INSERT
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name}, created_at, updated_at) VALUES (?, " >> "$REPO_IMPL_FILE"
        else
            echo "		${field_name}, " >> "$REPO_IMPL_FILE"
        fi
    done
    
    # Build placeholder list for INSERT
    for i in "${!FIELDS[@]}"; do
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		?, ?, ?)\`" >> "$REPO_IMPL_FILE"
        else
            echo "		?, " >> "$REPO_IMPL_FILE"
        fi
    done
    
    # Add Exec call
    {
        echo ""
        echo "	_, err := r.db.Exec(query,"
        echo "		entity.GetId(),"
    } >> "$REPO_IMPL_FILE"
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		entity.Get${field_name_pascal}()," >> "$REPO_IMPL_FILE"
    done
    
    # Close Save function
    {
        echo "		entity.GetCreatedAt(),"
        echo "		entity.GetUpdatedAt(),"
        echo "	)"
        echo "	return err"
        echo "}"
        echo ""
        echo "func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) FindById(id string) (*entities.${ENTITY_NAME_CAPITALIZED}, error) {"
        echo "	query := \`SELECT id,"
    } >> "$REPO_IMPL_FILE"
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name}, created_at, updated_at" >> "$REPO_IMPL_FILE"
        else
            echo "		${field_name}, " >> "$REPO_IMPL_FILE"
        fi
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
		FROM ${ENTITY_NAME_LOWER}s WHERE id = ?\`

	var dbEntity ${ENTITY_NAME_LOWER}Entity
	err := r.db.QueryRow(query, id).Scan(
		&dbEntity.Id,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		&dbEntity.${field_name_pascal}," >> "$REPO_IMPL_FILE"
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
		&dbEntity.CreatedAt,
		&dbEntity.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return r.mapToDomain(dbEntity)
}

func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) FindAll(limit, offset int) ([]*entities.${ENTITY_NAME_CAPITALIZED}, error) {
	query := \`SELECT id, 
EOF
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name}, created_at, updated_at" >> "$REPO_IMPL_FILE"
        else
            echo "		${field_name}, " >> "$REPO_IMPL_FILE"
        fi
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
		FROM ${ENTITY_NAME_LOWER}s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?\`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*entities.${ENTITY_NAME_CAPITALIZED}
	for rows.Next() {
		var dbEntity ${ENTITY_NAME_LOWER}Entity
		err := rows.Scan(
			&dbEntity.Id,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "			&dbEntity.${field_name_pascal}," >> "$REPO_IMPL_FILE"
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
			&dbEntity.CreatedAt,
			&dbEntity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		entity, err := r.mapToDomain(dbEntity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) Count() (int, error) {
	query := "SELECT COUNT(*) FROM ${ENTITY_NAME_LOWER}s"
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) Update(entity *entities.${ENTITY_NAME_CAPITALIZED}) error {
	query := \`UPDATE ${ENTITY_NAME_LOWER}s SET 
EOF
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name} = ?, updated_at = ?" >> "$REPO_IMPL_FILE"
        else
            echo "		${field_name} = ?, " >> "$REPO_IMPL_FILE"
        fi
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
		WHERE id = ?\`

	_, err := r.db.Exec(query,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		entity.Get${field_name_pascal}()," >> "$REPO_IMPL_FILE"
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
		entity.GetUpdatedAt(),
		entity.GetId(),
	)
	return err
}

func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) Delete(id string) error {
	query := "DELETE FROM ${ENTITY_NAME_LOWER}s WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *${ENTITY_NAME_CAPITALIZED}MySQLRepository) mapToDomain(dbEntity ${ENTITY_NAME_LOWER}Entity) (*entities.${ENTITY_NAME_CAPITALIZED}, error) {
	return entities.Restore${ENTITY_NAME_CAPITALIZED}(
		dbEntity.Id,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		dbEntity.${field_name_pascal}," >> "$REPO_IMPL_FILE"
    done
    
    cat >> "$REPO_IMPL_FILE" <<EOF
		dbEntity.CreatedAt,
		dbEntity.UpdatedAt,
	)
}
EOF
    
    print_success "Created repository implementation: $REPO_IMPL_FILE"
    
    # ==================== 4. Create Use Cases ====================
    print_info "Creating use cases..."
    
    # Create use case
    CREATE_UC_FILE="$MODULE_DIR/core/application/usecases/create_${ENTITY_NAME_LOWER}.go"
    cat > "$CREATE_UC_FILE" <<EOF
package usecases

import (
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$CREATE_UC_FILE"
        echo '' >> "$CREATE_UC_FILE"
    fi
    
    cat >> "$CREATE_UC_FILE" <<EOF
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/repositories"
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/domain/entities"
)

type Create${ENTITY_NAME_CAPITALIZED}InputDTO struct {
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_pascal} ${field_type_go}" >> "$CREATE_UC_FILE"
    done
    
    cat >> "$CREATE_UC_FILE" <<EOF
}

type Create${ENTITY_NAME_CAPITALIZED}OutputDTO struct {
	Id string \`json:"id"\`
}

type Create${ENTITY_NAME_CAPITALIZED}UseCase struct {
	repository repositories.${ENTITY_NAME_CAPITALIZED}Repository
}

func NewCreate${ENTITY_NAME_CAPITALIZED}UseCase(repo repositories.${ENTITY_NAME_CAPITALIZED}Repository) *Create${ENTITY_NAME_CAPITALIZED}UseCase {
	return &Create${ENTITY_NAME_CAPITALIZED}UseCase{repository: repo}
}

func (uc *Create${ENTITY_NAME_CAPITALIZED}UseCase) Execute(input Create${ENTITY_NAME_CAPITALIZED}InputDTO) (*Create${ENTITY_NAME_CAPITALIZED}OutputDTO, error) {
	entity, err := entities.New${ENTITY_NAME_CAPITALIZED}(
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		input.${field_name_pascal}," >> "$CREATE_UC_FILE"
    done
    
    cat >> "$CREATE_UC_FILE" <<EOF
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repository.Save(entity); err != nil {
		return nil, err
	}

	return &Create${ENTITY_NAME_CAPITALIZED}OutputDTO{Id: entity.GetId()}, nil
}
EOF
    
    print_success "Created use case: $CREATE_UC_FILE"
    
    # Get use case
    GET_UC_FILE="$MODULE_DIR/core/application/usecases/get_${ENTITY_NAME_LOWER}.go"
    cat > "$GET_UC_FILE" <<EOF
package usecases

import (
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$GET_UC_FILE"
        echo '' >> "$GET_UC_FILE"
    fi
    
    cat >> "$GET_UC_FILE" <<EOF
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/repositories"
)

type Get${ENTITY_NAME_CAPITALIZED}InputDTO struct {
	Id string
}

type Get${ENTITY_NAME_CAPITALIZED}OutputDTO struct {
	Id        string    \`json:"id"\`
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        json_tag=$(echo "$field" | cut -d':' -f1)
        echo "	${field_name_pascal} ${field_type_go} \`json:\"${json_tag}\"\`" >> "$GET_UC_FILE"
    done
    
    cat >> "$GET_UC_FILE" <<EOF
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
}

type Get${ENTITY_NAME_CAPITALIZED}UseCase struct {
	repository repositories.${ENTITY_NAME_CAPITALIZED}Repository
}

func NewGet${ENTITY_NAME_CAPITALIZED}UseCase(repo repositories.${ENTITY_NAME_CAPITALIZED}Repository) *Get${ENTITY_NAME_CAPITALIZED}UseCase {
	return &Get${ENTITY_NAME_CAPITALIZED}UseCase{repository: repo}
}

func (uc *Get${ENTITY_NAME_CAPITALIZED}UseCase) Execute(input Get${ENTITY_NAME_CAPITALIZED}InputDTO) (*Get${ENTITY_NAME_CAPITALIZED}OutputDTO, error) {
	entity, err := uc.repository.FindById(input.Id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	return &Get${ENTITY_NAME_CAPITALIZED}OutputDTO{
		Id:        entity.GetId(),
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		${field_name_pascal}: entity.Get${field_name_pascal}()," >> "$GET_UC_FILE"
    done
    
    cat >> "$GET_UC_FILE" <<EOF
		CreatedAt: entity.GetCreatedAt(),
		UpdatedAt: entity.GetUpdatedAt(),
	}, nil
}
EOF
    
    print_success "Created use case: $GET_UC_FILE"
    
    # List use case
    LIST_UC_FILE="$MODULE_DIR/core/application/usecases/list_${ENTITY_NAME_LOWER}.go"
    cat > "$LIST_UC_FILE" <<EOF
package usecases

import (
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$LIST_UC_FILE"
        echo '' >> "$LIST_UC_FILE"
    fi
    
    cat >> "$LIST_UC_FILE" <<EOF
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/repositories"
	"${MODULE_PATH}/internal/shared/dto"
)

type List${ENTITY_NAME_CAPITALIZED}InputDTO struct {
	Page  int
	Limit int
}

type List${ENTITY_NAME_CAPITALIZED}ItemDTO struct {
	Id        string    \`json:"id"\`
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        json_tag=$(echo "$field" | cut -d':' -f1)
        echo "	${field_name_pascal} ${field_type_go} \`json:\"${json_tag}\"\`" >> "$LIST_UC_FILE"
    done
    
    cat >> "$LIST_UC_FILE" <<EOF
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
}

type List${ENTITY_NAME_CAPITALIZED}OutputDTO struct {
	Items      []*List${ENTITY_NAME_CAPITALIZED}ItemDTO \`json:"items"\`
	Pagination *dto.PaginationResponseDTO               \`json:"pagination"\`
}

type List${ENTITY_NAME_CAPITALIZED}UseCase struct {
	repository repositories.${ENTITY_NAME_CAPITALIZED}Repository
}

func NewList${ENTITY_NAME_CAPITALIZED}UseCase(repo repositories.${ENTITY_NAME_CAPITALIZED}Repository) *List${ENTITY_NAME_CAPITALIZED}UseCase {
	return &List${ENTITY_NAME_CAPITALIZED}UseCase{repository: repo}
}

func (uc *List${ENTITY_NAME_CAPITALIZED}UseCase) Execute(input List${ENTITY_NAME_CAPITALIZED}InputDTO) (*List${ENTITY_NAME_CAPITALIZED}OutputDTO, error) {
	// Calculate offset
	offset := (input.Page - 1) * input.Limit

	// Get total count
	totalCount, err := uc.repository.Count()
	if err != nil {
		return nil, err
	}

	// Get entities
	entities, err := uc.repository.FindAll(input.Limit, offset)
	if err != nil {
		return nil, err
	}

	// Build items
	items := make([]*List${ENTITY_NAME_CAPITALIZED}ItemDTO, 0, len(entities))
	for _, entity := range entities {
		items = append(items, &List${ENTITY_NAME_CAPITALIZED}ItemDTO{
			Id:        entity.GetId(),
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "			${field_name_pascal}: entity.Get${field_name_pascal}()," >> "$LIST_UC_FILE"
    done
    
    cat >> "$LIST_UC_FILE" <<EOF
			CreatedAt: entity.GetCreatedAt(),
			UpdatedAt: entity.GetUpdatedAt(),
		})
	}

	// Build pagination
	pagination := dto.NewPaginationResponseDTO(input.Page, input.Limit, totalCount)

	return &List${ENTITY_NAME_CAPITALIZED}OutputDTO{
		Items:      items,
		Pagination: pagination,
	}, nil
}
EOF
    
    print_success "Created use case: $LIST_UC_FILE"
    
    # Update use case
    UPDATE_UC_FILE="$MODULE_DIR/core/application/usecases/update_${ENTITY_NAME_LOWER}.go"
    cat > "$UPDATE_UC_FILE" <<EOF
package usecases

import (
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$UPDATE_UC_FILE"
        echo '' >> "$UPDATE_UC_FILE"
    fi
    
    cat >> "$UPDATE_UC_FILE" <<EOF
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/repositories"
)

type Update${ENTITY_NAME_CAPITALIZED}InputDTO struct {
	Id string
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_pascal} ${field_type_go}" >> "$UPDATE_UC_FILE"
    done
    
    cat >> "$UPDATE_UC_FILE" <<EOF
}

type Update${ENTITY_NAME_CAPITALIZED}OutputDTO struct {
	Success bool \`json:"success"\`
}

type Update${ENTITY_NAME_CAPITALIZED}UseCase struct {
	repository repositories.${ENTITY_NAME_CAPITALIZED}Repository
}

func NewUpdate${ENTITY_NAME_CAPITALIZED}UseCase(repo repositories.${ENTITY_NAME_CAPITALIZED}Repository) *Update${ENTITY_NAME_CAPITALIZED}UseCase {
	return &Update${ENTITY_NAME_CAPITALIZED}UseCase{repository: repo}
}

func (uc *Update${ENTITY_NAME_CAPITALIZED}UseCase) Execute(input Update${ENTITY_NAME_CAPITALIZED}InputDTO) (*Update${ENTITY_NAME_CAPITALIZED}OutputDTO, error) {
	entity, err := uc.repository.FindById(input.Id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "	entity.Set${field_name_pascal}(input.${field_name_pascal})" >> "$UPDATE_UC_FILE"
    done
    
    cat >> "$UPDATE_UC_FILE" <<EOF

	if err := uc.repository.Update(entity); err != nil {
		return nil, err
	}

	return &Update${ENTITY_NAME_CAPITALIZED}OutputDTO{Success: true}, nil
}
EOF
    
    print_success "Created use case: $UPDATE_UC_FILE"
    
    # Delete use case
    DELETE_UC_FILE="$MODULE_DIR/core/application/usecases/delete_${ENTITY_NAME_LOWER}.go"
    cat > "$DELETE_UC_FILE" <<EOF
package usecases

import (
	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/repositories"
)

type Delete${ENTITY_NAME_CAPITALIZED}InputDTO struct {
	Id string
}

type Delete${ENTITY_NAME_CAPITALIZED}OutputDTO struct {
	Success bool \`json:"success"\`
}

type Delete${ENTITY_NAME_CAPITALIZED}UseCase struct {
	repository repositories.${ENTITY_NAME_CAPITALIZED}Repository
}

func NewDelete${ENTITY_NAME_CAPITALIZED}UseCase(repo repositories.${ENTITY_NAME_CAPITALIZED}Repository) *Delete${ENTITY_NAME_CAPITALIZED}UseCase {
	return &Delete${ENTITY_NAME_CAPITALIZED}UseCase{repository: repo}
}

func (uc *Delete${ENTITY_NAME_CAPITALIZED}UseCase) Execute(input Delete${ENTITY_NAME_CAPITALIZED}InputDTO) (*Delete${ENTITY_NAME_CAPITALIZED}OutputDTO, error) {
	if err := uc.repository.Delete(input.Id); err != nil {
		return nil, err
	}

	return &Delete${ENTITY_NAME_CAPITALIZED}OutputDTO{Success: true}, nil
}
EOF
    
    print_success "Created use case: $DELETE_UC_FILE"
    
    # ==================== 5. Create Controller ====================
    print_info "Creating controller..."
    
    CONTROLLER_FILE="$MODULE_DIR/infra/web/controllers/${ENTITY_NAME_LOWER}_controller.go"
    cat > "$CONTROLLER_FILE" <<EOF
package controllers

import (
	"net/http"

	"${MODULE_PATH}/internal/${MODULE_NAME}/core/application/usecases"
	"${MODULE_PATH}/internal/shared/dto"
	"${MODULE_PATH}/internal/shared/web/advisor"
	"${MODULE_PATH}/internal/shared/web/context"
)

type ${ENTITY_NAME_CAPITALIZED}Controller struct {
	createUseCase *usecases.Create${ENTITY_NAME_CAPITALIZED}UseCase
	getUseCase    *usecases.Get${ENTITY_NAME_CAPITALIZED}UseCase
	listUseCase   *usecases.List${ENTITY_NAME_CAPITALIZED}UseCase
	updateUseCase *usecases.Update${ENTITY_NAME_CAPITALIZED}UseCase
	deleteUseCase *usecases.Delete${ENTITY_NAME_CAPITALIZED}UseCase
}

func New${ENTITY_NAME_CAPITALIZED}Controller(
	createUseCase usecases.Create${ENTITY_NAME_CAPITALIZED}UseCase,
	getUseCase usecases.Get${ENTITY_NAME_CAPITALIZED}UseCase,
	listUseCase usecases.List${ENTITY_NAME_CAPITALIZED}UseCase,
	updateUseCase usecases.Update${ENTITY_NAME_CAPITALIZED}UseCase,
	deleteUseCase usecases.Delete${ENTITY_NAME_CAPITALIZED}UseCase,
) *${ENTITY_NAME_CAPITALIZED}Controller {
	return &${ENTITY_NAME_CAPITALIZED}Controller{
		createUseCase: &createUseCase,
		getUseCase:    &getUseCase,
		listUseCase:   &listUseCase,
		updateUseCase: &updateUseCase,
		deleteUseCase: &deleteUseCase,
	}
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Create(ctx context.WebContext) {
	var request usecases.Create${ENTITY_NAME_CAPITALIZED}InputDTO
	if err := ctx.BindJSON(&request); err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	output, err := c.createUseCase.Execute(request)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, output)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Get(ctx context.WebContext) {
	id := ctx.Param("id")
	
	output, err := c.getUseCase.Execute(usecases.Get${ENTITY_NAME_CAPITALIZED}InputDTO{Id: id})
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}
	
	if output == nil {
		advisor.ReturnNotFoundError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) List(ctx context.WebContext) {
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	pagination, err := dto.NewPaginationRequestDTO(pageStr, limitStr)
	if err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	output, err := c.listUseCase.Execute(usecases.List${ENTITY_NAME_CAPITALIZED}InputDTO{
		Page:  pagination.Page,
		Limit: pagination.Limit,
	})
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Update(ctx context.WebContext) {
	id := ctx.Param("id")
	
	var request usecases.Update${ENTITY_NAME_CAPITALIZED}InputDTO
	if err := ctx.BindJSON(&request); err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}
	request.Id = id

	output, err := c.updateUseCase.Execute(request)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}
	
	if output == nil {
		advisor.ReturnNotFoundError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, output)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Delete(ctx context.WebContext) {
	id := ctx.Param("id")

	output, err := c.deleteUseCase.Execute(usecases.Delete${ENTITY_NAME_CAPITALIZED}InputDTO{Id: id})
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, output)
}
EOF
    
    print_success "Created controller: $CONTROLLER_FILE"
    
    # ==================== 6. Update module.go ====================
    print_info "Updating module.go..."
    
    MODULE_FILE="$MODULE_DIR/infra/module.go"
    
    # Add to struct
    sed -i.bak "/type ${MODULE_NAME_CAPITALIZED}Module struct {/a\\
	${ENTITY_NAME_CAPITALIZED}Controller *controllers.${ENTITY_NAME_CAPITALIZED}Controller
" "$MODULE_FILE"
    
    # Add to New function
    sed -i.bak "/return &${MODULE_NAME_CAPITALIZED}Module{/i\\
\\
	// ${ENTITY_NAME_CAPITALIZED} CRUD\\
	${ENTITY_NAME_LOWER}Repo := repositories.New${ENTITY_NAME_CAPITALIZED}MySQLRepository(db)\\
	create${ENTITY_NAME_CAPITALIZED}UC := usecases.NewCreate${ENTITY_NAME_CAPITALIZED}UseCase(${ENTITY_NAME_LOWER}Repo)\\
	get${ENTITY_NAME_CAPITALIZED}UC := usecases.NewGet${ENTITY_NAME_CAPITALIZED}UseCase(${ENTITY_NAME_LOWER}Repo)\\
	list${ENTITY_NAME_CAPITALIZED}UC := usecases.NewList${ENTITY_NAME_CAPITALIZED}UseCase(${ENTITY_NAME_LOWER}Repo)\\
	update${ENTITY_NAME_CAPITALIZED}UC := usecases.NewUpdate${ENTITY_NAME_CAPITALIZED}UseCase(${ENTITY_NAME_LOWER}Repo)\\
	delete${ENTITY_NAME_CAPITALIZED}UC := usecases.NewDelete${ENTITY_NAME_CAPITALIZED}UseCase(${ENTITY_NAME_LOWER}Repo)\\
	${ENTITY_NAME_LOWER}Controller := controllers.New${ENTITY_NAME_CAPITALIZED}Controller(*create${ENTITY_NAME_CAPITALIZED}UC, *get${ENTITY_NAME_CAPITALIZED}UC, *list${ENTITY_NAME_CAPITALIZED}UC, *update${ENTITY_NAME_CAPITALIZED}UC, *delete${ENTITY_NAME_CAPITALIZED}UC)
" "$MODULE_FILE"
    
    sed -i.bak "/return &${MODULE_NAME_CAPITALIZED}Module{/a\\
		${ENTITY_NAME_CAPITALIZED}Controller: ${ENTITY_NAME_LOWER}Controller,
" "$MODULE_FILE"
    
    rm -f "${MODULE_FILE}.bak"
    print_success "Updated module.go"
    
    # ==================== 7. Update routes.go ====================
    print_info "Updating routes.go..."
    
    ROUTES_FILE="$MODULE_DIR/infra/web/routes.go"
    
    sed -i.bak "/func RegisterRoutes/a\\
\\
	// ${ENTITY_NAME_CAPITALIZED} routes\\
	router.POST(\"/${ENTITY_NAME_LOWER}s\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Create(context.NewGinContextAdapter(ctx))\\
	})\\
	router.GET(\"/${ENTITY_NAME_LOWER}s/:id\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Get(context.NewGinContextAdapter(ctx))\\
	})\\
	router.GET(\"/${ENTITY_NAME_LOWER}s\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.List(context.NewGinContextAdapter(ctx))\\
	})\\
	router.PUT(\"/${ENTITY_NAME_LOWER}s/:id\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Update(context.NewGinContextAdapter(ctx))\\
	})\\
	router.DELETE(\"/${ENTITY_NAME_LOWER}s/:id\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Delete(context.NewGinContextAdapter(ctx))\\
	})
" "$ROUTES_FILE"
    
    rm -f "${ROUTES_FILE}.bak"
    print_success "Updated routes.go"
    
else
    # ==================== 4-TIER ARCHITECTURE ====================
    print_info "Generating 4-tier structure..."
    
    # ==================== 1. Create Model ====================
    print_info "Creating model..."
    
    MODEL_FILE="$MODULE_DIR/models/${ENTITY_NAME_LOWER}.go"
    cat > "$MODEL_FILE" <<EOF
package models

EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo 'import "time"' >> "$MODEL_FILE"
        echo '' >> "$MODEL_FILE"
    fi
    
    cat >> "$MODEL_FILE" <<EOF
type ${ENTITY_NAME_CAPITALIZED} struct {
	ID        string
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_pascal} ${field_type_go}" >> "$MODEL_FILE"
    done
    
    cat >> "$MODEL_FILE" <<EOF
	CreatedAt time.Time
	UpdatedAt time.Time
}
EOF
    
    print_success "Created model: $MODEL_FILE"
    
    # ==================== 2. Create Repository ====================
    print_info "Creating repository..."
    
    REPO_FILE="$MODULE_DIR/repositories/${ENTITY_NAME_LOWER}_repository.go"
    cat > "$REPO_FILE" <<EOF
package repositories

import (
	"database/sql"
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$REPO_FILE"
        echo '' >> "$REPO_FILE"
    fi
    
    cat >> "$REPO_FILE" <<EOF
	"${MODULE_PATH}/internal/${MODULE_NAME}/models"
)

type ${ENTITY_NAME_CAPITALIZED}Repository struct {
	db *sql.DB
}

func New${ENTITY_NAME_CAPITALIZED}Repository(db *sql.DB) *${ENTITY_NAME_CAPITALIZED}Repository {
	return &${ENTITY_NAME_CAPITALIZED}Repository{db: db}
}

func (r *${ENTITY_NAME_CAPITALIZED}Repository) FindById(id string) (*models.${ENTITY_NAME_CAPITALIZED}, error) {
	query := \`SELECT id, 
EOF
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name}, created_at, updated_at" >> "$REPO_FILE"
        else
            echo "		${field_name}, " >> "$REPO_FILE"
        fi
    done
    
    cat >> "$REPO_FILE" <<EOF
		FROM ${ENTITY_NAME_LOWER}s WHERE id = ?\`

	var entity models.${ENTITY_NAME_CAPITALIZED}
	err := r.db.QueryRow(query, id).Scan(
		&entity.ID,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		&entity.${field_name_pascal}," >> "$REPO_FILE"
    done
    
    cat >> "$REPO_FILE" <<EOF
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

func (r *${ENTITY_NAME_CAPITALIZED}Repository) FindAll(limit, offset int) ([]*models.${ENTITY_NAME_CAPITALIZED}, error) {
	query := \`SELECT id, 
EOF
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name}, created_at, updated_at" >> "$REPO_FILE"
        else
            echo "		${field_name}, " >> "$REPO_FILE"
        fi
    done
    
    cat >> "$REPO_FILE" <<EOF
		FROM ${ENTITY_NAME_LOWER}s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?\`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*models.${ENTITY_NAME_CAPITALIZED}
	for rows.Next() {
		var entity models.${ENTITY_NAME_CAPITALIZED}
		err := rows.Scan(
			&entity.ID,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "			&entity.${field_name_pascal}," >> "$REPO_FILE"
    done
    
    cat >> "$REPO_FILE" <<EOF
			&entity.CreatedAt,
			&entity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}

	return entities, nil
}

func (r *${ENTITY_NAME_CAPITALIZED}Repository) Count() (int, error) {
	query := "SELECT COUNT(*) FROM ${ENTITY_NAME_LOWER}s"
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *${ENTITY_NAME_CAPITALIZED}Repository) Save(entity *models.${ENTITY_NAME_CAPITALIZED}) error {
	query := \`INSERT INTO ${ENTITY_NAME_LOWER}s (id, 
EOF
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name}, created_at, updated_at) VALUES (?, " >> "$REPO_FILE"
        else
            echo "		${field_name}, " >> "$REPO_FILE"
        fi
    done
    
    for i in "${!FIELDS[@]}"; do
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		?, ?, ?)\`" >> "$REPO_FILE"
        else
            echo "		?, " >> "$REPO_FILE"
        fi
    done
    
    cat >> "$REPO_FILE" <<EOF

	_, err := r.db.Exec(query,
		entity.ID,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		entity.${field_name_pascal}," >> "$REPO_FILE"
    done
    
    cat >> "$REPO_FILE" <<EOF
		entity.CreatedAt,
		entity.UpdatedAt,
	)
	return err
}

func (r *${ENTITY_NAME_CAPITALIZED}Repository) Update(entity *models.${ENTITY_NAME_CAPITALIZED}) error {
	query := \`UPDATE ${ENTITY_NAME_LOWER}s SET 
EOF
    
    for i in "${!FIELDS[@]}"; do
        field_name=$(echo "${FIELDS[$i]}" | cut -d':' -f1)
        if [ $i -eq $((${#FIELDS[@]} - 1)) ]; then
            echo "		${field_name} = ?, updated_at = ?" >> "$REPO_FILE"
        else
            echo "		${field_name} = ?, " >> "$REPO_FILE"
        fi
    done
    
    cat >> "$REPO_FILE" <<EOF
		WHERE id = ?\`

	_, err := r.db.Exec(query,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		entity.${field_name_pascal}," >> "$REPO_FILE"
    done
    
    cat >> "$REPO_FILE" <<EOF
		entity.UpdatedAt,
		entity.ID,
	)
	return err
}

func (r *${ENTITY_NAME_CAPITALIZED}Repository) Delete(id string) error {
	query := "DELETE FROM ${ENTITY_NAME_LOWER}s WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
EOF
    
    print_success "Created repository: $REPO_FILE"
    
    # ==================== 3. Create Service ====================
    print_info "Creating service..."
    
    SERVICE_FILE="$MODULE_DIR/services/${ENTITY_NAME_LOWER}_service.go"
    cat > "$SERVICE_FILE" <<EOF
package services

import (
	"fmt"
EOF
    
    if [ "$HAS_TIME_FIELD" = true ]; then
        echo '	"time"' >> "$SERVICE_FILE"
        echo '' >> "$SERVICE_FILE"
    fi
    
    cat >> "$SERVICE_FILE" <<EOF
	"${MODULE_PATH}/internal/shared"
	"${MODULE_PATH}/internal/${MODULE_NAME}/models"
	"${MODULE_PATH}/internal/${MODULE_NAME}/repositories"
)

type ${ENTITY_NAME_CAPITALIZED}Service struct {
	repository *repositories.${ENTITY_NAME_CAPITALIZED}Repository
}

func New${ENTITY_NAME_CAPITALIZED}Service(repo *repositories.${ENTITY_NAME_CAPITALIZED}Repository) *${ENTITY_NAME_CAPITALIZED}Service {
	return &${ENTITY_NAME_CAPITALIZED}Service{repository: repo}
}

func (s *${ENTITY_NAME_CAPITALIZED}Service) Get${ENTITY_NAME_CAPITALIZED}(id string) (*models.${ENTITY_NAME_CAPITALIZED}, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	entity, err := s.repository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ${ENTITY_NAME_LOWER}: %w", err)
	}

	if entity == nil {
		return nil, fmt.Errorf("${ENTITY_NAME_LOWER} not found")
	}

	return entity, nil
}

func (s *${ENTITY_NAME_CAPITALIZED}Service) List${ENTITY_NAME_CAPITALIZED}s(page, limit int) (map[string]interface{}, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Get total count
	totalCount, err := s.repository.Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count ${ENTITY_NAME_LOWER}s: %w", err)
	}

	// Get entities
	entities, err := s.repository.FindAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list ${ENTITY_NAME_LOWER}s: %w", err)
	}

	// Calculate total pages
	totalPages := 0
	if totalCount > 0 {
		totalPages = (totalCount + limit - 1) / limit
	}

	return map[string]interface{}{
		"items": entities,
		"pagination": map[string]int{
			"page":        page,
			"limit":       limit,
			"total_items": totalCount,
			"total_pages": totalPages,
		},
	}, nil
}

func (s *${ENTITY_NAME_CAPITALIZED}Service) Create${ENTITY_NAME_CAPITALIZED}(
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_camel} ${field_type_go}," >> "$SERVICE_FILE"
    done
    
    cat >> "$SERVICE_FILE" <<EOF
) (*models.${ENTITY_NAME_CAPITALIZED}, error) {
	// TODO: Add validation logic here

	now := time.Now().UTC()
	entity := &models.${ENTITY_NAME_CAPITALIZED}{
		ID:        shared.GenerateId(),
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		${field_name_pascal}: ${field_name_camel}," >> "$SERVICE_FILE"
    done
    
    cat >> "$SERVICE_FILE" <<EOF
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Save(entity); err != nil {
		return nil, fmt.Errorf("failed to create ${ENTITY_NAME_LOWER}: %w", err)
	}

	return entity, nil
}

func (s *${ENTITY_NAME_CAPITALIZED}Service) Update${ENTITY_NAME_CAPITALIZED}(
	id string,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        echo "	${field_name_camel} ${field_type_go}," >> "$SERVICE_FILE"
    done
    
    cat >> "$SERVICE_FILE" <<EOF
) (*models.${ENTITY_NAME_CAPITALIZED}, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	existing, err := s.repository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ${ENTITY_NAME_LOWER}: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("${ENTITY_NAME_LOWER} not found")
	}

	// TODO: Add validation logic here

EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_camel=$(echo "$field" | cut -d':' -f4)
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "	existing.${field_name_pascal} = ${field_name_camel}" >> "$SERVICE_FILE"
    done
    
    cat >> "$SERVICE_FILE" <<EOF
	existing.UpdatedAt = time.Now().UTC()

	if err := s.repository.Update(existing); err != nil {
		return nil, fmt.Errorf("failed to update ${ENTITY_NAME_LOWER}: %w", err)
	}

	return existing, nil
}

func (s *${ENTITY_NAME_CAPITALIZED}Service) Delete${ENTITY_NAME_CAPITALIZED}(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete ${ENTITY_NAME_LOWER}: %w", err)
	}

	return nil
}
EOF
    
    print_success "Created service: $SERVICE_FILE"
    
    # ==================== 4. Create Controller ====================
    print_info "Creating controller..."
    
    CONTROLLER_FILE="$MODULE_DIR/controllers/${ENTITY_NAME_LOWER}_controller.go"
    cat > "$CONTROLLER_FILE" <<EOF
package controllers

import (
	"net/http"

	"${MODULE_PATH}/internal/shared/dto"
	"${MODULE_PATH}/internal/shared/web/advisor"
	"${MODULE_PATH}/internal/shared/web/context"
	"${MODULE_PATH}/internal/${MODULE_NAME}/services"
)

type ${ENTITY_NAME_CAPITALIZED}Controller struct {
	service *services.${ENTITY_NAME_CAPITALIZED}Service
}

func New${ENTITY_NAME_CAPITALIZED}Controller(service *services.${ENTITY_NAME_CAPITALIZED}Service) *${ENTITY_NAME_CAPITALIZED}Controller {
	return &${ENTITY_NAME_CAPITALIZED}Controller{service: service}
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Get(ctx context.WebContext) {
	id := ctx.Param("id")

	entity, err := c.service.Get${ENTITY_NAME_CAPITALIZED}(id)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) List(ctx context.WebContext) {
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	pagination, err := dto.NewPaginationRequestDTO(pageStr, limitStr)
	if err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	result, err := c.service.List${ENTITY_NAME_CAPITALIZED}s(pagination.Page, pagination.Limit)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Create(ctx context.WebContext) {
	var request struct {
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        json_tag=$(echo "$field" | cut -d':' -f1)
        echo "		${field_name_pascal} ${field_type_go} \`json:\"${json_tag}\"\`" >> "$CONTROLLER_FILE"
    done
    
    cat >> "$CONTROLLER_FILE" <<EOF
	}

	if err := ctx.BindJSON(&request); err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	entity, err := c.service.Create${ENTITY_NAME_CAPITALIZED}(
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		request.${field_name_pascal}," >> "$CONTROLLER_FILE"
    done
    
    cat >> "$CONTROLLER_FILE" <<EOF
	)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, entity)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Update(ctx context.WebContext) {
	id := ctx.Param("id")

	var request struct {
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        field_type_go=$(echo "$field" | cut -d':' -f3)
        json_tag=$(echo "$field" | cut -d':' -f1)
        echo "		${field_name_pascal} ${field_type_go} \`json:\"${json_tag}\"\`" >> "$CONTROLLER_FILE"
    done
    
    cat >> "$CONTROLLER_FILE" <<EOF
	}

	if err := ctx.BindJSON(&request); err != nil {
		advisor.ReturnBadRequestError(ctx, err)
		return
	}

	entity, err := c.service.Update${ENTITY_NAME_CAPITALIZED}(
		id,
EOF
    
    for field in "${FIELDS[@]}"; do
        field_name_pascal=$(echo "$field" | cut -d':' -f5)
        echo "		request.${field_name_pascal}," >> "$CONTROLLER_FILE"
    done
    
    cat >> "$CONTROLLER_FILE" <<EOF
	)
	if err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, entity)
}

func (c *${ENTITY_NAME_CAPITALIZED}Controller) Delete(ctx context.WebContext) {
	id := ctx.Param("id")

	if err := c.service.Delete${ENTITY_NAME_CAPITALIZED}(id); err != nil {
		advisor.ReturnApplicationError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]bool{"success": true})
}
EOF
    
    print_success "Created controller: $CONTROLLER_FILE"
    
    # ==================== 5. Update module.go ====================
    print_info "Updating module.go..."
    
    MODULE_FILE="$MODULE_DIR/module.go"
    
    sed -i.bak "/type ${MODULE_NAME_CAPITALIZED}Module struct {/a\\
	${ENTITY_NAME_CAPITALIZED}Controller *controllers.${ENTITY_NAME_CAPITALIZED}Controller
" "$MODULE_FILE"
    
    sed -i.bak "/return &${MODULE_NAME_CAPITALIZED}Module{/i\\
\\
	// ${ENTITY_NAME_CAPITALIZED} CRUD\\
	${ENTITY_NAME_LOWER}Repo := repositories.New${ENTITY_NAME_CAPITALIZED}Repository(db)\\
	${ENTITY_NAME_LOWER}Service := services.New${ENTITY_NAME_CAPITALIZED}Service(${ENTITY_NAME_LOWER}Repo)\\
	${ENTITY_NAME_LOWER}Controller := controllers.New${ENTITY_NAME_CAPITALIZED}Controller(${ENTITY_NAME_LOWER}Service)
" "$MODULE_FILE"
    
    sed -i.bak "/return &${MODULE_NAME_CAPITALIZED}Module{/a\\
		${ENTITY_NAME_CAPITALIZED}Controller: ${ENTITY_NAME_LOWER}Controller,
" "$MODULE_FILE"
    
    rm -f "${MODULE_FILE}.bak"
    print_success "Updated module.go"
    
    # ==================== 6. Update routes.go ====================
    print_info "Updating routes.go..."
    
    ROUTES_FILE="$MODULE_DIR/routes.go"
    
    sed -i.bak "/func RegisterRoutes/a\\
\\
	// ${ENTITY_NAME_CAPITALIZED} routes\\
	router.POST(\"/${ENTITY_NAME_LOWER}s\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Create(context.NewGinContextAdapter(ctx))\\
	})\\
	router.GET(\"/${ENTITY_NAME_LOWER}s/:id\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Get(context.NewGinContextAdapter(ctx))\\
	})\\
	router.GET(\"/${ENTITY_NAME_LOWER}s\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.List(context.NewGinContextAdapter(ctx))\\
	})\\
	router.PUT(\"/${ENTITY_NAME_LOWER}s/:id\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Update(context.NewGinContextAdapter(ctx))\\
	})\\
	router.DELETE(\"/${ENTITY_NAME_LOWER}s/:id\", func(ctx *gin.Context) {\\
		module.${ENTITY_NAME_CAPITALIZED}Controller.Delete(context.NewGinContextAdapter(ctx))\\
	})
" "$ROUTES_FILE"
    
    rm -f "${ROUTES_FILE}.bak"
    print_success "Updated routes.go"
fi

echo ""
print_success "Entity '$ENTITY_NAME_CAPITALIZED' created successfully!"
echo ""
print_info "SQL to create the table:"
echo ""
echo "CREATE TABLE ${ENTITY_NAME_LOWER}s ("
echo "    id VARCHAR(36) PRIMARY KEY,"

for field in "${FIELDS[@]}"; do
    field_name_snake=$(echo "$field" | cut -d':' -f1)
    field_type_mysql=$(echo "$field" | cut -d':' -f2)
    echo "    ${field_name_snake} ${field_type_mysql},"
done

echo "    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,"
echo "    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"
echo ");"
echo ""

print_info "API Endpoints created:"
echo "  POST   /${ENTITY_NAME_LOWER}s          - Create new ${ENTITY_NAME_LOWER}"
echo "  GET    /${ENTITY_NAME_LOWER}s/:id      - Get ${ENTITY_NAME_LOWER} by ID"
echo "  GET    /${ENTITY_NAME_LOWER}s          - List all ${ENTITY_NAME_LOWER}s (with pagination)"
echo "  PUT    /${ENTITY_NAME_LOWER}s/:id      - Update ${ENTITY_NAME_LOWER}"
echo "  DELETE /${ENTITY_NAME_LOWER}s/:id      - Delete ${ENTITY_NAME_LOWER}"
echo ""
