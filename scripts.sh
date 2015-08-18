for i in $(seq 1 600); do echo "$i / 600" | bc -l 1>&2 ; curl -s 'localhost:8000/dir/data'; done > count.txt

cat count.txt | sort | uniq -c

