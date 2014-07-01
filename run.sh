#!/bin/bash
source="$2"
if [[ "$1" = "python" ]]; then
    cmd="./crawl.py"
else
    cmd="./crawl" 
fi

for i in 1,10 3,10 3,100 5,100 10,100 4,1000 10,1000
do
    IFS=","
    set $i
    run="time $cmd $source $1 $2"
    echo $run
    echo $(time $cmd $source $1 $2)
done
