#!/bin/bash
set -e
# Feature at providers/oauth/google but providers has no base.patch
mkdir -p "$1/templates/providers/oauth/google"
cat > "$1/templates/providers/oauth/base.patch" << 'PATCH'
diff --git a/oauth.txt b/oauth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/oauth.txt
@@ -0,0 +1 @@
+oauth feature
PATCH
cat > "$1/templates/providers/oauth/google/base.patch" << 'PATCH'
diff --git a/google.txt b/google.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/google.txt
@@ -0,0 +1 @@
+google feature
PATCH
