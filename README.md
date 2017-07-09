# crongo

## Dead simple cron job monitoring


* No dependencies
* Easy to deploy static binaries
* Written in Go

# Usage

TODO


# How it works


/usr/bin/crongo -cmd /home/myuser/myjob.sh

crongo executes the command, and does an HTTP POST request to the crogod server specified in the config file 

crongod saves a copy of this job on the filesystem:
  YYYYMMDDT_HH_SS_hostname_command.json
