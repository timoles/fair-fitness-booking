package main

import (
	"gopkg.in/headzoo/surf.v1"
	"fmt"
	"strings"
	"io/ioutil"
	"flag"
	"net/http"
	"crypto/tls"
	"os"
	"strconv"
	"log"
)

var wantedCourses map[string]string
var bow = surf.NewBrowser()
var debug = false
/*
   <div class="sclogin-joomla-login vertical span12">
                <form action="/fitness-clubs/neu-ulm" method="post" id="sclogin-form312">
                <div class="control-group" id="form-sclogin-username">
                            <input name="username" tabindex="0" id="sclogin-username" class="input-block-level" alt="username" type="text" placeholder="Benutzername" required="" aria-required="true"/>
                <div class="control-group" id="form-sclogin-password">
                            <input name="password" tabindex="0" id="sclogin-passwd" class="input-block-level" alt="password" type="password" placeholder="Passwort" required="" aria-required="true"/>

*/
func main() {
	//os.Setenv("HTTP_PROXY", "http://127.0.0.1:8080")
	// TODO decrease HTTP timeout
	// go run fairfitnes.go 'Yoga,10' 'TRX,19'
	// Ignore https certs
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	_, err := http.Get("https://golang.org/")
	checkErrorPanic(err)

	fmt.Println("TODO USAGE")
	// TODO Usage
	flag.Parse()
	// Read from args what courses the user wants to book
	wantedCourses = make(map[string]string)
	for i := 0; i < len(flag.Args()); i++ {
		splitted := strings.Split(flag.Args()[i], ",")
		if len(splitted) != 2 || splitted[0] == "" || splitted[1] == "" {
			log.Fatal("Wrong Parameters (Example: fairfitnes.go 'Yoga,10' 'TRX,19')")
		}
		wantedCourses[splitted[0]] = splitted[1]
	}
	// Open the booking site
	err = bow.Open("https://fair-fitness.com/fitness-clubs/neu-ulm")
	if err != nil {
		panic(err)
	}
	test := bow.Body()
	if !readKW(test) {
		os.Exit(1)
	}

	// Read password file and submit form with login credentials
	dat, err := ioutil.ReadFile("password.txt")
	checkErrorPanic(err)
	pw := string(dat)
	pw = strings.TrimSuffix(pw, "\n")
	fm, err := bow.Form("[id='sclogin-form312']")
	fm.Input("username", "Timo")
	fm.Input("password", pw)

	if fm.Submit() != nil {
		panic(err)
	}

	bow.Back()
	body := bow.Body()

	checkPage(body)
	fmt.Println("Done")
}

/*
Returns true if a new week is detected, false if not

IMPORTANT!!! New and Old are switched up
 */
func readKW(body string) bool {
	// Read kw file and compare to current kw
	dat, err := ioutil.ReadFile("kw.txt")
	checkErrorPanic(err)
	kwNew := string(dat)
	kwNew = strings.TrimSuffix(kwNew, "\n")
	kwNew = strings.TrimSpace(kwNew)

	// fmt.Println("New: ", kwNew)
	kwIndex := strings.Index(body, "Kursplan KW:")
	kwNewInt, err := strconv.Atoi(kwNew)

	checkErrorPanic(err)
	kwOld := body[kwIndex+12 : kwIndex+14]
	kwOld = strings.TrimSuffix(kwOld, "\n")
	kwOld = strings.TrimSpace(kwOld)
	// fmt.Println("Old: ", kwOld)
	kwOldInt, err := strconv.Atoi(kwOld)
	checkErrorPanic(err)
	if kwNewInt < kwOldInt || kwNewInt == 1 {
		fmt.Println("New Week")
		writeKW(kwOld)
		return true
	}
	if debug == true {
		return true
	}
	fmt.Println("No new course plan detected, exiting...")
	return false
}


func writeKW(kw string) {
	f, err := os.Create("kw.txt")
	checkErrorPanic(err)
	defer f.Close()
	_, err = f.WriteString(kw) // doesn't matter if KW is only one digit because there's a trailing space
	checkErrorPanic(err)
}


func checkPage(body string) {
	// fmt.Println(body)
	coursesAvailable := strings.Contains(body, "cpf_tr3")
	// Check if a new week is available, TODO doesn't recheck if booking is possible in the middle of the week


	for coursesAvailable {
		// Get Course
		indexCourseBegin := strings.Index(body, "cpf_tr3")
		indexCourseEnd := strings.Index(body[indexCourseBegin:], "</tr>")
		courseString := body[indexCourseBegin : indexCourseBegin+indexCourseEnd]
		// fmt.Println(courseString)
		timeIndexStart := strings.Index(courseString, "<td class=\"cpf_td4a\">")
		timeIndexEnd := strings.Index(courseString[timeIndexStart+21:], ":")
		timeStart := courseString[timeIndexStart+21 : timeIndexStart+21+timeIndexEnd]
		timeStart = strings.TrimSpace(timeStart)
		timeStartInt, err := strconv.Atoi(timeStart)
		checkErrorPanic(err)
		// Get course Type
		indexTypeBegin := strings.Index(courseString, "cpf_sp1a") + 10
		indexTypeEnd := strings.Index(courseString[indexTypeBegin:], "</span")
		courseType := courseString[indexTypeBegin : indexTypeBegin+indexTypeEnd]
		fmt.Println("Coursetype found: " + courseType)

		// Check if check-in possible and extract submit value
		if strings.Contains(courseString, "Reservieren!") {
			valueIndexBegin := strings.Index(courseString, "value=") + 7
			valueIndexEnd := strings.Index(courseString[valueIndexBegin:], "\"")
			submitValue := courseString[valueIndexBegin : valueIndexBegin+valueIndexEnd]
			submitValue = strings.TrimSpace(submitValue)
			fmt.Println("	[+] Course can be booked with value: " + submitValue)
			timeArgs, bookingRequested := wantedCourses[courseType]

			if bookingRequested {
				// Only if booking is requested we have a valid timeArgs
				timeArgsInt, err := strconv.Atoi(timeArgs)
				checkErrorPanic(err)
				if timeStartInt >= timeArgsInt {
					fmt.Println("	[*] Trying to book course")
					bookCourse(submitValue)
				}
			}
		} else {
			fmt.Println("	[-] Course full")
		}
		body = body[indexCourseBegin+indexCourseEnd:]
		coursesAvailable = strings.Contains(body, "cpf_tr3")
	}

}

func bookCourse(value string) {
	// Send Post Booking request to course with {value}
	// cp_res=%C2%A0%C2%A0%C2%A0%C2%A0%C2%A03259*7
	postBody := "cp_res=%C2%A0%C2%A0%C2%A0%C2%A0%C2%A0" + value
	bow.Post("https://fair-fitness.com/fitness-clubs/neu-ulm", "application/x-www-form-urlencoded", strings.NewReader(postBody))

	fmt.Println("Tried to book, mby worked TODO confirmation") // TODO

}

func checkErrorPanic(err error) {
	if err != nil {
		f, _ := os.Create("error.txt")
		defer f.Close()
		log.SetOutput(f) // TODO error file
		panic(err)

	}
}
