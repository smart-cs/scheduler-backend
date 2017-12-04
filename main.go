package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const port = ":8080"

var helper = NewCourseHelper()

func main() {
	http.HandleFunc("/schedules", schedulesHandler)
	http.HandleFunc("/courses", coursesHandler)

	fmt.Printf("Listening on port %s\n", port)

	http.ListenAndServe(port, nil)
}

func schedulesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	fmt.Printf("QUERY: %v\n", r.URL.Query())
	coursesStr := r.URL.Query().Get("courses")
	if coursesStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Requires "courses" query parameter`))
		return
	}
	courses := strings.Split(coursesStr, ",")
	for i, c := range courses {
		fmt.Printf("%d: %s\n", i, c)
	}
	data := helper.CreateSchedules(courses)
	j, _ := json.MarshalIndent(data, "", "    ")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func coursesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	data := helper.ValidCourses()
	j, _ := json.MarshalIndent(data, "", "    ")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
