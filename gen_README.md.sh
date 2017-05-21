#!/usr/bin/env bash
echo "## a funny middleware onion for [fasthttp](github.com/valyala/fasthttp). inspired by [Alice](https://github.com/justinas/alice)"
echo "[![GoDoc](https://godoc.org/github.com/hanjm/fasthttpmiddleware?status.svg)](https://godoc.org/github.com/hanjm/fasthttpmiddleware)"
echo "[![Go Report Card](https://goreportcard.com/badge/github.com/hanjm/fasthttpmiddleware)](https://goreportcard.com/report/github.com/hanjm/fasthttpmiddleware)"
echo "[![code-coverage](http://gocover.io/_badge/github.com/hanjm/fasthttpmiddleware)](http://gocover.io/github.com/hanjm/fasthttpmiddleware)"
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