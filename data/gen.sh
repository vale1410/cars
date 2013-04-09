#!/bin/zsh

while read line
do
    a=("${(s/ /)line}")
    for x in {00..$a[2]}
    do
        echo gen -file $a[1] -e3 -sym -add $x to $(basename $a[1] .txt)_lb_$x.cnf
    done
done
