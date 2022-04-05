package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Appointment struct {
	EndsAt string `json:"ends_at"`
	Id int `json:"id"`
	UserId int `json:"user_id"`
	StartsAt string `json:"starts_at"`
	TrainerId int `json:"trainer_id"`
}

var Appointments []Appointment

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllAppointments(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllAppointments")
	json.NewEncoder(w).Encode(Appointments)
}

// returns a single appointment based on the appointment Id
func returnSingleAppointmentBasedOnId(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["id"]
	intKey, _ := strconv.Atoi(key)
	fmt.Println("id key: ", key)

	for _, appointment := range Appointments {
		if appointment.Id == intKey {
			json.NewEncoder(w).Encode(appointment)
		}
	}

	fmt.Fprintf(w, "Id: " + key)
}

// returns all the appointments for a particular Trainer Id
func returnAppointmentsForTrainer(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["trainer_id"]
	intKey, _ := strconv.Atoi(key)
	fmt.Println("trainer id key: ", key)

	for _, appointment := range Appointments {
		if appointment.TrainerId == intKey {
			json.NewEncoder(w).Encode(appointment)
		}
	}

	fmt.Fprintf(w, "Trainer Id: " + key)
}

// returns all available appointments for a particular trainer between two dates
func returnAvailableAppointmentsForTrainer(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["trainer_id"]

	startTimeString := r.URL.Query().Get("starts_at")
	endTimeString := r.URL.Query().Get("ends_at")
	fmt.Println(startTimeString)
	fmt.Println(endTimeString)

	intKey, _ := strconv.Atoi(key)
	fmt.Println("Current trainer id: ", key)

	//const layout = "2019-01-24T09:30:00-08:00"
	startTime, _ := time.Parse(time.RFC3339, startTimeString)
	endTime, _ := time.Parse(time.RFC3339, endTimeString)

	// first create an array of all times between those two time ranges
	var allTimes []time.Time
	allTimes = append(allTimes, startTime)
	nextTime := startTime.Add(time.Minute * 30)

	for !nextTime.Equal(endTime) {
		allTimes = append(allTimes, nextTime)
		nextTime = nextTime.Add(time.Minute * 30)
	}

	//create an array containing the times available for appointments
	var availableAppts []time.Time

	for _, curApptTime := range allTimes { // iterate through all possible times
		fmt.Println(curApptTime)
		currTimeTaken := false
		for _, appointment := range Appointments { // iterate through all scheduled times to see if there's a match (if there is, don't add)
			if appointment.TrainerId == intKey { // if the trainer Id matches, check to see if start time matches
				takenApptTime, _ := time.Parse(time.RFC3339, appointment.StartsAt)
				if curApptTime.Equal(takenApptTime) {
					currTimeTaken = true
				}
			}
		}
		if !currTimeTaken {
			availableAppts = append(availableAppts, curApptTime)
			json.NewEncoder(w).Encode(curApptTime)
		}
	}

	fmt.Println("Trainer ", key," has the following times available for appointments between ", startTime, " and ", endTime)
	for _, time := range availableAppts {
		fmt.Println(time)
	}
}

func createNewAppointment(w http.ResponseWriter, r *http.Request) {
	endsAt := r.URL.Query().Get("ends_at")
	fmt.Println(endsAt)
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	fmt.Println("id: ", id)
	userId, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	fmt.Println("user id: ", userId)
	startsAt := r.URL.Query().Get("starts_at")
	fmt.Println(startsAt)
	trainerId, _ := strconv.Atoi(r.URL.Query().Get("trainer_id"))
	fmt.Println("trainer id: ", trainerId)

	newAppointment := Appointment{
		EndsAt: endsAt,
		Id: id,
		UserId: userId,
		StartsAt: startsAt,
		TrainerId: trainerId,
	}

	json.NewEncoder(w).Encode(newAppointment)
	Appointments = append(Appointments, newAppointment)

	// print all the appointments
	fmt.Println("-----")
	for _, appt := range Appointments {
		fmt.Println(appt)
	}

	// write appointments to file
	writeAppointmentsToFile()
}

func handleRequests() {
	// creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homePage)
	router.HandleFunc("/all", returnAllAppointments)
	router.HandleFunc("/appointment", createNewAppointment).Methods("POST")
	router.HandleFunc("/appointment/{id}", returnSingleAppointmentBasedOnId)
	router.HandleFunc("/appointment/trainer/{trainer_id}", returnAppointmentsForTrainer)
	router.HandleFunc("/available/trainer/{trainer_id}", returnAvailableAppointmentsForTrainer)

	log.Fatal(http.ListenAndServe(":10000", router))
}

func readAppointmentsFromFile() {
	jsonFile, err := os.Open("appointments.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened appointments.json")
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &Appointments)
}

func writeAppointmentsToFile() {
	file, _ := json.MarshalIndent(Appointments, "", " ")
	_ = ioutil.WriteFile("appointments.json", file, 0644)
}


func main() {
	fmt.Println("Rest API using Mux")
	readAppointmentsFromFile()
	handleRequests()
}
