package racetools

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type Client struct { //struct to recieve data from velocidrone csv
	VelocidronName string `csv:"Player Name"`
	QualifyingTime string `csv:"Lap Time"`
	X_Pos          string `csv:"-"`
	ModelName      string `csv:"Model Name"`
	X_Country      string `csv:"-"`
}

type Racers struct {
	RacerName      string `csv:"Display Name"`
	VdName         string `csv:"VdName"`
	QualifyingTime string
	ModelName      string
	Id             string `csv:"ID"`
}

type DiscordRacers struct {
	DiscordId string `csv:"DiscordId"`
	VdName    string `csv:"VdName"`
}

// take a list of racers and returns group sets of racers

func GetVdRacers(filename string) []*Client {

	var Clients = []*Client{}
	var OkRaceClass = []*Client{}

	raceFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer raceFile.Close()
	if err := gocsv.UnmarshalFile(raceFile, &Clients); err != nil { // Load clients from file
		panic(err)
	}
	for _, client := range Clients { //clients are the master qual times
		if client.ModelName == "TBS Spec" || client.ModelName == "Twig XL 3" {
			OkRaceClass = append(OkRaceClass, client) // checkedIn seperates the class of quads from the master list
		}
	}
	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return OkRaceClass
}

// For using the voice chat in FMV discord as base group for pairing.
func GetFMVvoice(fileCsv string) []*Racers {

	var FmvRacers = []*Racers{}

	fmvVoiceFile, err := os.OpenFile(fileCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fmvVoiceFile.Close()
	if err := gocsv.UnmarshalFile(fmvVoiceFile, &FmvRacers); err != nil { // Load clients from file
		fmt.Printf("Something broke with FMV CSV: %v", err) //csv needs to be in same folder as main.go for now
	}
	if _, err := fmvVoiceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return FmvRacers
}

func GetDiscordId(discordCsv string) []*DiscordRacers {
	var DiscordRecords = []*DiscordRacers{}

	discordIdFile, err := os.OpenFile(discordCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer discordIdFile.Close()
	if err := gocsv.UnmarshalFile(discordIdFile, &DiscordRecords); err != nil { // Load clients from file
		fmt.Printf("Something broke with discord ID CSV: %v", err) //csv needs to be in same folder as main.go for now
	}
	if _, err := discordIdFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return DiscordRecords
}

func RaceArray(vdList [][]string) [][][]string {
	var maxGroupsize = 8
	var grouplength int
	var totalGroups int
	var modulus int

	racers := len(vdList)
	if racers > 40 {
		maxGroupsize = 10
	}

	for i := 1; i <= maxGroupsize; i++ {
		if racers/i <= maxGroupsize {
			totalGroups = i
			modulus = racers % i
			if modulus == 0 {
				grouplength = racers / i
			} else {
				grouplength = (racers - modulus) / i
			}
			break
		}
	}

	var groupStructure = make([][][]string, totalGroups)
	var c int
	x := modulus

	for i := 1; i <= totalGroups; i++ {
		if x > 0 { // distribues the modulus between the lower teir groups
			racers := vdList[c : i*(grouplength+1)]
			groupStructure[i-1] = racers
			x--
			c += grouplength + 1
		} else { // groups that don't take a modulus
			racers := vdList[c : c+grouplength]
			groupStructure[i-1] = racers
			c += grouplength
		}
	}
	return groupStructure
}

func BindLists(vdl []*Client, fmvl []*Racers, dcl []*DiscordRacers) []*Racers {
	var bound []*Racers

	for _, f := range fmvl {
		for _, d := range dcl {
			if d.DiscordId == f.Id {
				d.VdName = f.VdName

			}
		}
	}
	for _, f := range fmvl {
		for _, v := range vdl {
			if v.VelocidronName == f.VdName || f.RacerName == v.VelocidronName {
				f.QualifyingTime = v.QualifyingTime
				f.ModelName = v.ModelName
				f.RacerName = v.VelocidronName
				bound = append(bound, f)
				break
			}
		}
		if f.VdName == "" {
			f.QualifyingTime = "CHECK IN Please!"
			bound = append(bound, f)
		}
	}
	return bound
}
