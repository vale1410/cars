#!/bin/zsh

instances=$1
basefolder=exp-$2
timelimit=$3

mkdir -p $basefolder

#encoding=('minisat' 'microsat' 'glucose' 'lingeling' 'cmsat' 'clasp')
#solver=('glucose' 'lingeling' 'cmsat' 'clasp')

encoding=('e1' 'e2' 'e3')
sol=glucose

case $4 in
    1) seed=(142) ;;
    2) seed=(142 321) ;;
    3) seed=(142 321 832) ;;
    4) seed=(142 321 832 999) ;;
esac 

for se in $seed
do 
    for enc in $encoding
    do 
       folder=$basefolder/$enc
       mkdir -p $folder
       for x in $instances/$enc/*
       do 
           output=$folder/$(basename $x .txt)-$se.log
           echo msat -seed $se -solver $sol -time $timelimit '<' $x '1>' $output '2>&1'
       done	
    done 
done
