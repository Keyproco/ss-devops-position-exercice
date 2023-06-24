package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func listRunningPods() {
	host := os.Getenv("HOST")
	token := os.Getenv("TOKEN")
	// TODO use the token mounted within the pod
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := host + "/api/v1/namespaces/default/pods"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	var podList PodList
	err = json.Unmarshal(body, &podList)

	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, pod := range podList.Items {
		if pod.Status.Name == "Running" {
			fmt.Println("Pod Name:", pod.Metadata.Name+" "+pod.Status.Name)

		}

	}

	//fmt.Println("Response:", string(body))

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

	listRunningPods()

}

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
