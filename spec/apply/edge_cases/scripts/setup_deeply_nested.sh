#!/bin/bash
set -e
mkdir -p "$1/templates/a/b/c/d/e"
cat > "$1/templates/a/base.patch" << 'PATCH'
diff --git a/a.txt b/a.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/a.txt
@@ -0,0 +1 @@
+a feature
PATCH
cat > "$1/templates/a/b/base.patch" << 'PATCH'
diff --git a/b.txt b/b.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/b.txt
@@ -0,0 +1 @@
+b feature
PATCH
cat > "$1/templates/a/b/c/base.patch" << 'PATCH'
diff --git a/c.txt b/c.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/c.txt
@@ -0,0 +1 @@
+c feature
PATCH
cat > "$1/templates/a/b/c/d/base.patch" << 'PATCH'
diff --git a/d.txt b/d.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/d.txt
@@ -0,0 +1 @@
+d feature
PATCH
cat > "$1/templates/a/b/c/d/e/base.patch" << 'PATCH'
diff --git a/e.txt b/e.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/e.txt
@@ -0,0 +1 @@
+e feature
PATCH
