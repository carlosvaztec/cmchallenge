package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/imroc/req/v3"
)

const (
	GoalFile1   = "01_cross_goal.json"
	GoalFile2   = "02_butterfly_goal.json"
	BaseUrl     = "https://challenge.crossmint.io/api/"
	CandidateId = "2d10f365-0233-400b-b680-c4ab1c322de9"
)

func main() {
	runFirstChallenge()
	//runSecondChallenge()
	//resetCandidateMegaworldUsingGoalApiResource()
}

func runFirstChallenge() {
	fmt.Println("Running first challenge")
	megaverse := createMegaverse(11)
	solveFirstChallenge(&megaverse)
	//megaverse.print()
	megaverse.publish()
}

func runSecondChallenge() {
	fmt.Println("Running second challenge")
	megaverse := CreateMegaverseFromFile(GoalFile2)
	//megaverse.print()
	megaverse.publish()
}

func resetCandidateMegaworldUsingGoalApiResource() {
	megaverse := createMegaverseFromGoal()
	megaverse.resetAndPublish()
}

// populate map for first challenge
func solveFirstChallenge(megaverse *Megaverse) {
	start := 5
	content := *megaverse.Content
	content[start][start] = 1
	for i := 1; i < start-1; i++ {
		content[start-i][start-i] = 1
		content[start-i][start+i] = 1
		content[start+i][start-i] = 1
		content[start+i][start+i] = 1
	}
}

// MODEL

// megaverse
type Megaverse struct {
	Content *[][]int
}

// publish map using API
func (megaverse *Megaverse) publish() {
	api := CreateApi(BaseUrl, CandidateId)
	for row, columns := range *megaverse.Content {
		for column, _ := range columns {
			astralObjectId := (*megaverse.Content)[row][column]
			if astralObjectId != 0 {
				api.getApiResourceFor(astralObjectId).createAt(row, column)
				time.Sleep(700 * time.Millisecond) // avoid Too Many Requests error
			}
		}
	}
}

func (megaverse *Megaverse) print() {
	for i, a := range *megaverse.Content {
		//fmt.Print(i)
		//fmt.Print("  ")
		for j, _ := range a {
			fmt.Print((*megaverse.Content)[i][j])
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

// reset candidate map by deleting all astral objects using API
func (m *Megaverse) resetAndPublish() {
	for r, cls := range *m.Content {
		for c, _ := range cls {
			value := (*m.Content)[r][c]
			if value != 0 {
				response, error := CreateApi(BaseUrl, CandidateId).getApiResourceFor(value).deleteAt(r, c)
				if error != nil {
					log.Fatal(error)
				}
				if !response.IsSuccess() {
					log.Fatal(response.GetStatus())
				}
				time.Sleep(700 * time.Millisecond) // lets give the server a break
			}
		}
	}
}

// API

type Api struct {
	BaseUrl     string
	CandidateId string
}

type ApiResource interface {
	createAt(row int, column int) (*req.Response, error)
	deleteAt(row int, column int) (*req.Response, error)
}

func CreateApi(baseUrl string, candidateId string) *Api {
	return &Api{
		BaseUrl:     baseUrl,
		CandidateId: candidateId,
	}
}

func (api *Api) getGoal() (*req.Response, error) {
	return req.EnableDumpAll().R().Get(BaseUrl + "/map/" + CandidateId + "/goal")
}

// API resource names
const (
	Polyanets = "polyanets"
	Soloons   = "soloons"
	Comeths   = "comeths"
)

// return correct API resource for an astral object id
func (api Api) getApiResourceFor(astralObjectId int) ApiResource {
	switch astralObjectId {
	case 1:
		return PolyanetResource{name: Polyanets}
	case 2:
		return ComethResource{name: Comeths, direction: "up"}
	case 3:
		return ComethResource{name: Comeths, direction: "down"}
	case 4:
		return ComethResource{name: Comeths, direction: "left"}
	case 5:
		return ComethResource{name: Comeths, direction: "right"}
	case 6:
		return SoloonResource{name: Soloons, color: "red"}
	case 7:
		return SoloonResource{name: Soloons, color: "white"}
	case 8:
		return SoloonResource{name: Soloons, color: "blue"}
	case 9:
		return SoloonResource{name: Soloons, color: "purple"}
	default:
		log.Fatalf("no api resource for id %d", astralObjectId)
	}
	return PolyanetResource{name: Polyanets}
}

func (api *Api) createRequest(params *RequestParams) *req.Request {
	paramsJson, err := json.Marshal(params)

	if err != nil {
		fmt.Println(err)
	}

	return req.EnableDumpAll().DevMode().R().SetContentType("application/json").SetBody(string(paramsJson))
}

// API resource Polyanets
type PolyanetResource struct {
	name string
}

func (resource PolyanetResource) createAt(row int, column int) (*req.Response, error) {
	api := CreateApi(BaseUrl, CandidateId)

	params := &RequestParams{
		CandidateId: api.CandidateId,
		Row:         row,
		Column:      column,
	}

	return api.createRequest(params).Post(api.BaseUrl + resource.name)
}

func (resource PolyanetResource) deleteAt(row int, column int) (*req.Response, error) {
	api := CreateApi(BaseUrl, CandidateId)

	params := &RequestParams{
		CandidateId: api.CandidateId,
		Row:         row,
		Column:      column,
	}

	return api.createRequest(params).Delete(api.BaseUrl + resource.name)
}

// API resource Saloons
type SoloonResource struct {
	name  string
	color string
}

func (resource SoloonResource) createAt(row int, column int) (*req.Response, error) {
	api := CreateApi(BaseUrl, CandidateId)

	params := &RequestParams{
		CandidateId: api.CandidateId,
		Row:         row,
		Column:      column,
		Color:       resource.color,
	}

	return api.createRequest(params).Post(api.BaseUrl + resource.name)
}

func (resource SoloonResource) deleteAt(row int, column int) (*req.Response, error) {
	api := CreateApi(BaseUrl, CandidateId)

	params := &RequestParams{
		CandidateId: api.CandidateId,
		Row:         row,
		Column:      column,
		Color:       resource.color,
	}

	return api.createRequest(params).Delete(api.BaseUrl + resource.name)
}

// API resource Cometh
type ComethResource struct {
	name      string
	direction string
}

func (resource ComethResource) createAt(row int, column int) (*req.Response, error) {
	api := CreateApi(BaseUrl, CandidateId)

	params := &RequestParams{
		CandidateId: api.CandidateId,
		Row:         row,
		Column:      column,
		Direction:   resource.direction,
	}

	return api.createRequest(params).Post(api.BaseUrl + resource.name)
}

func (resource ComethResource) deleteAt(row int, column int) (*req.Response, error) {
	api := CreateApi(BaseUrl, CandidateId)

	params := &RequestParams{
		CandidateId: api.CandidateId,
		Row:         row,
		Column:      column,
		Direction:   resource.direction,
	}

	return api.createRequest(params).Delete(api.BaseUrl + resource.name)
}

type RequestParams struct {
	CandidateId string `json:"candidateId"`
	Row         int    `json:"row"`
	Column      int    `json:"column"`
	Color       string `json:"color,omitempty"`
	Direction   string `json:"direction,omitempty"`
}

// for goal map parsing
type MegaverseGoal struct {
	Goal [][]string `json:goal`
}

const (
	Polyanet     = "POLYANET"
	ComethUp     = "UP_COMETH"
	ComethDown   = "DOWN_COMETH"
	ComethLeft   = "LEFT_COMETH"
	ComethRight  = "RIGHT_COMETH"
	SoloonRed    = "RED_SOLOON"
	SoloonWhite  = "WHITE_SOLOON"
	SoloonBlue   = "BLUE_SOLOON"
	SoloonPurple = "PURPLE_SOLOON"
	Space        = "SPACE"
)

func astralObjectIdFor(content string) int {
	switch content {
	case Polyanet:
		return 1
	case ComethUp:
		return 2
	case ComethDown:
		return 3
	case ComethLeft:
		return 4
	case ComethRight:
		return 5
	case SoloonRed:
		return 6
	case SoloonWhite:
		return 7
	case SoloonBlue:
		return 8
	case SoloonPurple:
		return 9
	default:
		log.Printf("astral object %s not supported", content)
	}
	return -1 // never called; compiler was not happy
}

// return a Megaverse
func createMegaverse(dimention int) Megaverse {
	// init rows
	content := make([][]int, dimention)

	// init columns
	for i := range content {
		content[i] = make([]int, dimention)
	}

	return Megaverse{
		Content: &content,
	}
}

// goal map parsing
func CreateMegaverseFromFile(jsonFilename string) Megaverse {
	jsonFile, err := os.Open(jsonFilename)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	json, _ := ioutil.ReadAll(jsonFile)

	return createMegaverseFromBytes(json)
}

// fetch map from goal endpoint and populate megaworld
func createMegaverseFromGoal() Megaverse {
	response, err := CreateApi(BaseUrl, CandidateId).getGoal()

	if err != nil || !response.IsSuccess() {
		log.Fatal("fail to fetch goal", err)
	}

	//log.Print(response.String())

	return createMegaverseFromBytes(response.Bytes())
}

func createMegaverseFromBytes(jsonContent []byte) Megaverse {
	// parse json
	var goal MegaverseGoal
	json.Unmarshal(jsonContent, &goal)

	// populate megaverse
	megaverse := createMegaverse(len(goal.Goal))

	for row, a := range goal.Goal {
		for column, _ := range a {
			astralObjectName := goal.Goal[row][column]
			if astralObjectName != Space {
				(*megaverse.Content)[row][column] = astralObjectIdFor(astralObjectName)
			}
		}
	}

	return megaverse
}
