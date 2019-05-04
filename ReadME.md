# CSV > JSON converter
The application converts csv files to json formatted files.

CSV files inside the given directory are being converted into JSON formatted files. I decided to declare a cron job for observing the directory every minute for new CSV files. If there is a unconverted CSV file inside the directory, then the conversion will take place.

At first, I decided to create a file converter on server which receives csv files from a message queue or by file upload. Then the each minute file conversion will take place.

It is not a production scale application obviously. My intention is to learn and experiment something on my own. If anyone has an idea,the contributions and feedbacks are welcomed.

# How to
CSV files should be added to the given directory. The folder path should be passed as a flag parameter as given below
```
$ go build -o converter.exe

$ converter.exe -folder=C:\Users\ekin\Desktop\test\ -filetype=csv -targetType=json
```

For now the applications tracks the folder each minute for csv files. After that it converts csv files to json formatted ones.