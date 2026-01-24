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
echo -e "   - internal/example/core/application/repositories/example_repository.go"
echo -e "   - internal/example/core/application/usecases/get_example.go"
echo -e "   - internal/example/core/domain/entities/example.go"
echo -e "   - internal/example/infra/repositories/example_mysql_repository.go"
echo -e "   - internal/example/infra/web/controllers/example_controller.go"
echo -e "   - internal/example/ (entire directory)"
echo ""
echo -e "${YELLOW}   And the following files will be modified:${NC}"
echo -e "   - cmd/server/container/container.go"
echo -e "   - internal/infra/web/register_routes.go"
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

# Remove entire example module directory
if [ -d "internal/example" ]; then
    rm -rf "internal/example"
    echo -e "${GREEN}  ✓ Removed internal/example/${NC}"
else
    echo -e "${YELLOW}  ⚠ Directory not found: internal/example${NC}"
fi

echo ""
echo -e "${BLUE}🔧 Updating container.go...${NC}"

# Update container.go
container_file="cmd/server/container/container.go"
if [ -f "$container_file" ]; then
    # Create a temporary file without Example lines
    grep -v "Example" "$container_file" > "$container_file.tmp"
    mv "$container_file.tmp" "$container_file"
    
    # Format the file
    go fmt "$container_file" > /dev/null 2>&1 || true
    
    echo -e "${GREEN}  ✓ Updated $container_file${NC}"
else
    echo -e "${RED}  ❌ File not found: $container_file${NC}"
fi

echo ""
echo -e "${BLUE}🔧 Updating register_routes.go...${NC}"

# Update register_routes.go
routes_file="internal/infra/web/register_routes.go"
if [ -f "$routes_file" ]; then
    # Remove the example route registration line
    sed -i.bak '/exampleWeb.RegisterRoutes/d' "$routes_file" && rm "$routes_file.bak"
    
    # Format the file
    go fmt "$routes_file" > /dev/null 2>&1 || true
    
    echo -e "${GREEN}  ✓ Updated $routes_file${NC}"
else
    echo -e "${RED}  ❌ File not found: $routes_file${NC}"
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
