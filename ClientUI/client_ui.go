// Copyright to TechNinja Team
//
//

package main

import (
	"fmt"
	"strconv"

	"strings"

	"github.com/go-redis/redis"
	"github.com/icza/gowut/gwu"
)

const (
	Install  = 1
	Rollback = 2
	Update   = 3
	NoAction = 4
)

const (
	Success = 1
	Failure = 2
)

type SoftwareDB struct {
	Name         string
	Version      string
	AvailVersion string
	Action       int
	Status       int
}

func SetDataListInDB(client *redis.Client, SDBList []*SoftwareDB) {

	for _, SDB := range SDBList {

		SetDataInDB(client, SDB)
	}
}

func SetDataInDB(client *redis.Client, SDB *SoftwareDB) {
	client.HSet(SDB.Name, "Name", SDB.Name)
	client.HSet(SDB.Name, "Version", SDB.Version)
	client.HSet(SDB.Name, "AvailVersion", SDB.AvailVersion)
	client.HSet(SDB.Name, "Action", SDB.Action)
	client.HSet(SDB.Name, "Status", SDB.Status)
}

func PrepareKubernetesSetupDummyData() []*SoftwareDB {
	SDBList := []*SoftwareDB{}

	SDB1 := &SoftwareDB{}
	SDB1.Status = Success
	SDB1.Action = Update
	SDB1.Name = "kubernetes"
	SDB1.Version = "1.9.3-00"
	SDB1.AvailVersion = "1.10.3"

	SDB2 := &SoftwareDB{}
	SDB2.Status = Success
	SDB2.Action = NoAction
	SDB2.Name = "docker-ce"
	SDB2.Version = "17.03.2~ce-0~ubuntu-xenial"
	SDB2.AvailVersion = "17.03.2~ce-0~ubuntu-xenial"

	SDBList = append(SDBList, SDB1)
	SDBList = append(SDBList, SDB2)
	return SDBList
}

func PrepareKubernetesKeyList() (keyList []string) {

	keyList = []string{"kubernetes", "docker-ce"}
	return
}

func DatabaseOperation(DBClient *redis.Client, ClientUI gwu.Window) {
	SDPDataLIst := PrepareKubernetesSetupDummyData()
	SetDataListInDB(DBClient, SDPDataLIst)

	// Display software details at Ninja Client UI
	keyList := PrepareKubernetesKeyList()
	DisplayAtNinjaClientUI(DBClient, ClientUI, keyList)
}

func GetDataFromDataBase(key string, client *redis.Client) (out SoftwareDB) {
	fmt.Printf("PNP Client DATABASE ")
	outSDB := SoftwareDB{}
	outSDB.Name = client.HGet(key, "Name").Val()
	outSDB.Version = client.HGet(key, "Version").Val()
	outSDB.AvailVersion = client.HGet(key, "AvailVersion").Val()
	outSDB.Action, _ = strconv.Atoi(client.HGet(key, "Action").Val())
	outSDB.Status, _ = strconv.Atoi(client.HGet(key, "Status").Val())

	fmt.Printf("*** Software Details # Name: %s, Version: %s AvailVersion: %s Action:%v, Status: %v",
		outSDB.Name, outSDB.Version, outSDB.AvailVersion, outSDB.Action, outSDB.Status)

	data := client.HGetAll(key)
	fmt.Printf("All Value: %+v", data)
	return outSDB
}

func DisplayAtNinjaClientUI(DBClient *redis.Client, win gwu.Window, keyList []string) {
	// Fetching data from Database for all keys
	var sdb [2]SoftwareDB
	sdb[0] = GetDataFromDataBase(keyList[0], DBClient)
	sdb[1] = GetDataFromDataBase(keyList[1], DBClient)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(20)

	t := gwu.NewTable()
	t.Style().SetBorder2(10, gwu.BrdStyleSolid, gwu.ClrNavy)
	t.SetAlign(gwu.HARight, gwu.VATop)
	t.Style().SetSize("1000", "500")
	t.EnsureSize(2, 2)
	t.RowFmt(0).Style().SetBackground(gwu.ClrNavy)

	t.RowFmt(0).SetAlign(gwu.HADefault, gwu.VAMiddle)
	t.RowFmt(1).SetAlign(gwu.HADefault, gwu.VAMiddle)


	img := gwu.NewImage(fmt.Sprintf("Installed Software"), "http://www2.multilizer.com/wp-content/uploads/2014/07/tool.jpg")
	img.Style().SetSize("70", "50")
	t.Add(img, 0, 0)

	lb1 := gwu.NewLabel(fmt.Sprintf("Current Version"))
	//lb1.Style().SetBackground("blue")
	lb1.Style().SetColor("white")
	lb1.Style().SetWidth("20")

	lb2 := gwu.NewLabel(fmt.Sprintf("Available Version"))
	lb2.Style().SetColor("white")
	lb2.Style().SetWidth("20")

	lb3 := gwu.NewLabel(fmt.Sprintf("Status"))
	lb3.Style().SetColor("white")
	lb3.Style().SetWidth("20")

	lb4 := gwu.NewLabel(fmt.Sprintf("Action"))
	lb4.Style().SetColor("white")
	lb4.Style().SetWidth("20")

	t.Add(lb1, 0, 1)
	t.Add(lb2, 0, 2)
	t.Add(lb3, 0, 3)
	t.Add(lb4, 0, 4)

	btnsPanel := gwu.NewNaturalPanel()
	for row := 1; row < 3; row++ {
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb[row-1].Name)), row, 0)
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb[row-1].Version)), row, 1)
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", sdb[row-1].AvailVersion)), row, 2)

		var statusStr string
		if sdb[row-1].Status == Success {
			statusStr = "Operation Success"
		} else {
			statusStr = "Operation Failure"
		}
		t.Add(gwu.NewLabel(fmt.Sprintf("%s", statusStr)), row, 3)

		var actionStr string
		if sdb[row-1].Action == Install {
			actionStr = "INSTALL"
		} else if sdb[row-1].Action == Rollback {
			actionStr = "ROLLBACK"
		} else if sdb[row-1].Action == Update {
			actionStr = "UPDATE"
		} else if sdb[row-1].Action == NoAction {
			actionStr = "NOACTION"
		} else {
			actionStr = ""
		}

		butn1 := gwu.NewButton(fmt.Sprintf("%s", actionStr))
		butn1.Style().SetColor("white")
		butn1.Style().SetBackground("green")
		//name := "button"+ strconv.Itoa(row)
		name := "button" + sdb[row-1].Name
		butn1.SetAttr("ID", name)

		butn1.AddEHandlerFunc(func(e gwu.Event) {
			if butn1.Text() == "UPDATE" {
				fmt.Printf("UPDATE button pressed!")
				val := butn1.Attr("ID")
				if strings.Contains(val, "kubernetes") {
					fmt.Printf("Hey last UPDATE action was for Kubernetes Software")
				} else if strings.Contains(val, "docker-ce") {
					fmt.Printf("Hey last UPDATE action was for Docker Software")
				}
			} else if butn1.Text() == "INSTALL" {
				fmt.Printf("INSTALL button pressed!")
			} else if butn1.Text() == "ROLLBACK" {
				fmt.Printf("ROLLBACK button pressed!")
				val := butn1.Attr("ID")
				if strings.Contains(val, "kubernetes") {
					fmt.Printf("Hey last ROLLBACK action was for Kubernetes Software")
				} else if strings.Contains(val, "docker-ce") {
					fmt.Printf("Hey last ROLLBACK action was for Docker Software")
				}
			} else if butn1.Text() == "NOACTION" {
				fmt.Printf("NOACTION button pressed!")
			} else {
				fmt.Printf("UNKNOWN button pressed!")
			}
		}, gwu.ETypeClick)

		t.Add(butn1, row, 4)
	}

	p.Add(t)
	p.Add(btnsPanel)
	win.Add(p)
}

func main() {
	//  Master window
	masterWin := gwu.NewWindow("web-ui-dashboard", "TECH-NINJA CLIENT GUI !")
	masterWin.Style().SetFullSize()
	masterWin.SetAlign(gwu.HACenter, gwu.VAMiddle)

	/* Master window */
	p4 := gwu.NewPanel()
	p4.SetHAlign(gwu.HACenter)
	p4.SetCellPadding(2)
	l1 := gwu.NewLabel("Welcome to TechNinja Dashboard")
	l1.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("300%")
	l1.Style().SetColor("green")
	l1.Style().SetBackground("while")
	p4.Add(l1)
	masterWin.Add(p4)

	/* Display window for software Catalog */
	ClientWin := gwu.NewWindow("display-ui", "TECH-NINJA CLIENT GUI!")
	ClientWin.Style().SetFullWidth()
	ClientWin.SetHAlign(gwu.HACenter)
	ClientWin.SetCellPadding(2)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HACenter)
	p.SetCellPadding(2)
	l2 := gwu.NewLabel("PNP Software Catalog")
	l2.Style().SetFontWeight(gwu.FontWeightBold).SetFontSize("300%")
	l2.Style().SetColor("green")
	l2.Style().SetBackground("while")
	p.Add(l2)
	ClientWin.Add(p)

	// Database object creation
	DBClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6389",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := DBClient.Ping().Result()
	fmt.Println(pong, err)
	DatabaseOperation(DBClient, ClientWin)

	// Adding all windows to server
	server := gwu.NewServer("techninja.com", "localhost:8081")
	server.SetText("Starting Tech Ninja!!")
	server.AddWin(ClientWin)
	server.AddWin(masterWin)

	//server.Start()
	server.Start("display-ui")
}
