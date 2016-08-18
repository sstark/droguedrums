#!/bin/bash
if [[ $1 == "fix" ]]
then
    diff="-w"
else
    diff="-d"
fi
gofmt $diff src/*.go
