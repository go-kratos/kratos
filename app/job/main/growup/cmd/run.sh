dates=("2017-06-01" "2017-07-01" "2017-08-01" "2017-09-01" "2017-10-01" "2017-11-01" "2017-12-01" "2018-01-01" "2018-02-01" "2018-03-01" "2018-04-01" "2018-05-01" "2018-06-01" "2018-07-01" "2018-08-01" "2018-09-01" "2018-10-01")

for date in ${dates[@]}
do
	echo "start del $date"
	curl -XPOST localhost:8000/x/internal/job/up-rating/score/del?date=$date
	if [[ $? != 0 ]]; then
		exit
	fi
	echo ------------\n
done

for date in ${dates[@]}
do
	echo "start run $date"
	for ((i=0;i<=22;i++)); do
		echo "start run past $date", $i
		curl -XPOST localhost:8000/x/internal/job/up-rating/past/score?date=$date
		if [[ $? != 0 ]]; then
			exit
		fi
	done

	curl -XPOST localhost:8000/x/internal/job/up-rating/score?date=$date
	if [[ $? != 0 ]]; then
		exit
	fi
	echo ------------\n
done
