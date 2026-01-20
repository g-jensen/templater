#!/bin/bash
set -e
mkdir -p "$1/templates/auth/oauth/google"
mkdir -p "$1/templates/auth/oauth/github"
mkdir -p "$1/templates/database/migrations"
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
cat > "$1/templates/auth/oauth/github/base.patch" << 'PATCH'
diff --git a/github.txt b/github.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/github.txt
@@ -0,0 +1 @@
+github feature
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
cat > "$1/templates/database/migrations/base.patch" << 'PATCH'
diff --git a/migrations.txt b/migrations.txt
new file mode 100644
index 0000000..e69de29
--- /dev/null
+++ b/migrations.txt
@@ -0,0 +1 @@
+migrations feature
PATCH
