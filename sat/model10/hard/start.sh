#!/bin/bash

timeout=100000

timeout $timeout lingeling pb_200_02.cnf --verbose | grep '^[^v]' >  pb_200_02.log &
timeout $timeout lingeling pb_200_06.cnf --verbose | grep '^[^v]' >  pb_200_06.log &
timeout $timeout lingeling pb_200_08.cnf --verbose | grep '^[^v]' >  pb_200_08.log &
timeout $timeout lingeling pb_300_02.cnf --verbose | grep '^[^v]' >  pb_300_02.log &
timeout $timeout lingeling pb_300_06.cnf --verbose | grep '^[^v]' >  pb_300_06.log &

sleep $timeout

timeout $timeout lingeling pb_300_09.cnf --verbose | grep '^[^v]' >  pb_300_09.log &
timeout $timeout lingeling pb_400_01.cnf --verbose | grep '^[^v]' >  pb_400_01.log &
timeout $timeout lingeling pb_400_02.cnf --verbose | grep '^[^v]' >  pb_400_02.log &
timeout $timeout lingeling pb_400_07.cnf --verbose | grep '^[^v]' >  pb_400_07.log &
timeout $timeout lingeling pb_400_08.cnf --verbose | grep '^[^v]' >  pb_400_08.log &













