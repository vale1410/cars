#!/bin/zsh

# start.sh
# 1: folder of benchmarks
# 2: folder of output [output]
# 3: runtime in seconds [60]

source prll.sh

data=${1-data/test}
output=${2-output/model9_14400}
time=${3-14400}

mkdir -p $output

argument=()

all=()
all+=data/set2 
all+=data/set3
all+=data/set1
all+=data/set4

echo $all

for x in {1,2,3}
do
    r=$[${RANDOM}%100000]
    for data in $all
    do
        echo $data
        for model in 9
        do
            for conf in {4,3,2}
            do
                if [[ -d $data ]]; then
                    for f in $(ls $data)
                    do
                        mkdir -p $output/$(basename $data)'_'$r
                        a=$data/$f' '$model' '$time' '$conf' '$r' '$output/$(basename $data)'_'$r/$(basename $f .lp)'_'$model'_'$time'_'$conf'.txt'
                        argument+=$a
                    done
                fi
            done
        done
    done
done

myfn() {
    x1=$(echo $1 | cut -d' ' -f1)
    x2=$(echo $1 | cut -d' ' -f2)
    x3=$(echo $1 | cut -d' ' -f3)
    x4=$(echo $1 | cut -d' ' -f4)
    x5=$(echo $1 | cut -d' ' -f5)
    x6=$(echo $1 | cut -d' ' -f6)
    ./run.sh $x1 $x2 $x3 $x4 $x5 $x6
}

prll -c 6 myfn $argument

