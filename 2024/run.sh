#!/bin/bash

set -e

total=0
# Dir only at depth 1
for dir in $(find . -maxdepth 1 -type d | sort); do
    if [ "$dir" != "." ]; then
        echo "$dir"
        cd "$dir" || exit
        go build -o sol .
        ./sol full > tmp.txt 2>&1
        exec_time_raw=$(grep -e "main in seconds" tmp.txt | cut -d ' ' -f 6)
        exec_time_ms=$(echo "$exec_time_raw * 1000" | bc)
        total=$(echo "$total + $exec_time_ms" | bc)
        echo "Execution time: ${exec_time_ms::-3} ms"
        rm sol
        rm tmp.txt
        echo "---------------------------------"
        cd ..
    fi
done
echo "Total execution time: ${total::-3} ms"
