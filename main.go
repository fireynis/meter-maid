package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB

type User struct {
	Id        uint64 `gorm:"primary_key"`
	Name             string `gorm:"type:varchar(128);"`
	SlackTeamId      string `gorm:"type:varchar(128);unique_index"`
	MeterFrequencies []MeterMaidFrequency
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type MeterMaidFrequency struct {
	Id        uint64 `gorm:"primary_key"`
	UserId    uint64
	User      User
	Hour      int
	DayOfWeek int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func init() {
	var err error
	db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&MeterMaidFrequency{})
	db.Model(&MeterMaidFrequency{}).AddForeignKey("user_id", "users(id)", "SET NULL", "CASCADE")
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler).Name("home").Methods("POST")

	_, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil {
		log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
	} else {
		log.Fatalln(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	//Should do err check against parse form.
	data := r.Form
	w.Header().Add("Content-Type", "application/json")
	response := "Got it!"
	command := strings.Fields(data.Get("text"))
	if len(command) == 0 {
		go AlertChannel(data)
		response += " Time stored."
	}
	jsonResp, _ := json.Marshal(JsonResponse{Type: "ephemeral", Text: strings.TrimSpace(response)})
	fmt.Fprint(w, string(jsonResp))
}

func AlertChannel(values url.Values) {
	user := User{}
	db.FirstOrInit(&user, User{SlackTeamId:values.Get("team_id")})
	fmt.Println(user)
	user.Name = values.Get("team_domain")
	db.Save(&user)
	db.Create(&MeterMaidFrequency{UserId:user.Id, Hour:time.Now().Hour(), DayOfWeek:int(time.Now().Weekday())})
}

type JsonResponse struct {
	Type string `json:"response_type"`
	Text string `json:"text"`
}
