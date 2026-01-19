#!/bin/bash
set -e
mkdir -p "$1/templates/auth"
mkdir -p "$1/templates/database"
cat > "$1/templates/auth/base.patch" << 'PATCH'
diff --git a/auth.txt b/auth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/auth.txt
@@ -0,0 +1 @@
+auth feature
PATCH
cat > "$1/templates/database/base.patch" << 'PATCH'
diff --git a/database.txt b/database.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/database.txt
@@ -0,0 +1 @@
+database feature
PATCH
