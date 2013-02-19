#!/bin/bash

# run.sh 
# 1 file default test.lp
# 2 <model number 1-10> default 1
# 3 <timelimit in seconds> default 60
# 4 <option 0-5> default crafty
# 5 random number
# 6 output file

random=${5-1982}
output=${6-output.txt}

option='--time-limit='${3-1800}' -t 1 --stats '
option=$option'--trans-ext=all '
option=$option'--seed='$random' ' 
#quiet:
#option=$option'-q ' 

case ${4-3} in
    0) option=$option'--configuration=frumpy ';;
    1) option=$option'--configuration=jumpy ';;
    2) option=$option'--configuration=handy ';;
    3) option=$option'--configuration=crafty ';;
    4) option=$option'--configuration=trendy ';;
    5) option=$option'--configuration=chatty ';;
esac

instance=${1-test.lp}

cat $instance | grep '^%'
echo -e 'Instance\t:' $instance | tee -a  $output
echo -e 'Model\t\t:' model${2-7}.lp | tee -a $output
echo -e 'Option\t\t: ' $option | tee -a $output
gringo $instance model${2-7}.lp | clasp $option  | tee -a $output


if grep -q '^SAT' $output
then
    solution=/tmp/solution_$(basename $output .txt)_$RANDOM.pl
    prettyOutput=/tmp/pretty_$(basename $output .txt)_$RANDOM.pl
    rm -fr $solution
    rm -fr $prettyOutput
    cat $output | grep 'is(' |  tail -n 1 | sed 's/ /\n/g' | sed 's/$/./g' | sort  > $solution
    cat print.pl >> $solution
    cat $instance | sort >> $solution
    prolog -f print.pl -f $solution -g start -t halt > $prettyOutput
    column -t -s ',' $prettyOutput
fi
echo
