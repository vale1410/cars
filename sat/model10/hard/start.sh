#!/bin/bash

timeout=100000
solver=5

../../sat.sh pb_200_02.cnf $solver $timeout >  pb_200_02_crypto.log &
../../sat.sh pb_200_06.cnf $solver $timeout >  pb_200_06_crypto.log &
../../sat.sh pb_200_08.cnf $solver $timeout >  pb_200_08_crypto.log &
../../sat.sh pb_300_02.cnf $solver $timeout >  pb_300_02_crypto.log &
../../sat.sh pb_300_06.cnf $solver $timeout >  pb_300_06_crypto.log &

sleep $timeout                                              

../../sat.sh pb_300_09.cnf $solver $timeout >  pb_300_09_crypto.log &
../../sat.sh pb_400_01.cnf $solver $timeout >  pb_400_01_crypto.log &
../../sat.sh pb_400_02.cnf $solver $timeout >  pb_400_02_crypto.log &
../../sat.sh pb_400_07.cnf $solver $timeout >  pb_400_07_crypto.log &
../../sat.sh pb_400_08.cnf $solver $timeout >  pb_400_08_crypto.log &
