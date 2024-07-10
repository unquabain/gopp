---
title: Test File 01
data:
    - name: Test 01
      value: 1
    - name: Test 02
      value: 2
    - name: Test 03
      value: 3
---
# This is a test file

The title is {{.title}}.

Here is the rendered test file 02:

{{ include "testfile02.md" .data }}
    
