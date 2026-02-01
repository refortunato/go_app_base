# Documentation Index

This directory contains comprehensive documentation for the Go App Base project.

## 📚 Documentation Structure

### Implementation Guides
Located in [`implementation/`](./implementation/)

- **[Swagger Guide](./implementation/swagger-guide.md)** - Complete guide for API documentation using Swagger/OpenAPI
  - How to document endpoints
  - Annotation reference
  - Best practices
  - Troubleshooting

- **[Routes Management](./implementation/routes-management.md)** - Route registration patterns
  - Module-based routing
  - Centralized orchestration
  - Best practices

- **[Dependency Management](./implementation/dependency-management.md)** - Dependency injection patterns
  - Container structure
  - Module factories
  - Wiring dependencies

### Development Scripts
Located in [`scripts/`](./scripts/)

- **[Create Module Guide](./scripts/create-module-guide.md)** - Complete guide for creating new modules
  - DDD vs 4-tier architecture
  - Directory structure
  - Step-by-step walkthrough

- **[Create Entity Guide](./scripts/create-entity-guide.md)** - Guide for adding entities to modules
  - CRUD scaffolding
  - Repository patterns
  - Controller generation

## 🚀 Quick Links

### Getting Started
- [Main README](../README.md) - Project overview and setup
- [Copilot Instructions](../.github/copilot-instructions.md) - AI assistant guidelines

### Common Tasks

**Create New Module**
```bash
./scripts/create-module.sh
```
📖 [Full Guide](./scripts/create-module-guide.md)

**Add Entity to Module**
```bash
./scripts/create-entity.sh
```
📖 [Full Guide](./scripts/create-entity-guide.md)

**Generate Swagger Docs**
```bash
make swagger
```
📖 [Full Guide](./implementation/swagger-guide.md)

**View API Documentation**
```
http://localhost:8080/swagger/index.html
```

## 📖 Reading Order for New Developers

1. **[Main README](../README.md)** - Understand project structure and tech stack
2. **[Copilot Instructions](../.github/copilot-instructions.md)** - Learn architecture patterns and conventions
3. **[Create Module Guide](./scripts/create-module-guide.md)** - Create your first module
4. **[Swagger Guide](./implementation/swagger-guide.md)** - Document your API
5. **[Routes Management](./implementation/routes-management.md)** - Understand routing
6. **[Dependency Management](./implementation/dependency-management.md)** - Master DI patterns

## 🔍 Finding What You Need

| I want to... | Read this... |
|-------------|--------------|
| Set up the project | [Main README](../README.md) |
| Understand the architecture | [Copilot Instructions](../.github/copilot-instructions.md) |
| Create a new feature module | [Create Module Guide](./scripts/create-module-guide.md) |
| Add CRUD for an entity | [Create Entity Guide](./scripts/create-entity-guide.md) |
| Document my API | [Swagger Guide](./implementation/swagger-guide.md) |
| Add new routes | [Routes Management](./implementation/routes-management.md) |
| Wire dependencies | [Dependency Management](./implementation/dependency-management.md) |

## 📝 Contributing to Documentation

When adding new features or patterns to the project:

1. **Update relevant guides** - Keep documentation synchronized with code
2. **Add examples** - Show practical usage, not just theory
3. **Update this index** - Add new documentation files here
4. **Follow the format** - Match existing documentation style
5. **Test instructions** - Verify steps actually work

## 🏗️ Documentation Principles

All documentation in this project follows these principles:

✅ **Practical** - Real examples, not abstract theory  
✅ **Complete** - Step-by-step, no assumed knowledge  
✅ **Updated** - Changed with code, not after  
✅ **Searchable** - Clear headings, good structure  
✅ **Concise** - Essential info only, no fluff  

---

**Need help?** Check the guides above or refer to the [Copilot Instructions](../.github/copilot-instructions.md) for comprehensive architecture details.
