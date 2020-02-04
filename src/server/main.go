package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type counters struct {
	sync.Mutex
	view            int
	click           int
	contentSelected string
}

type values struct {
	// JSON encoder can only see field starts with capital letter
	Views  string `json:"Views"`
	Clicks string `json:"Clicks"`
}

type userInfo struct {
	sync.Mutex
	timeStamp time.Time
	reqCount  int
}

var (
	c        = counters{}
	contents = []string{"sports", "entertainment", "business", "education"}
)

// Always use pointer when work with mutexes!
var (
	mutex        = &sync.Mutex{}
	countersChan chan *counters
	mockStore    map[string]values
	userStore    map[string](*userInfo) // Stimulating a databse that keeps users' info
)

func getSelectInfo(content string) string {
	t := time.Now()
	dd, mm, yy, hh, min, ss := t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()
	key := fmt.Sprintf("%s %d/%02d/%02d %02d:%02d:%02d ", content, dd, mm, yy, hh, min, ss)
	return key
}

//These should be static and no need for rate limit
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to EQ Works 😎")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	data := contents[rand.Intn(len(contents))]

	c.Lock()
	c.view++
	c.contentSelected = getSelectInfo(data)
	c.Unlock()

	err := processRequest(r)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	// simulate random click call
	if rand.Intn(100) < 50 {
		processClick(data)
	}

	c.Lock()
	countersChan <- &c
	c.Unlock()

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(data)
	return
}

func processRequest(r *http.Request) error {
	time.Sleep(time.Duration(rand.Int31n(50)) * time.Millisecond)
	return nil
}

func processClick(data string) error {
	c.Lock()
	c.click++
	c.Unlock()

	return nil
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if !isAllowed(r.RemoteAddr) {
		w.WriteHeader(429)
		fmt.Fprint(w, "Slow down, we don't do it here Flash")
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(mockStore)
	return
}

const limit = 20
const period = 60

func isAllowed(ip string) bool {
	// Rate limiter
	if user, ok := userStore[ip]; ok {
		currentTime := time.Now()
		if user.reqCount < limit {
			user.reqCount++
			userStore[ip] = user
			return true
		}
		timePassed := currentTime.Sub(user.timeStamp).Seconds()
		if timePassed <= period {
			return false
		}
		user.reqCount = 1
		user.timeStamp = currentTime
		userStore[ip] = user
		return true
	}
	userStore[ip] = &userInfo{timeStamp: time.Now(), reqCount: 1}
	return true
}

func uploadCounters() error {
	// Clean implement of repetitive task
	// https://stackoverflow.com/questions/16466320/is-there-a-way-to-do-repetitive-tasks-at-intervals
	interval := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-interval.C:
			data := <-countersChan
			mutex.Lock()
			mockStore[data.contentSelected] = values{Views: strconv.Itoa(data.view), Clicks: strconv.Itoa(data.click)}
			mutex.Unlock()
		}
	}
}

func main() {
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/stats/", statsHandler)

	countersChan = make(chan *counters)
	mockStore = make(map[string]values)
	userStore = make(map[string](*userInfo))

	go uploadCounters()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
