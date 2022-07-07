# Easy Games
A twitch based stream recommendation server coded by Golang.

Click [here](http://13.59.49.252) to explore...

## Features
- User Log in/Log out/Register
  - Session-based Authentication and Authorization
- Search top games from Twitch backend server
- Add/Delete favorite items
  - in three types, VEDIO, CLIP and STEAM
- Get user's all favorite items
- Content based recommendation, recommend 
  - by user's history(favorite history)
  - by default(topgames) 

## Code brief diagram
![code diagram](https://user-images.githubusercontent.com/66594541/176646519-df0ec1d7-a00c-4006-a3e7-dd34423c03c7.jpg)

## Frontend structure
![Component Tree](https://user-images.githubusercontent.com/66594541/177716388-857b839d-32d1-45ae-a3c9-85a8ba3c8fb2.jpg)

## APIs - postman api collection link
https://www.getpostman.com/collections/bdaa61a62fad141adde4

## Notes:
- `netstat` command
```
netstat -ln
netstat -a
```
- screen - keep golang server long-runniung
```
screen -S session_name        # create session
screen -r session_name        # restore session
screen -ls                    # list the current running screen sessions
screen -wipe session_name     # delete the session
screen -X -S session_id quit  # kill the session
```
