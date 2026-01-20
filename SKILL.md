---
name: git-templater
description: Creating patch-based template repositories for the git-templater CLI tool
---

# Patch-Based Template Repository Architecture

This skill teaches you how to build repositories designed for **compositional templating** using git-templater. Users will select feature branches to compose into new projects, so architecture must minimize merge conflicts while maximizing modularity.

## Core Principle: Spatial Isolation

**Each feature must occupy its own filesystem space.** The primary strategy for avoiding conflicts is ensuring features touch different files.

### File Organization Strategies

#### ✅ Feature-Per-Directory (Preferred)
```
src/
  authentication/
    auth.ext
    middleware.ext
  logging/
    logger.ext
    config.ext
  api-client/
    client.ext
    retry.ext
```

Each feature lives in its own directory. Merging multiple features means adding directories, not modifying shared files.

#### ✅ Feature-Per-File (When Appropriate)
```
config/
  database.ext
  cache.ext
  queue.ext
  email.ext
```

When features are small or naturally file-scoped, one file per feature works well.

#### ❌ Shared Files (Avoid)
```
src/
  core.ext          # Multiple features modify this
  utils.ext         # Multiple features add helpers here
  registry.ext      # Auto-discovery or central registration
```

Any file touched by multiple features creates merge conflicts.

## Forbidden Patterns

### 1. **NO Automatic Module Discovery**

Module discovery requires a central registry or import manifest that every feature modifies.

❌ **Bad: Auto-discovery**
```python
# app.py - touched by every feature
def discover_plugins():
    for module in os.listdir('plugins'):
        import_module(f'plugins.{module}')
```

✅ **Good: Explicit imports**
```python
# features/web_server/app.py
from features.authentication import auth  # explicit
from features.logging import logger        # explicit

def create_app():
    app = Framework()
    auth.configure(app)
    logger.configure(app)
    return app
```

**If a feature needs something, it imports it explicitly.** Dependencies are declared at the call site, not discovered centrally.

### 2. **NO Shared Configuration Objects**

❌ **Bad: Single config file modified by all features**
```yaml
# config.yml - conflicts inevitable
database:
  host: localhost
  port: 5432
cache:
  enabled: true
  ttl: 3600
queue:
  workers: 4
```

✅ **Good: Config per feature**
```
config/
  database.yml      # only database feature touches this
  cache.yml         # only cache feature touches this
  queue.yml         # only queue feature touches this
```

## Git Workflow

### Initialize and Create Base

```bash
git init
git checkout -b master
```

**The master/base branch contains:**
- **Directory structure** - empty directories showing project organization
- **Build/tooling files** - package.json, go.mod, Makefile, etc. (minimal)
- **Entry point** - minimal main file that does nothing but run successfully
- **Testing setup** - test framework config, passing placeholder tests
- **Optional:** Common dev tooling (linting, formatting) if universal to all template uses

**What NOT to put in base:**
- Business logic
- Feature implementations
- Libraries/dependencies specific to features
- Configuration for optional capabilities

The base should answer: "What does every instance of this template absolutely need?"

Commit the base:
```bash
git add .
git commit -m "Initial project skeleton"
```

### Creating Feature Branches

Branch names should use `/` to show hierarchy:

```bash
# Independent features (branch from master)
git checkout master
git checkout -b database/postgresql
git checkout -b auth/jwt
git checkout -b logging/structured

# Alternative implementations (branch from master)
git checkout master
git checkout -b auth/oauth              # alternative to auth/jwt
git checkout -b database/mysql          # alternative to database/postgresql

# Hierarchical features (branch from parent feature)
git checkout database/postgresql
git checkout -b database/postgresql/migrations   # extends postgresql

git checkout auth/jwt
git checkout -b auth/jwt/refresh-tokens          # extends jwt

# Composition features (branch from master, import others)
git checkout master
git checkout -b app/web-api             # imports auth + database + logging
```

**Branching Strategies:**

1. **Independent Features** - Branch from master, add capability in isolation
2. **Alternative Implementations** - Branch from master, mutually exclusive with alternatives
3. **Hierarchical Features** - Branch from parent feature, extend it
4. **Composition Features** - Branch from master, explicitly import multiple features

### Commit Strategy

#### Atomic Feature Commits

Each feature branch should have clear, atomic commits:

```bash
# Good commit messages
"Add PostgreSQL connection module"
"Add JWT authentication middleware"
"Add structured logging with JSON output"
"Add user authentication with JWT and PostgreSQL"  # composition feature
```

#### What Goes in a Feature Commit

**Add:**
- New directories for this feature
- New files specific to this feature
- Dependencies in build files (acceptable conflict)

**Modify:**
- Only files created by this feature OR its parent (if hierarchical)
- Build files when adding dependencies (small, acceptable conflict)

**Never Modify:**
- Files created by sibling features
- Shared registries or manifests
- Other features' configuration files

### Acceptable vs. Unacceptable Conflicts

#### ✅ Acceptable (Small, Mechanical)

```json
// package.json
{
  "dependencies": {
<<<<<<< HEAD
    "express": "^4.18.0"
=======
    "jsonwebtoken": "^9.0.0"
>>>>>>>
  }
}
```

Resolution: Keep both. Mechanical merge.

#### ❌ Unacceptable (Logic Conflicts)

```javascript
// app.js
<<<<<<< HEAD
function initialize() {
    setupAuth()
}
=======
function initialize() {
    setupLogging()
}
>>>>>>>
```

This means two features modified the same logical file. **Architecture failure** - redesign to isolate.

### Testing Feature Independence

Before finalizing a feature branch, verify it merges cleanly to parent:

```bash
# On feature branch: auth/jwt
git checkout master
git merge --no-commit --no-ff auth/jwt

# Check conflicts - should be none or only acceptable (dependencies, docs)
git diff --check
git status

# Abort without committing
git merge --abort

# Return to feature
git checkout auth/jwt
```

**If you find unacceptable conflicts:** Refactor the feature to use spatial isolation.

**Do NOT test combinations** - git-templater users will test feature compositions.

## Summary

The key to patch-based templates is **spatial isolation**:

1. **Each feature occupies unique filesystem space** (directory/file)
2. **Dependencies are explicit, never discovered**
3. **Small conflicts in build files are acceptable**
4. **Logic conflicts indicate architecture failure**
5. **Base branch is minimal skeleton**
6. **Feature branches are independently testable**
7. **Hierarchical branches extend parents**
9. **Work stays local (no push)**