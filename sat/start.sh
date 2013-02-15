#!/bin/zsh

# start.sh
# 1: folder of benchmarks
# 2: folder of output [output]
# 3: runtime in seconds [60]

source prll.sh

data=${1-data/test}
output=${2-output}
time=${3-18000}

mkdir -p $output

argutment=()

for data in *.cnf
do
    for solver in {1,2,3,4,5}
    #for solver in 2
    do
        a=$data' '$solver' '$time' '$output/$(basename $data .cnf)'_'$solver'_'$time'.log'
        argument+=($a)
    done
done

myfn() {
    x1=$(echo $1 | cut -d' ' -f1)
    x2=$(echo $1 | cut -d' ' -f2)
    x3=$(echo $1 | cut -d' ' -f3)
    x4=$(echo $1 | cut -d' ' -f4)
    echo x4
    ./sat.sh $x1 $x2 $x3 | tee $x4
}

prll -c 6 myfn $argument

