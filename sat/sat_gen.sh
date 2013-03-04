#!/bin/zsh 

source prll.sh

argument=()

time=1800

for model in {1,2,7,9}
do
    folder=model$model
    mkdir -p $folder
    mkdir -p $folder/cnf
    mkdir -p $folder/log
    for data in ../data/set*/**
    do
        a=$data' '$model' '$folder' '$time
        argument+=($a)
    done
done

myfn() {
    x1=$(echo $1 | cut -d' ' -f1)
    x2=$(echo $1 | cut -d' ' -f2)
    x3=$(echo $1 | cut -d' ' -f3)
    x4=$(echo $1 | cut -d' ' -f4)

    instance=$x3/cnf/$(basename $x1 .lp).cnf
    log=$x3/log/$(basename $x1 .lp).log
    echo $instance
    echo $log

    lp2sat $x1 ../model$x2.lp > $instance
    ./sat.sh $instance 3 $x4 > $log
}

prll -c 6 myfn $argument

