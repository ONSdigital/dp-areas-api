#!/bin/bash

count=0;
t_pri_ops_per_second=0; 

for i in $( awk '{ print $1; }' b/t-pri.txt )
   do 
     t_pri_ops_per_second=$(echo $t_pri_ops_per_second+$i | bc )
     ((count++))
   done

echo "Topic Private Average Ops per second:"
echo "scale=2; $t_pri_ops_per_second / $count" | bc


count=0;
d_pri_ops_per_second=0; 

for i in $( awk '{ print $1; }' b/d-pri.txt )
   do 
     d_pri_ops_per_second=$(echo $d_pri_ops_per_second+$i | bc )
     ((count++))
   done

echo " "
echo "Dataset Private Average Ops per second:"
echo "scale=2; $d_pri_ops_per_second / $count" | bc

echo " "
echo "Topic Private improvement :"
echo "scale=4; $t_pri_ops_per_second / $d_pri_ops_per_second" | bc

######

count=0;
t_pub_ops_per_second=0; 

for i in $( awk '{ print $1; }' b/t-pub.txt )
   do 
     t_pub_ops_per_second=$(echo $t_pub_ops_per_second+$i | bc )
     ((count++))
   done

echo " "
echo "#####"
echo " "

echo "Topic Public Average Ops per second:"
echo "scale=2; $t_pub_ops_per_second / $count" | bc


count=0;
d_pub_ops_per_second=0; 

for i in $( awk '{ print $1; }' b/d-pub.txt )
   do 
     d_pub_ops_per_second=$(echo $d_pub_ops_per_second+$i | bc )
     ((count++))
   done

echo " "
echo "Dataset Public Average Ops per second:"
echo "scale=2; $d_pub_ops_per_second / $count" | bc

echo " "
echo "Topic Public improvement :"
echo "scale=4; $t_pub_ops_per_second / $d_pub_ops_per_second" | bc

# awk '{ total += $1 } END { print total/NR }' b/d-pub.txt