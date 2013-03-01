#!/bin/bash

timeout=100000

clasp pb_200_02.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_200_02_clasp.log &
clasp pb_200_06.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_200_06_clasp.log &
clasp pb_200_08.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_200_08_clasp.log &
clasp pb_300_02.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_300_02_clasp.log &
clasp pb_300_06.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_300_06_clasp.log &
                                                            
sleep $timeout                                              
                                                            
clasp pb_300_09.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_300_09_clasp.log &
clasp pb_400_01.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_400_01_clasp.log &
clasp pb_400_02.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_400_02_clasp.log &
clasp pb_400_07.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_400_07_clasp.log &
clasp pb_400_08.cnf --time-limit=$timeout --configuration=crafty --stats >  pb_400_08_clasp.log &













