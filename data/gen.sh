#!/bin/zsh

while read line
do
    for e in {e1,e2,e3,e4,e5}
    do 
        mkdir -p $e
        mkdir -p $e/cnf
        a=("${(s/ /)line}")
        for x in {00..$a[2]}
        do
            typeset -i v=$x
            echo $a[1] $x
            ../gen/gen_sat -file $a[1] -$e -sym -add $v > $e/cnf/$(basename $a[1] .txt)_lb_$x.cnf
        done
    done
done
