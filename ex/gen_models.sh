#!/bin/zsh

mkdir -p $1/cnf

for x in ../data/txt/hard/*; 
do 
    ../bin/gen_alpha -file $x -$1 >  $1/cnf/$(basename $x .txt).cnf
#    minisat+ $1/cnf/$(basename $x .txt).pbo -cnf $1/cnf/$(basename $x .txt).cnf
done
    


