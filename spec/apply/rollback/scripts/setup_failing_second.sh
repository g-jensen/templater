#!/bin/bash
set -e
mkdir -p "$1/templates/auth"
mkdir -p "$1/templates/database"
# First patch will work
cat > "$1/templates/auth/base.patch" << 'PATCH'
diff --git a/auth.txt b/auth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/auth.txt
@@ -0,0 +1 @@
+auth feature
PATCH
# Second patch will fail - references a file that doesn't exist
cat > "$1/templates/database/base.patch" << 'PATCH'
diff --git a/nonexistent.txt b/nonexistent.txt
index 1234567..abcdefg 100644
--- a/nonexistent.txt
+++ b/nonexistent.txt
@@ -1 +1 @@
-old content
+new content
PATCH
