#!/bin/sh
#
i3status | while : 
do

	read line
	emails=$(./check-emails ./email.yml 2> /dev/null)
	if [ $emails -eq 0 ]; then
		echo $line || exit 1
	else
		echo "GMAIL: Unread ("$emails") | $line" || exit 1
	fi

done
