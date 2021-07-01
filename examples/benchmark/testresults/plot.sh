#!/bin/bash
m_path=$(dirname $0)
m_path=${m_path/\./$(pwd)}
cd $m_path

./transpose.sh

gnuplot benchmark.gnu

rm -fr t_*.csv
