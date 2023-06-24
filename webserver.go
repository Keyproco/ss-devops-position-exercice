package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type PodList struct {
	Kind       string    `json:"kind"`
	APIVersion string    `json:"apiVersion"`
	Metadata   Metadata  `json:"metadata"`
	Items      []PodItem `json:"items"`
}

type Metadata struct {
	ResourceVersion string `json:"resourceVersion"`
}

type PodItem struct {
	Metadata PodMetadata `json:"metadata"`
	Status   PodStatus   `json:"status"`
}
type PodStatus struct {
	Name string `json:"phase"`
}

type PodMetadata struct {
	Name string `json:"name"`
}

func listPods() []byte {
	host := os.Getenv("HOST")
	//token := os.Getenv("TOKEN")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := host + "/api/v1/namespaces/default/pods"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return nil
	}

	// retrieve token
	bToken, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+string(bToken))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
	}
	return body
}

func filterPodsByStatus(status string, items []PodItem) []PodItem {

	var filtered []PodItem

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Status.Name), strings.ToLower(status)) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func podsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	body := listPods()

	var podList PodList
	err := json.Unmarshal(body, &podList)

	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	podsRunning := filterPodsByStatus("Running", podList.Items)

	if len(podsRunning) == 0 {
		http.Error(w, "There are no pods running", http.StatusOK)
		return
	}

	tmpl := template.Must(template.ParseFiles("pods-running-tpl.html"))
	err = tmpl.Execute(w, podsRunning)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// /var/run/secrets/kubernetes.io/serviceaccount
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", podsHandler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
