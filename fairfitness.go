package main

import (
	"gopkg.in/headzoo/surf.v1"
	"fmt"
	"strings"
	"io/ioutil"
	"flag"
	"net/http"
	"crypto/tls"
)

var wantedCourses map[string]struct{}
var bow = surf.NewBrowser()
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
	// go run fairfitnes.go 'Yoga' 'TRX'
	// Ignore https certs
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    _, err := http.Get("https://golang.org/")
    if err != nil {
        fmt.Println(err)
    }

	fmt.Println("TODO USAGE")
	// TODO Usage
	flag.Parse()
	//fmt.Println("tail", flag.Args())
	wantedCourses = make(map[string]struct{})
	for i:=0; i< len(flag.Args()) ; i++{
		wantedCourses[flag.Args()[i]] = struct{}{}
	}
	
	
	
	err = bow.Open("https://fair-fitness.com/fitness-clubs/neu-ulm")
	if err != nil {
		panic(err)
	}

	// Read password file and submit form with login credentials
	dat, err := ioutil.ReadFile("password.txt")
	if err != nil {
		panic(err)
	}
	pw := string(dat)
	pw = strings.TrimSuffix(pw,"\n")
	fm, err := bow.Form("[id='sclogin-form312']")
	fm.Input("username", "Timo")
	fm.Input("password",pw)

    if fm.Submit() != nil {
    	panic(err)
    }

    bow.Back()
    body := bow.Body()

    fmt.Println("loggedin")

 
 /*
    // TODO this could 
    // be used for further checks for smarter tool to not create to much 
    // reserve noise

	// Check if Friday is available to reserve  
	test := strings.Index(body,"Freitag")
	if test < 0 {
			return
		}
		
	}
	// is Friday in future?
	test2 := body[test:test+20]
	parts := strings.Split(test2, ",")
	test3 := strings.TrimSpace(parts[1])

	// 24.07.2018
	// t, err := time.Parse("2006-01-02", "2011-01-19")
	t, err := time.Parse("02.01.2006", test3)

	if err != nil{
		fmt.Println(err)
	}
	until := time.Until(t)
*/

	until := 10 
	if until > 0{
		fmt.Println("in the future")
		future(body)
	}else{
		fmt.Println("No Future Sleep TODO")
		past()
	}
	
	fmt.Println("Done")

}
func future(body string){
	// fmt.Println(body)
	coursesAvailable := strings.Contains(body,"cpf_tr3")
	for coursesAvailable {
		// Get Course
		indexCourseBegin := strings.Index(body,"cpf_tr3")
		indexCourseEnd := strings.Index(body[indexCourseBegin:],"</tr>") 
		courseString := body[indexCourseBegin:indexCourseBegin+indexCourseEnd]

		// Get course Type
		indexTypeBegin := strings.Index(courseString,"cpf_sp1a")+10
		indexTypeEnd := strings.Index(courseString[indexTypeBegin:],"</span")
		courseType := courseString[indexTypeBegin:indexTypeBegin+indexTypeEnd]
		fmt.Println("Coursetype found: " + courseType)

		// Check if check-in possible and extract submit value
		if strings.Contains(courseString,"Reservieren!"){
			valueIndexBegin := strings.Index(courseString,"value=")+7
			valueIndexEnd := strings.Index(courseString[valueIndexBegin:],"\"")
			submitValue := courseString[valueIndexBegin:valueIndexBegin+valueIndexEnd]
			submitValue = strings.TrimSpace(submitValue)
			fmt.Println("	[+] Course can be booked with value: " + submitValue)
			_,bookingRequested := wantedCourses[courseType]
			if bookingRequested{
				fmt.Println("	[*] Trying to book course")
				bookCourse(submitValue)
			}
		}else{
			fmt.Println("	[-] Course full")
		}
		body = body[indexCourseBegin+indexCourseEnd:]
		coursesAvailable = strings.Contains(body,"cpf_tr3")
			
	}
}

func bookCourse(value string){
	// Send Post Booking request to course with {value}
	// cp_res=%C2%A0%C2%A0%C2%A0%C2%A0%C2%A03259*7
	postBody := "cp_res=%C2%A0%C2%A0%C2%A0%C2%A0%C2%A0" + value
	bow.Post("https://fair-fitness.com/fitness-clubs/neu-ulm","application/x-www-form-urlencoded",strings.NewReader(postBody))

    fmt.Println("Tried to book, mby worked TODO confirmation") // TODO

}

func past(){

}
