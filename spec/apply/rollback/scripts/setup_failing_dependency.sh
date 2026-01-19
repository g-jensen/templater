#!/bin/bash
set -e
mkdir -p "$1/templates/auth/oauth/google"
# auth will work
cat > "$1/templates/auth/base.patch" << 'PATCH'
diff --git a/auth.txt b/auth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/auth.txt
@@ -0,0 +1 @@
+auth feature
PATCH
# oauth will fail
cat > "$1/templates/auth/oauth/base.patch" << 'PATCH'
diff --git a/nonexistent.txt b/nonexistent.txt
index 1234567..abcdefg 100644
--- a/nonexistent.txt
+++ b/nonexistent.txt
@@ -1 +1 @@
-old content
+new content
PATCH
# google would work if oauth worked
cat > "$1/templates/auth/oauth/google/base.patch" << 'PATCH'
diff --git a/google.txt b/google.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/google.txt
@@ -0,0 +1 @@
+google feature
PATCH
