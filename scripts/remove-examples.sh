#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=====================================${NC}"
echo -e "${BLUE}  Go App Base - Remove Examples${NC}"
echo -e "${BLUE}=====================================${NC}"
echo ""

echo -e "${YELLOW}⚠️  This will remove all example files from the project.${NC}"
echo -e "${YELLOW}   The following files will be deleted:${NC}"
echo -e "   - internal/core/application/repositories/example_repository.go"
echo -e "   - internal/core/application/usecases/get_example.go"
echo -e "   - internal/core/domain/entities/example.go"
echo -e "   - internal/infra/repositories/example_mysql_repository.go"
echo -e "   - internal/infra/web/controllers/example_controller.go"
echo ""
echo -e "${YELLOW}   And the following files will be modified:${NC}"
echo -e "   - internal/infra/dependencies/dependencies.go"
echo -e "   - internal/infra/web/webserver/gin_handler.go"
echo ""

# Ask for confirmation
echo -e "${YELLOW}Do you want to proceed? (y/n):${NC}"
read -r confirm

if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo -e "${RED}❌ Operation cancelled.${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}🗑️  Removing example files...${NC}"

# Remove example files
files_to_remove=(
    "internal/core/application/repositories/example_repository.go"
    "internal/core/application/usecases/get_example.go"
    "internal/core/domain/entities/example.go"
    "internal/infra/repositories/example_mysql_repository.go"
    "internal/infra/web/controllers/example_controller.go"
)

for file in "${files_to_remove[@]}"; do
    if [ -f "$file" ]; then
        rm "$file"
        echo -e "${GREEN}  ✓ Removed $file${NC}"
    else
        echo -e "${YELLOW}  ⚠ File not found: $file${NC}"
    fi
done

echo ""
echo -e "${BLUE}🔧 Updating dependencies.go...${NC}"

# Update dependencies.go
deps_file="internal/infra/dependencies/dependencies.go"
if [ -f "$deps_file" ]; then
    # Create a temporary file without Example lines
    grep -v "Example" "$deps_file" > "$deps_file.tmp"
    mv "$deps_file.tmp" "$deps_file"
    
    # Format the file
    go fmt "$deps_file" > /dev/null 2>&1 || true
    
    echo -e "${GREEN}  ✓ Updated $deps_file${NC}"
else
    echo -e "${RED}  ❌ File not found: $deps_file${NC}"
fi

echo ""
echo -e "${BLUE}🔧 Updating gin_handler.go...${NC}"

# Update gin_handler.go
handler_file="internal/infra/web/webserver/gin_handler.go"
if [ -f "$handler_file" ]; then
    # Remove the example route (the entire block with 3 lines)
    sed -i.bak '/router.GET("\/examples\/:id"/,/})/d' "$handler_file" && rm "$handler_file.bak"
    
    # Format the file
    go fmt "$handler_file" > /dev/null 2>&1 || true
    
    echo -e "${GREEN}  ✓ Updated $handler_file${NC}"
else
    echo -e "${RED}  ❌ File not found: $handler_file${NC}"
fi

echo ""
echo -e "${BLUE}🔄 Running go mod tidy...${NC}"
go mod tidy

echo ""
echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}✅ Example files removed successfully!${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""
echo -e "${BLUE}📌 Next steps:${NC}"
echo -e "   1. Review the changes"
echo -e "   2. Test your application: make dev"
echo -e "   3. Commit the changes:"
echo -e "      git add ."
echo -e "      git commit -m 'Remove example files'"
echo ""
