# Goal of the code
This simple GO app ( my first one in GO ) is aimed to seamlessly ( without user intervention ) update the hosts file on Windows / Linux / MacOS / RaspberryPi computers in a small company or in advanced home network.

# Usage
Just build it for your platform with:
```go build```
and than just run it or schedule to run it periodically ( ie. with *crontab* ).
Keep in mind, that you need to run it with higher privileges ( ie. sudo or RunAs Adm ).

# Config
The config needed by this app is provided by simple package named *psikuta*, which was ( for obvious reasons ) git ignored for this repo.

# ToDo
Not yet tested on MacOS and ARM platforms ....

# Licence
This code is released under BSD licence. 
