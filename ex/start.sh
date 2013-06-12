#!/bin/zsh

# start.sh
# 1: folder of benchmarks
# 2: folder of output [output]
# 3: runtime in seconds [60]
# 4: seed

source prll.sh
time=1800

seed=123
argument=()
for encoding in {e1,e2,e3}
    do
    mkdir -p $encoding/log
    for instance in $encoding/cnf/*
        do
            a=$instance' '$time' '$encoding/log/$(basename $instance .cnf)'_'$time'_'$seed'.log '$seed
            argument+=($a)
        done
    done

seed=4321

for encoding in {e1,e2,e3}
    do
    mkdir -p $encoding/log
    for instance in $encoding/cnf/*
        do
            a=$instance' '$time' '$encoding/log/$(basename $instance .cnf)'_'$time'_'$seed'.log '$seed
            argument+=($a)
        done
    done

seed=96783

for encoding in {e1,e2,e3}
    do
    mkdir -p $encoding/log
    for instance in $encoding/cnf/*
        do
            a=$instance' '$time' '$encoding/log/$(basename $instance .cnf)'_'$time'_'$seed'.log '$seed
            argument+=($a)
        done
    done


myfn() {
    x1=$(echo $1 | cut -d' ' -f1)
    x2=$(echo $1 | cut -d' ' -f2)
    x3=$(echo $1 | cut -d' ' -f3)
    x4=$(echo $1 | cut -d' ' -f4)
    #echo $x1 $x2 $x3
    minisat $x1 -cpu-lim=$x2 -rnd-init -rnd-seed=$x4 | tee $x3
}

prll -c 3 myfn $argument
