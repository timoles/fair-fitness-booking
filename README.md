# FairFitness Booking

Tool written in golang to automatically book choosen courses on https://fair-fitness.com/fitness-clubs/neu-ulm

## Build

go get gopkg.in/headzoo/surf.v1
go build fairfitness.go

## Example Usage

./book 'TRX' 'Yoga'

## TODO

-  Currently trying to book courses that are already booked, need to change that

- !!! Timeout in HTTP Client !!!

- Properly check if courses are in the future

- Make it possible to set a time from which on courses will be booked e.g. 1600 (can do this in parsing of future() )

- Handle invalid certs a better way

- Improve parsing

- Improve login with form submit

- Properly check if the course booking worked

- Don't book courses right before they start (to avoid reserving not used seats)

- Better logging

- Usage Messages
