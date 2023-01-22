#!/bin/bash

#rm -rf /tmp/home1
#mkdir -p "/tmp/home1"
gok -hm /tmp/home1 -u testov signin
gok -hm /tmp/home1 -u testov login
gok -hm /tmp/home1 -u testov init
echo "create first record from home1"
gok -hm /tmp/home1 -u testov desc add First record
gok -hm /tmp/home1 -u testov desc ls


#rm -rf /tmp/home2
#mkdir -p "/tmp/home2"
gok -hm /tmp/home2 -u testov login
gok -hm /tmp/home2 -u testov init
gok -hm /tmp/home2 -u testov desc ls
echo "create second record from home2"
gok -hm /tmp/home2 -u testov desc add Second record
gok -hm /tmp/home2 -u testov desc ls

