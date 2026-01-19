---
name: templater
description: Creating patch-based template repositories for the templater CLI tool
---

# Templater Template Generation

Reference for creating template repositories that work with the `templater` CLI.

## Core Concept

Templates are **patch-based features** organized in a directory tree. Each feature is a directory containing a `base.patch` file. Directory hierarchy defines dependencies—a feature depends on its parent directory's feature.

## Template Repository Structure

```
my-templates/
  base.patch                    # Optional root feature (applied first if exists)
  auth/
    base.patch                  # "auth" feature, depends on root
    oauth/
      base.patch                # "auth/oauth" feature, depends on "auth"
      google/
        base.patch              # "auth/oauth/google", depends on "auth/oauth"
      github/
        base.patch              # "auth/oauth/github", depends on "auth/oauth"
  database/
    base.patch                  # "database" feature, depends on root
    migrations/
      base.patch                # "database/migrations", depends on "database"
```

## Dependency Rules

1. A feature's dependency is its parent directory's `base.patch` (if it exists)
2. If a parent directory has no `base.patch`, walk up until one is found or reach root
3. Root `base.patch` is optional; if absent, top-level features have no dependencies
4. Sibling features (e.g., `google` and `github`) are independent of each other

## Creating Patches

### Generate a Patch from Git Changes

```bash
# Stage changes for a feature
git add src/auth.js src/middleware.js

# Generate patch file
git diff --cached > templates/auth/base.patch

# Or from committed changes
git diff HEAD~1 > templates/auth/base.patch
```

### Patch File Format

Standard unified diff format:

```diff
diff --git a/src/auth.js b/src/auth.js
new file mode 100644
index 0000000..a1b2c3d
--- /dev/null
+++ b/src/auth.js
@@ -0,0 +1,15 @@
+export function authenticate(token) {
+  if (!token) {
+    throw new Error('Token required');
+  }
+  return validateToken(token);
+}
```

### Patch Types

**New file:**
```diff
diff --git a/newfile.js b/newfile.js
new file mode 100644
--- /dev/null
+++ b/newfile.js
@@ -0,0 +1,3 @@
+content
+goes
+here
```

**Modify existing file:**
```diff
diff --git a/existing.js b/existing.js
--- a/existing.js
+++ b/existing.js
@@ -10,6 +10,8 @@ function existing() {
   return true;
 }
 
+function newFunction() {
+  return 'added';
+}
```

**Delete file:**
```diff
diff --git a/oldfile.js b/oldfile.js
deleted file mode 100644
--- a/oldfile.js
+++ /dev/null
@@ -1,5 +0,0 @@
-content
-that
-was
-here
-before
```

## Designing Feature Hierarchies

### Think in Layers

Organize features from general to specific:

```
api/
  base.patch           # Core API setup (routing, middleware)
  rest/
    base.patch         # REST conventions
  graphql/
    base.patch         # GraphQL setup
    subscriptions/
      base.patch       # WebSocket subscriptions
```

### Share Common Dependencies

When multiple features need the same base:

```
auth/
  base.patch           # Core auth (session, cookies)
  jwt/
    base.patch         # JWT tokens (depends on auth)
  oauth/
    base.patch         # OAuth foundation (depends on auth)
    google/
      base.patch       # Google provider (depends on oauth)
    github/
      base.patch       # GitHub provider (depends on oauth)
```

Applying `auth/oauth/google` automatically applies: `auth` → `auth/oauth` → `auth/oauth/google`

### Keep Features Focused

Each `base.patch` should do one thing well:

```
# GOOD - focused features
database/
  base.patch           # Connection setup only
  migrations/
    base.patch         # Migration tooling
  seeding/
    base.patch         # Seed data utilities

# BAD - monolithic feature
database/
  base.patch           # Connection + migrations + seeds + ORM + caching
```

## Patch Composition Rules

### Patches Must Apply Cleanly

Templater uses `git apply`. Patches fail if:
- Context lines don't match the target file
- File already exists (for new file patches)
- File doesn't exist (for modify patches)

### Order Matters

Patches apply in dependency order (ancestors first). Design patches assuming:
1. Parent features are already applied
2. Sibling features may or may not be applied

### Avoid Conflicts Between Siblings

Sibling features should not modify the same lines:

```
# BAD - both modify the same config section
auth/oauth/google/base.patch  → adds to line 15 of config.js
auth/oauth/github/base.patch  → adds to line 15 of config.js

# GOOD - each modifies distinct sections or files
auth/oauth/google/base.patch  → adds google-config.js
auth/oauth/github/base.patch  → adds github-config.js
```

### Use Extension Points

Design parent features with clear extension points:

```javascript
// In auth/base.patch - creates auth.js with:
const providers = [];

export function registerProvider(provider) {
  providers.push(provider);
}

// In auth/oauth/google/base.patch - creates google.js:
import { registerProvider } from '../auth.js';
registerProvider(googleProvider);
```

## Testing Templates

### Verify with Dry Run

```bash
templater apply ./templates ./test-project auth/oauth/google --dry-run

# Output:
# Would apply:
#   1. auth
#   2. auth/oauth
#   3. auth/oauth/google
```

### Test Fresh Application

```bash
# Create empty target
mkdir test-project && cd test-project && git init

# Apply features
templater apply ../templates . auth/oauth/google

# Verify files exist and work
ls -la src/
```

### Test Incremental Application

```bash
# Apply base first
templater apply ./templates ./project auth

# Later, add more
templater apply ./templates ./project auth/oauth/google
# Should skip "auth" (already applied)
```

## Common Patterns

### Configuration Files

Use patches that append to config arrays:

```diff
diff --git a/config/providers.js b/config/providers.js
--- a/config/providers.js
+++ b/config/providers.js
@@ -1,4 +1,5 @@
 export const providers = [
   // existing providers
+  'google',
 ];
```

### Package Dependencies

Include package.json changes in patches:

```diff
diff --git a/package.json b/package.json
--- a/package.json
+++ b/package.json
@@ -10,6 +10,7 @@
   "dependencies": {
     "express": "^4.18.0",
+    "passport-google-oauth20": "^2.0.0",
   }
 }
```

### Environment Variables

Add to .env.example (never .env):

```diff
diff --git a/.env.example b/.env.example
--- a/.env.example
+++ b/.env.example
@@ -5,3 +5,6 @@ DATABASE_URL=postgres://localhost/myapp
+# Google OAuth
+GOOGLE_CLIENT_ID=
+GOOGLE_CLIENT_SECRET=
```

## Anti-Patterns

### Overlapping Patches

```
# BAD - siblings modify same file section
feature-a/base.patch  → modifies lines 10-15 of index.js
feature-b/base.patch  → modifies lines 12-18 of index.js
```

### Implicit Dependencies

```
# BAD - oauth/google assumes database exists but doesn't depend on it
auth/
  oauth/
    google/
      base.patch  # Uses database connection without declaring dependency
database/
  base.patch
```

### Hardcoded Paths

```diff
# BAD - hardcoded absolute path
+const config = require('/home/user/project/config');

# GOOD - relative path
+const config = require('./config');
```

### Giant Patches

```
# BAD - single patch with 500+ lines
auth/
  base.patch  # Creates 20 files, modifies 10 more

# GOOD - decomposed into focused features
auth/
  base.patch           # Core auth only
  sessions/
    base.patch         # Session management
  passwords/
    base.patch         # Password hashing
```

## Target Project Tracking

Applied features are tracked in `.templater/applied.yml`:

```yaml
applied:
  - auth
  - auth/oauth
  - auth/oauth/google
```

This file is sorted alphabetically. Templater reads it to skip already-applied features.

## CLI Quick Reference

```bash
# List available features
templater list ./templates

# Check what's applied to a project
templater status ./my-project

# Apply features (with dependencies)
templater apply ./templates ./project auth/oauth/google

# Apply multiple features
templater apply ./templates ./project auth/oauth/google database/migrations

# Apply from file
templater apply ./templates ./project -f features.txt

# Preview without applying
templater apply ./templates ./project auth --dry-run
```

## Workflow for Creating a New Template Repository

1. **Start with a working project** - Build the feature manually first
2. **Identify layers** - What depends on what?
3. **Create patches bottom-up** - Start with leaf features, work toward root
4. **Test each patch in isolation** - Apply to fresh project
5. **Test combinations** - Apply sibling features together
6. **Document expected usage** - Which features go together?
