#!/bin/bash

echo "Close ALL apps before running these benchmarks !"
echo " "
echo "AND make sure your MacBook is externally powered !"
echo " "
echo "thkis may take a while to run ..."

mkdir -p b

TEST_COUNT=10

# BenchmarkGetTopicPrivateHandler

# NOTE: benchmarks are run in their own loop to ensure use of cache for most repeatable measurements
x=$TEST_COUNT
while [ $x -gt 0 ];
do
  go test -run=api/topic_test.go -bench=GetTopicPrivateHandler >b/t-pri-$x.txt
  x=$(($x-1))
done

# get the results from each benchtest & delete
cat b/t-pri-1.txt | grep "allocs/op" >b/t-pri.txt
rm b/t-pri-1.txt
x=$TEST_COUNT
while [ $x -gt 1 ];
do
  cat b/t-pri-$x.txt | grep "allocs/op" >>b/t-pri.txt
  rm b/t-pri-$x.txt
  x=$(($x-1))
done

echo "b/t-pri :"
cat b/t-pri.txt


# BenchmarkGetDatasetPrivate

# NOTE: benchmarks are run in their own loop to ensure use of cache for most repeatable measurements
x=$TEST_COUNT
while [ $x -gt 0 ];
do
  go test -run=api/topic_test.go -bench=GetDatasetPrivate >b/d-pri-$x.txt
  x=$(($x-1))
done

# get the results from each benchtest & delete
cat b/d-pri-1.txt | grep "allocs/op" >b/d-pri.txt
rm b/d-pri-1.txt
x=$TEST_COUNT
while [ $x -gt 1 ];
do
  cat b/d-pri-$x.txt | grep "allocs/op" >>b/d-pri.txt
  rm b/d-pri-$x.txt
  x=$(($x-1))
done

echo "b/d-pri :"
cat b/d-pri.txt


# BenchmarkGetTopicPublicHandler

# NOTE: benchmarks are run in their own loop to ensure use of cache for most repeatable measurements
x=$TEST_COUNT
while [ $x -gt 0 ];
do
  go test -run=api/topic_test.go -bench=GetTopicPublicHandler >b/t-pub-$x.txt
  x=$(($x-1))
done

# get the results from each benchtest & delete
cat b/t-pub-1.txt | grep "allocs/op" >b/t-pub.txt
rm b/t-pub-1.txt
x=$TEST_COUNT
while [ $x -gt 1 ];
do
  cat b/t-pub-$x.txt | grep "allocs/op" >>b/t-pub.txt
  rm b/t-pub-$x.txt
  x=$(($x-1))
done

echo "t-pub :"
cat b/t-pub.txt


# BenchmarkGetDatasetPublic

# NOTE: benchmarks are run in their own loop to ensure use of cache for most repeatable measurements
x=$TEST_COUNT
while [ $x -gt 0 ];
do
  go test -run=api/topic_test.go -bench=GetDatasetPublic >b/d-pub-$x.txt
  x=$(($x-1))
done

# get the results from each benchtest & delete
cat b/d-pub-1.txt | grep "allocs/op" >b/d-pub.txt
rm b/d-pub-1.txt
x=$TEST_COUNT
while [ $x -gt 1 ];
do
  cat b/d-pub-$x.txt | grep "allocs/op" >>b/d-pub.txt
  rm b/d-pub-$x.txt
  x=$(($x-1))
done

echo "d-pub :"
cat b/d-pub.txt


# compute and show averages:
echo " "
./b-total.sh
