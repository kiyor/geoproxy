#!/bin/bash
############################################

# File Name : pack.sh

# Purpose :

# Creation Date : 05-19-2017

# Last Modified : Fri 19 May 2017 11:46:53 PM UTC

# Created By : Kiyor 

############################################

go build
tar cvzf geoproxy.tar.gz \
	GeoLite2-City.mmdb \
	pac.tpl \
	conf \
	geoproxy
stf
