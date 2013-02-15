#!/bin/zsh

for s in $(ls -1 | grep 'set'| sed 's/^.*_//g' | sort | uniq) 
do 
    echo seed: $s

    for x in *_$s* ; 
    do 
        echo $x; 
        echo sat:; cat $x/* | grep '^SAT'  | wc -l; 
        echo unsat:; cat $x/* | grep '^UNSAT' | wc -l ; 
        echo unknown:;
        cat $x/* | grep '^UNKN' | wc -l ; 
        echo; 
        echo;
    done 
done
