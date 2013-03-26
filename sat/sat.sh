#!/bin/bash

# run.sh 
# 1 file default test.lp
# 2 solver
# 3 timeout

instance=${1-test.cnf}
solver=${2-1}
timeout=${3-60}

case $solver in
    1) timeout $timeout minisat $instance  ;;
    2) timeout $timeout cryptominisat $instance  ;;
    3) timeout $timeout lingeling --verbose $instance  ;;
    4) clasp --time-limit=$timeout -q --stats --configuration=crafty $instance  ;;
    5) timeout $timeout glucose_static $instance | grep '^[^v]' ;;
esac
