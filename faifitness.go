package main

import (
	"gopkg.in/headzoo/surf.v1"
	"fmt"
	"strings"
	"time"
	"flag"
	"net/http"
	"net/url"
	"bytes"
	"io/ioutil"
)
var wantedCourses map[string]struct{}
/*
   <div class="sclogin-joomla-login vertical span12">
                <form action="/fitness-clubs/neu-ulm" method="post" id="sclogin-form312">
                <div class="control-group" id="form-sclogin-username">
                            <input name="username" tabindex="0" id="sclogin-username" class="input-block-level" alt="username" type="text" placeholder="Benutzername" required="" aria-required="true"/>
                <div class="control-group" id="form-sclogin-password">
                            <input name="password" tabindex="0" id="sclogin-passwd" class="input-block-level" alt="password" type="password" placeholder="Passwort" required="" aria-required="true"/>

*/
func main() {
	flag.Parse()
	//fmt.Println("tail", flag.Args())
	wantedCourses = make(map[string]struct{})
	for i:=0; i< len(flag.Args()) ; i++{
		wantedCourses[flag.Args()[i]] = struct{}{}
	}
	
	
	bow := surf.NewBrowser()
	err := bow.Open("https://fair-fitness.com/fitness-clubs/neu-ulm")
	if err != nil {
		panic(err)
	}

	// Outputs: "The Go Programming Language"
	
	fm, err := bow.Form("[id='sclogin-form312']")
	// fmt.Println(err,fm)
	fm.Input("username", "username")
	fm.Input("password","password")
    if fm.Submit() != nil {
    	panic(err)
    }
    
    body := bow.Body()
	//fmt.Println(bow.Body)
	/*
    fm.Input("sclogin-username", "JoeRedditorW")
    fm.Input("sclogin-passwd", "d234rlkasd")
    if fm.Submit() != nil {
    	panic(err)
    }
    */
    fmt.Println("loggedin")

    fmt.Println(body)

	test := strings.Index(body,"Dienstag")
	if test < 0 {
		return
	}
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
	fmt.Println(body)
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
	form := url.Values{
		"cp_res": {"				" + value},
	}
	return
	// func Post(url string, bodyType string, body io.Reader) (resp *Response, err error) {
	//对form进行编码
	body := bytes.NewBufferString(form.Encode())
	resp, err := http.Post("https://fair-fitness.com/fitness-clubs/neu-ulm", "application/x-www-form-urlencoded", body)
		if err != nil {
		panic(err)
	}
	if resp.StatusCode == 200{
		bodyBytes, _ :=ioutil.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		fmt.Println("	[+] Course booked successfully")
	}
}
func past(){

}

// cp_res=%C2%A0%%C2%A0%%C2%A0%%C2%A0%%C2%A0%
