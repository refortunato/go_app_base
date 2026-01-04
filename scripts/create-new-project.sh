#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=====================================${NC}"
echo -e "${BLUE}  Go App Base - Project Generator${NC}"
echo -e "${BLUE}=====================================${NC}"
echo ""

# Function to validate app name
validate_app_name() {
    local name=$1
    if [[ ! "$name" =~ ^[a-z0-9_-]+$ ]]; then
        echo -e "${RED}❌ Invalid name. Use only lowercase letters, numbers, hyphens (-) and underscores (_).${NC}"
        return 1
    fi
    return 0
}

# Function to convert app_name to Title Case for README
app_name_to_title() {
    local name=$1
    # Replace _ and - with spaces, then capitalize first letter of each word
    echo "$name" | sed 's/[_-]/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2));}1'
}

# Ask for application name
while true; do
    echo -e "${YELLOW}📝 Enter the application name (lowercase, no spaces, only _ and - allowed):${NC}"
    read -r app_name
    
    if [ -z "$app_name" ]; then
        echo -e "${RED}❌ Application name cannot be empty.${NC}"
        continue
    fi
    
    if validate_app_name "$app_name"; then
        break
    fi
done

echo -e "${GREEN}✅ Application name: $app_name${NC}"
echo ""

# Check if directory already exists
if [ -d "$app_name" ]; then
    echo -e "${RED}❌ Directory '$app_name' already exists. Please choose another name or remove the existing directory.${NC}"
    exit 1
fi

# Clone the repository
echo -e "${BLUE}📦 Cloning go_app_base repository...${NC}"
git clone git@github.com:refortunato/go_app_base.git "$app_name"

# Navigate to the new project
cd "$app_name"

# Remove git history and initialize new repository
echo -e "${BLUE}🔧 Configuring git...${NC}"
rm -rf .git
git init
git branch -M main

# Ask for remote repository
echo ""
echo -e "${YELLOW}🌐 Do you want to configure a remote repository now? (y/n):${NC}"
read -r configure_remote

remote_repository=""
github_space=""
final_module_path=""

if [[ "$configure_remote" == "y" || "$configure_remote" == "Y" ]]; then
    echo -e "${YELLOW}📝 Enter the remote repository URL (e.g., git@github.com:username/repo.git):${NC}"
    read -r remote_repository
    
    if [ -n "$remote_repository" ]; then
        # Convert git@github.com:user/repo.git to github.com/user/repo
        final_module_path=$(echo "$remote_repository" | sed 's/git@//' | sed 's/\.git$//' | sed 's/:/\//')
        git remote add origin "$remote_repository"
        echo -e "${GREEN}✅ Remote repository configured: $remote_repository${NC}"
    fi
fi

if [ -z "$final_module_path" ]; then
    echo ""
    echo -e "${YELLOW}📝 Enter your GitHub username or organization name:${NC}"
    read -r github_space
    
    if [ -z "$github_space" ]; then
        echo -e "${RED}❌ GitHub space cannot be empty.${NC}"
        exit 1
    fi
    
    final_module_path="github.com/$github_space/$app_name"
    echo -e "${BLUE}ℹ️  Module path will be: $final_module_path${NC}"
fi

echo ""
echo -e "${BLUE}🔄 Updating project files...${NC}"

# Update .env.example
if [ -f "cmd/server/.env.example" ]; then
    sed -i.bak "s/go_app_base/$app_name/g" cmd/server/.env.example && rm cmd/server/.env.example.bak
    echo -e "${GREEN}  ✓ Updated cmd/server/.env.example${NC}"
fi

# Update docker-compose.yaml
if [ -f "docker-compose.yaml" ]; then
    sed -i.bak "s/go_app_base/$app_name/g" docker-compose.yaml && rm docker-compose.yaml.bak
    echo -e "${GREEN}  ✓ Updated docker-compose.yaml${NC}"
fi

# Update schema.sql
if [ -f "schema.sql" ]; then
    sed -i.bak "s/go_app_base/$app_name/g" schema.sql && rm schema.sql.bak
    echo -e "${GREEN}  ✓ Updated schema.sql${NC}"
fi

# Update README.md title
if [ -f "README.md" ]; then
    app_title=$(app_name_to_title "$app_name")
    sed -i.bak "s/Go Application Base/$app_title/g" README.md && rm README.md.bak
    echo -e "${GREEN}  ✓ Updated README.md title to: $app_title${NC}"
fi

# Update all .go files
echo -e "${BLUE}🔄 Updating Go import paths...${NC}"
find . -type f -name "*.go" -exec sed -i.bak "s|github.com/refortunato/go_app_base|$final_module_path|g" {} \; -exec rm {}.bak \;
echo -e "${GREEN}  ✓ Updated all .go files${NC}"

# Update go.mod
if [ -f "go.mod" ]; then
    sed -i.bak "s|github.com/refortunato/go_app_base|$final_module_path|g" go.mod && rm go.mod.bak
    echo -e "${GREEN}  ✓ Updated go.mod${NC}"
fi

# Initial commit
git add .
git commit -m "Initial commit: Project created from go_app_base template"

echo ""
echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}✅ Project '$app_name' created successfully!${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""
echo -e "${BLUE}📁 Project location: $(pwd)${NC}"
echo -e "${BLUE}📦 Module path: $final_module_path${NC}"
echo ""

if [ -n "$remote_repository" ]; then
    echo -e "${YELLOW}📌 Next steps:${NC}"
    echo -e "   1. cd $app_name"
    echo -e "   2. Review the configuration files"
    echo -e "   3. Run: git push -u origin main"
    echo -e "   4. Run: make dev (to start development environment)"
else
    echo -e "${YELLOW}📌 Next steps:${NC}"
    echo -e "   1. cd $app_name"
    echo -e "   2. Configure your remote repository:"
    echo -e "      git remote add origin <your-repository-url>"
    echo -e "   3. Push your code:"
    echo -e "      git push -u origin main"
    echo -e "   4. Run: make dev (to start development environment)"
fi

echo ""
echo -e "${BLUE}💡 To remove example files, run:${NC}"
echo -e "   ./scripts/remove-examples.sh"
echo ""
