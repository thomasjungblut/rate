#!/bin/bash

for j in {1..100}; do
   for (( c=0; c<=5+j; c++ )); do
     echo "line ${j}"
   done
   sleep 1
done

