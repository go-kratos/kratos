#!/bin/bash

for file in `ls *.csv`  
do  
	awk -f tst.awk $file > t_$file
done  
