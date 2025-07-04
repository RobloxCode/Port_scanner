In this project, you give an ip and then a range of ports that you want to scan, the ports can only go form 0 to 65535 though,
the program will chech to which of those ports is the ip connected to and will display the results in a file named with the 
current date inside a directory called ScanResults, if the file is not created the program will create it automatically.

To scan your ip, in the main function you just have to modify the variables: 'ip', and 'protocol' to yours in the main function,
you can also use the functions DisplayClosedPorts()/DisplayOpenPorts() to display in the console only the ports that your ip is connected or not to
respctively.
In the function ScanRangePortsConcurrently() you pass the first port and the last port as the 4th and 5th parameters respctively, the functions is not inclusive on the last port,
so if you want to include the last port you must add one to that number
