#!/bin/bash

# run.sh 
# 1 file default test.lp
# 2 <model number 1-10> default 1
# 3 <timelimit in seconds> default 60
# 4 <option 0-5> default crafty
# 5 output file

Option='--time-limit='${3-60}' -t 1 --stats --trans-ext=all '
output=${5-output.txt}

case ${4-3} in
    0) Option=$Option'--configuration=frumpy ';;
    1) Option=$Option'--configuration=jumpy ';;
    2) Option=$Option'--configuration=handy ';;
    3) Option=$Option'--configuration=crafty ';;
    4) Option=$Option'--configuration=trendy ';;
    5) Option=$Option'--configuration=chatty ';;
esac

instance=${1-test.lp}

cat $instance | grep '^%'
echo -e 'Instance\t:' $instance | tee -a  $output
echo -e 'Model\t\t:' model${2-1}.lp | tee -a $output
echo -e 'Option\t\t: ' $Option | tee -a $output
gringo $instance model${2-1}.lp | clasp $Option  | tee -a $output
