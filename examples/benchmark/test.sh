#!/bin/bash

server_bin_name="gowebbenchmark"

. ./libs.sh

length=${#web_frameworks[@]}

test_result=()

cpu_cores=`cat /proc/cpuinfo|grep processor|wc -l`
if [ $cpu_cores -eq 0 ]
then
  cpu_cores=1
fi

test_web_framework()
{
  echo "testing web framework: $2"
  ./$server_bin_name $2 $3 &
  sleep 2

  throughput=`wrk -t$cpu_cores -c$4 -d30s http://127.0.0.1:8080/ | grep Requests/sec | awk '{print $2}'`
  echo "throughput: $throughput requests/second"
  test_result[$1]=$throughput

  pkill -9 $server_bin_name
  sleep 2
  echo "finsihed testing $2"
  echo
}

test_all()
{
  echo "###################################"
  echo "                                   "
  echo "      ProcessingTime  $1ms         "
  echo "      Concurrency     $2           "
  echo "                                   "
  echo "###################################"
  for ((i=0; i<$length; i++))
  do
  	test_web_framework $i ${web_frameworks[$i]} $1 $2
  done
}


pkill -9 $server_bin_name

echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > processtime.csv
test_all 0 5000
echo "0 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
test_all 10 5000
echo "10 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
test_all 100 5000
echo "100 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv
test_all 500 5000
echo "500 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime.csv


echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > concurrency.csv
test_all 30 100
echo "100,"$(IFS=$','; echo "${test_result[*]}" ) >> concurrency.csv
test_all 30 1000
echo "1000,"$(IFS=$','; echo "${test_result[*]}" ) >> concurrency.csv
test_all 30 5000
echo "5000,"$(IFS=$','; echo "${test_result[*]}" ) >> concurrency.csv


test_all -1 5000
echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > cpubound.csv
echo "cpu-bound,"$(IFS=$','; echo "${test_result[*]}" ) >> cpubound.csv

echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > cpubound-concurrency.csv
test_all -1 100
echo "100,"$(IFS=$','; echo "${test_result[*]}" ) >> cpubound-concurrency.csv
test_all -1 1000
echo "1000,"$(IFS=$','; echo "${test_result[*]}" ) >> cpubound-concurrency.csv
test_all -1 5000
echo "5000,"$(IFS=$','; echo "${test_result[*]}" ) >> cpubound-concurrency.csv


mv -f processtime.csv ./testresults
mv -f concurrency.csv ./testresults
mv -f cpubound.csv ./testresults
mv -f cpubound-concurrency.csv ./testresults
