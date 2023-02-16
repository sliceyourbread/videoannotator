package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"papercup/videoannotator/controllers"
	"papercup/videoannotator/models"
)

func main() {

	controllers.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	http.HandleFunc("/addAnnotation", addAnnotationHandler)
	http.HandleFunc("/updateAnnotation", updateHandler)

	http.HandleFunc("/addVideo", addVideoHandler)
	http.HandleFunc("/deleteVideo", deleteHandler)

	http.HandleFunc("/listAnnotations", listHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	u := models.User{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(500)
	}

	link := r.URL.Query().Get("link")
	x, err := u.UpdateAnnotation(controllers.DB.Conn, link)
	w.Write([]byte("updated rows: " + fmt.Sprintf("%v", x)))
}

func addAnnotationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	u := models.User{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		w.WriteHeader(500)
	}

	link := r.URL.Query().Get("link")
	u.AddAnnotation(controllers.DB.Conn, link)
}

func addVideoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	v := models.Video{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(500)
	}

	// check annotation is the correct length
	for i := range v.Annotation {
		if v.Annotation[i].End > v.Length || v.Annotation[i].Start < 0 {
			w.Write([]byte("Invalid annotation length"))
		}
	}

	x, err := v.AddVideo(controllers.DB.Conn)
	if err != nil {

	}

	w.Write([]byte("added video rows: " + string(x)))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	v := models.Video{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		w.WriteHeader(500)
	}

	x, err := v.DeleteVideo(controllers.DB.Conn)
	fmt.Println(err)
	fmt.Println(x)

}

func listHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	u := models.User{}
	err = json.Unmarshal(body, &u)
	if err != nil {
	}

	link := r.URL.Query().Get("link")
	a, err := u.ListAnnotations(controllers.DB.Conn, link)
	fmt.Println(a)
}
