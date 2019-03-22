#!/bin/bash

find ${pwd} -name "*\.go" -exec gofmt -l -w {} \;
