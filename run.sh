#!/bin/bash

# run.sh 
# 1 file default test.lp
# 2 <model number 1-10> default 1
# 3 <timelimit in seconds> default 60
# 4 <option 0-5> default crafty
# 5 output file

Option='--time-limit='${3-1800}' -t 1 --stats --trans-ext=all '
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
echo -e 'Model\t\t:' model${2-7}.lp | tee -a $output
echo -e 'Option\t\t: ' $Option | tee -a $output
gringo $instance model${2-7}.lp | clasp $Option  | tee -a $output


if grep -q '^SAT' $output
then
    solution=/tmp/solution_$(basename $output .txt)_$RANDOM.pl
    prettyOutput=/tmp/pretty_$(basename $output .txt)_$RANDOM.pl
    rm -fr $solution
    rm -fr $prettyOutput
    cat $output | grep 'is_car' |  tail -n 1 | sed 's/ /\n/g' | sed 's/$/./g' | sort  > $solution
    cat print.pl >> $solution
    cat $instance | sort >> $solution
    prolog -f print.pl -f $solution -g start -t halt > $prettyOutput
    column -t -s ',' $prettyOutput
fi
echo
