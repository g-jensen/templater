#!/bin/bash
set -e
mkdir -p "$1/templates/feature-with-dashes"
mkdir -p "$1/templates/feature_with_underscores"
cat > "$1/templates/feature-with-dashes/base.patch" << 'PATCH'
diff --git a/dashes.txt b/dashes.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/dashes.txt
@@ -0,0 +1 @@
+dashes feature
PATCH
cat > "$1/templates/feature_with_underscores/base.patch" << 'PATCH'
diff --git a/underscores.txt b/underscores.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/underscores.txt
@@ -0,0 +1 @@
+underscores feature
PATCH
