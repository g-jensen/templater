#!/bin/bash
set -e
mkdir -p "$1/templates/auth"
cat > "$1/templates/auth/base.patch" << 'PATCH'
---
diff --git a/auth.txt b/auth.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/auth.txt
@@ -0,0 +1 @@
+auth feature
--
2.43.0
PATCH
