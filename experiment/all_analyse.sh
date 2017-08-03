#!/bin/zsh

for x in res-*; do rm -fr $x; done 
for x in exp-*; do mkdir res-${x#exp-}; done 

for x in exp-*; do analyze -f $x -o res-${x#exp-}; done
for x in res-*; do cd $x; pdflatex all.tex; cd ..; done
