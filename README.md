# Assignment_2\_\_Security

This is the second assignement for security, made by Casper Storm Frøding.

## How to run

You need 4 terminals open in the Assignment_2\_\_Security directory. 3 of them will act as the patients and the last will act as the hospital.
These commands have been run in bash.
If you encounter problems, please run the command: go mod tidy

To start you run this command on one of the terminals, which will start the hospital:

```bash
go run hospital/hospital.go -port 5454
```

Then the patients can be started, by runnning these commands:
(OBS!!! Please be aware that the patients have a 10 second buffer to make sure that all patients have been initiliased at this time, so all three commands needs to be run in the span of AT MOST 10 seconds! )

```bash
go run patient/patient.go -id=1
go run patient/patient.go -id=2
go run patient/patient.go -id=3
```

## How to read the values

The private values the patients have, are a random number that is between 1 and 10 000, and can be seen in the first line printed, i.e. `2024/10/22 16:15:25 Patient 3 just started, with the value: ____`
At last the hospital prints the calculated sum of all the private values.
