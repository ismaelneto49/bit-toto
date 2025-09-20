## Run application

This benchmark create a node at 5000 port that will be the requester  
The remaining nodes will act as servers

The root files of the projects are in `root-files`  
The files available per node are in `files/NODE`

Run this command:  
`./init-benchmark.sh`

Then type `request file1.txt`

Clean TCP ports  
`sudo kill -9 $(sudo lsof -t -i:{5000..5050})`
