#!/bin/bash
############################################

# File Name : start.sh

# Purpose :

# Creation Date : 05-15-2017

# Last Modified : Mon May 15 11:25:17 2017

# Created By : Kiyor 

############################################

killall geoproxy
sleep 1
nohup ./geoproxy </dev/null >geoproxy.log 2>&1&
