#!/bin/bash

for f in *$1*.cnf 
do 
    x=$(minisat $f | grep "c conflicts" | sed "s/.*: //g" | sed "s/ .*//g")
    echo $f $x
done 
