#!/bin/zsh

for k in 3
do 
    file='out_solver'$k'.txt'
    rm -fr $file
    for y in {01..10}; 
    do 
        for x in {2..4}; 
        do 
            #echo $x - $y >> $file
            #cat 'pb_'$x'00_'$y'_'$k'_18000.log' | grep 'seconds,\|UNKNOWN\|SATISFIABLE\|UNSATISFIABLE' >> $file
            cat 'pb_'$x'00_'$y'_'$k'_18000.log' | grep 'seconds,' >> $file
        done
    done
done
