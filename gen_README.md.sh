#!/usr/bin/env bash
echo "## a funny middleware onion for [fasthttp](github.com/valyala/fasthttp). inspired by [Alice](https://github.com/justinas/alice)"
echo ""
echo "### Example"
echo ""
echo "\`\`\`go"
cat example/main.go
echo "\`\`\`"
echo ""
echo "### Document"
echo ""
echo "\`\`\`go"
godoc . | sed '1,8d' | sed 's/SUBDIRECTORIES//' | sed 's/example//'
echo "\`\`\`"
echo ""