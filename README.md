# Confsgo

This is a C/S program, to deploy configure files to clients. The server will start a http listener listen at 60000 by default, and server configure files. Client continuely try to get configure file from the server 3 second one time. When configure file be got server will rename the configure to other name.    

##Example:

**server side run:**
``` confsgo.exe -type server -after="after.bat" -befor="befor.bat" -path="." ```


**client side run:**
``` confsgo.exe -type client -url="http://localhost:60000/config.txt" -after=after.bat ```
