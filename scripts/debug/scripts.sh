#!/bin/bash
for i in {12001..12065};
do
	echo $i
	sed -n '/#/,/#/p' Node_127.0.0.1\:$i.graph > nodes_$i.txt
	sort -k2,2n -k3,3n nodes_$i.txt > nodes_$i.sorted
	rm nodes_$i.txt
done


