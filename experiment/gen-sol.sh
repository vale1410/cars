#!/bin/zsh

instances=$1
basefolder=exp-$2
timelimit=$3

mkdir -p $basefolder

solver=(minisat microsat glucose lingeling cmsat clasp)
#solver=(glucose lingeling cmsat clasp)

case $4 in
    1) seed=(142) ;;
    2) seed=(142 321) ;;
    3) seed=(142 321 832) ;;
    4) seed=(142 321 832 999) ;;
esac 


for se in $seed
do 
    for sol in $solver
    do 
       folder=$basefolder/$sol
       mkdir -p $folder
       for x in $instances/*
       do 
           output=$folder/$(basename $x .cnf)-$se.log
           echo msat -seed $se -solver $sol -time $timelimit '<' $x '1>' $output '2>&1'
       done	
    done 
done
