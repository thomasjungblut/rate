#!/bin/bash

for j in {1..100}; do
   rnx=$((1 + RANDOM % 10))
   for i in {1..10}; do
     echo "line"
   done
   sleep 1
done

