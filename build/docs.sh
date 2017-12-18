#!/bin/bash
# 
# Small script to build and update documentation to GitHub Pages
# 

npm install gitbook-plugin-insert-logo

rm -rf _book
gitbook install
gitbook build
cp docs/CNAME _book/CNAME
cd _book
git init
git add -A
git commit -m 'update book'
git push -f git@github.com:ernoaapa/eliot.git master:gh-pages