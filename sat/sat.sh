#!/bin/bash

# run.sh 
# 1 file default test.lp
# 2 solver
# 3 timeout

instance=${1-test.cnf}
solver=${2-1}
timeout=${3-60}

case $solver in
    1) minisat $instance -cpu-lim=$timeout  ;;
    2) timeout $timeout cryptominisat $instance --nosolprint  ;;
    3) timeout $timeout lingeling --verbose $instance  ;;
    4) clasp --time-limit=$timeout -q --stats --configuration=handy $instance  ;;
    5) timeout $timeout glucose_static $instance | grep '^[^v]' ;;
    6) clasp --time-limit=$timeout -q --stats --configuration=jumpy $instance  ;;
    7) clasp --time-limit=$timeout -q --stats --configuration=crafty $instance  ;;
    8) clasp --time-limit=$timeout -q --stats --configuration=trendy $instance  ;;
    9) clasp --time-limit=$timeout -q --stats --configuration=frumpy $instance  ;;
esac
