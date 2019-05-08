# CSV > JSON converter
This small converter helps us to convert CSV files to JSON file. Later on I am planning to make a general converter which supports different file formats

CSV files in the given directory are being converted into JSON formatted files. Each second application tracks the directory for new CSV files. So once the application is up, new CSV files can be added under directory.

# How to
CSV files should be added to the given directory. The folder path should be passed as a flag parameter as given below
```
$ go build -o converter.exe

$ converter.exe -folder=C:\Users\ekin\Desktop\test\ -filetype=csv -targetType=json
```