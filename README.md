# FairFitness Booking

Tool written in golang to automatically book choosen courses on https://fair-fitness.com/fitness-clubs/neu-ulm

## Build

go get gopkg.in/headzoo/surf.v1
go build fairfitness.go

## Example Usage

Set your password in the password.txt

Change the username withing fairfitness.go in the line (fm.Input("username", "[Your Username]"))

./book 'TRX' 'Yoga'

## Troubleshoot

* No booking confirmation

        Check your login credentials