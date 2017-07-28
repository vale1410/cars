#!/bin/zsh

while read line
do
    for e in {e1,e2,e3}
    do 
        mkdir -p $e
        mkdir -p $e/cnf
        a=("${(s/ /)line}")
        for x in {00..$a[2]}
        do
            typeset -i v=$x
            echo $a[1] $x
            ./encode -f $a[1] -$e -sym -add $v > cnf/$e/$(basename $a[1] .txt)_lb_$x.cnf
        done
    done
done
