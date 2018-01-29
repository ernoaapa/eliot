#!/bin/bash -e
# 
# Small script to build and update documentation to GitHub Pages
# 

rm -rf _book
gitbook install
gitbook build
cp docs/CNAME _book/CNAME
cd _book
touch .nojekyll
git init
git add -A
git commit -m 'update book'
git push -f git@github.com:ernoaapa/eliot.git master:gh-pages