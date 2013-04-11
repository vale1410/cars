#!/bin/zsh

while read line
do
    a=("${(s/ /)line}")
    for x in {00..$a[2]}
    do
        typeset -i v=$x
        echo $a[1] $x
        ../gen/gen1 -file $a[1] -e3 -re1 -re2 -add $v > cnf/$(basename $a[1] .txt)_lb_$x.cnf
    done
done
