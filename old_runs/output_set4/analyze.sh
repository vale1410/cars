#!/bin/zsh

tmp=/tmp/ana.txt

rm -fr $tmp

for y in ../data/set4/*
do  for x in $(basename $y .lp)*.txt; 
    do  echo $y >> $tmp
        cat $x| grep '^UNSAT\|^SAT' >> $tmp
    done
done

cat $tmp | grep -B1 '^UN\|^SAT' | grep 'pb_' | uniq | sort

