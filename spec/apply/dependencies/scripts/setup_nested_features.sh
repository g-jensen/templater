#!/bin/bash
set -e
mkdir -p "$1/templates/auth/oauth/google"
cat > "$1/templates/auth/base.patch" << 'PATCH'
diff --git a/auth.txt b/auth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/auth.txt
@@ -0,0 +1 @@
+auth feature
PATCH
cat > "$1/templates/auth/oauth/base.patch" << 'PATCH'
diff --git a/oauth.txt b/oauth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/oauth.txt
@@ -0,0 +1 @@
+oauth feature
PATCH
cat > "$1/templates/auth/oauth/google/base.patch" << 'PATCH'
diff --git a/google.txt b/google.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/google.txt
@@ -0,0 +1 @@
+google feature
PATCH
