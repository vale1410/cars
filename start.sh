#!/bin/zsh

# start.sh
# 1: folder of benchmarks
# 2: folder of output [output]

source prll.sh
time=7200
output=${2-output}

mkdir -p $output

argument=()

for model in {7,3,2,1,6}
do
    for conf in 3
    do
        if [[ -d $1 ]]; then
            for f in $(ls $1)
            do
                a=$1/$f' '$model' '$time' '$conf' '$output/$(basename $f .lp)'_'$model'_'$time'_'$conf'.txt'
                argument+=$a
            done
        elif [[ -f $1 ]]; then
            #argument+=$1' '$model' '$time' '$conf' '$(basename $f .lp) 
        fi
    done
done

myfn() {
    x1=$(echo $1 | cut -d' ' -f1)
    x2=$(echo $1 | cut -d' ' -f2)
    x3=$(echo $1 | cut -d' ' -f3)
    x4=$(echo $1 | cut -d' ' -f4)
    x5=$(echo $1 | cut -d' ' -f5)
    ./run.sh $x1 $x2 $x3 $x4 $x5
}

prll -c 6 myfn $argument

