package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"encoding/json"
	"io"
	"os"
	"strings"

	"path/filepath"

	"github.com/ZTPInstaller/ServerUI/libhttp"
)

type ConfigEnvironment struct {
	EnvironmentName     string   `json:"EnvironmentName"`
	Mac                 []string `json:"Mac"`
	InstructionFileName string   `json:"InstructionFileName"`
	AutoUpdate          bool     `json:"AutoUpdate"`
}

const ZTP_SERVER_REST_ENDPOINT = "172.16.128.147:9099"

func GetCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/create/create.html")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, "")
}

func GetEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/create/edit.html")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, "")
}

func RedirectToSubmissionPage(w http.ResponseWriter, err error) {

	if err == nil {
		tmpl, err := template.ParseFiles("templates/create/submit.html")
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

		tmpl.Execute(w, "")
	} else {
		tmpl, err := template.ParseFiles("templates/create/failure.html")
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

		tmpl.Execute(w, "")
	}
}

func FetchAndSubmitReqToPNPServer(w http.ResponseWriter, r *http.Request, reqType string) (err error) {
	var filePath string
	var UploadedFileName string

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		if reqType == "POST" {
			fmt.Println(err)
			return
		} else {
			fmt.Println("Ignoring Empty file received in Edit Request")
		}
	} else {
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)

		dir, err := filepath.Abs(filepath.Dir(handler.Filename))
		if err != nil {
			fmt.Println(err)
			return err
		}

		filePath = filepath.Join(dir, handler.Filename)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer f.Close()
		io.Copy(f, file)
		UploadedFileName = handler.Filename
	}

	fmt.Printf("Uploaded file:%s filePath:%s", UploadedFileName, filePath)
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
	fmt.Printf("********************")

	fmt.Printf(r.Form.Get("maclist"))
	macList := r.Form.Get("maclist")
	envName := r.Form.Get("installname")
	autoUpdate := r.Form.Get("Enable Auto Updates")

	mList := strings.Split(macList, ",")
	var isAutoUpdate bool
	if autoUpdate == "" {
		isAutoUpdate = false
	} else {
		isAutoUpdate = true
	}

	fmt.Printf("ENV Name:%s", envName)
	fmt.Printf("AutoUpdate:%v", isAutoUpdate)
	fmt.Printf("MAC List :%s", mList)

	CfgEnv := ConfigEnvironment{EnvironmentName: envName, Mac: mList, InstructionFileName: filePath, AutoUpdate: isAutoUpdate}
	mapB, err := json.Marshal(CfgEnv)
	if err != nil {
		err = fmt.Errorf("error in marshalling the request to json for applying token : %s", err)
		return
	}
	bodyStr := string(mapB)
	err = SendHTTPRequestToPNPServer(bodyStr, reqType)
	if err != nil {
		return err
	}

	return nil
}

func SendHTTPRequestToPNPServer(bodyStr string, reqType string) (err error) {
	body := strings.NewReader(bodyStr)
	url := "http://" + ZTP_SERVER_REST_ENDPOINT + "/pnp/environment"
	fmt.Println("REST api for Create ENV: %s", url)
	req, err := http.NewRequest(reqType, url, body)
	if err != nil {
		err = fmt.Errorf("error in forming the request: %s ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("failed to trigger HTTP post request to PNP Server")
		return err
	}
	defer resp.Body.Close()

	return
}

func ProcessCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("**** NINJA:  CREATE Handler ****")
	RedirectToSubmissionPage(w, nil)
	_ = FetchAndSubmitReqToPNPServer(w, r, "POST")
}

func ProcessEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("**** NINJA:  EDIT Handler ****")
	RedirectToSubmissionPage(w, nil)
	_ = FetchAndSubmitReqToPNPServer(w, r, "PUT")
}
