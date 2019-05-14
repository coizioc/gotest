package deathmatch

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// DeathmatchHeader is a string representing the header of a deathmatch message.
const DeathmatchHeader = "__**:anger:DEATHMATCH:anger:**__"

type fighter struct {
	name   string
	health int
}

type attack struct {
	power int
	text  string
}

// Deathmatch performs a deathmatch between two Members.
func Deathmatch(m1 *discordgo.Member, m2 *discordgo.Member) []string {
	rand.Seed(time.Now().UTC().UnixNano())
	attacks, err := getAttacks()
	if err != nil {
		return []string{"error opening attacks"}
	}

	f1 := initFighter(m1)
	f2 := initFighter(m2)
	isF1Turn := false

	dmMessages := []string{}
	currMessage := fmt.Sprintf("%s\n\n\n\n**%s**: %3d/100\n**%s**: %3d/100",
		DeathmatchHeader, f1.name, f1.health, f2.name, f2.health)

	prevAtkText := "\n"
	dmMessages = append(dmMessages, currMessage)

	for {
		// Randomly select the power of the next attack and find an attack that matches that power.
		power := getPower()
		var currAttack *attack
		for {
			currAttack = attacks[rand.Intn(len(attacks))]
			if currAttack.power == power {
				break
			}
		}

		// Calculate the damange of that attack.
		damage := 0
		if rand.Intn(20) != 0 && !(strings.Contains(currAttack.text, "Infinity Gauntlet") && rand.Intn(2) == 0) {
			damage = calculateDamage(power)
		}

		// Create the replace list for the out string and remove the damage from the player's health.
		var replaceList [3]string
		if isF1Turn {
			replaceList[0] = ":arrow_right:"
			replaceList[1] = f1.name
			replaceList[2] = f2.name
			f2.health -= damage
			if f2.health < 0 {
				f2.health = 0
			}
		} else {
			replaceList[0] = ":arrow_left:"
			replaceList[1] = f2.name
			replaceList[2] = f1.name
			f1.health -= damage
			if f1.health < 0 {
				f1.health = 0
			}
		}

		// Add the new attack text to dmMessages.
		currAtkText := getAttackText(currAttack, replaceList, damage, power)
		isF1Turn = !isF1Turn

		currMessage := fmt.Sprintf("%s\n\n%s%s\n**%s**: %3d/100\n**%s**: %3d/100",
			DeathmatchHeader, prevAtkText, currAtkText, f1.name, f1.health, f2.name, f2.health)
		dmMessages = append(dmMessages, currMessage)

		// Check if there is a winner. If so, break the loop and return dmMessages.
		prevAtkText = currAtkText
		if f1.health < 1 {
			currMessage += fmt.Sprintf("\n:trophy: **%s has won!**", f2.name)
			dmMessages = append(dmMessages, currMessage)
			break
		}
		if f2.health < 1 {
			currMessage += fmt.Sprintf("\n:trophy: **%s has won!**", f1.name)
			dmMessages = append(dmMessages, currMessage)
			break
		}
	}

	return dmMessages
}

func initFighter(m *discordgo.Member) *fighter {
	currFighter := new(fighter)
	if m.Nick != "" {
		currFighter.name = m.Nick
	} else {
		currFighter.name = m.User.Username
	}
	currFighter.health = 100
	return currFighter
}

func calculateDamage(power int) int {
	if power == 0 {
		return 0
	} else if power > 10 {
		return 100
	} else {
		damage := power*3 + (rand.Intn(2*power+1) - power)
		return damage
	}
}

func getPower() int {
	u := rand.Intn(123)
	// TODO is messy as hell need to learn cleaner way to recreate the relevant python code equivalent.
	switch {
	case u < 5:
		return 0
	case u < 30:
		return 1
	case u < 52:
		return 2
	case u < 70:
		return 3
	case u < 85:
		return 4
	case u < 97:
		return 5
	case u < 106:
		return 6
	case u < 112:
		return 7
	case u < 116:
		return 8
	case u < 119:
		return 9
	case u < 121:
		return 10
	default:
		return 11
	}
}

func getAttacks() ([]*attack, error) {
	file, err := os.Open("deathmatch.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var attacks []*attack
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currLine := scanner.Text()
		fields := strings.Split(currLine, ";")
		currAttack := new(attack)
		currAttack.power, _ = strconv.Atoi(fields[1])
		currAttack.text = fields[0]

		attacks = append(attacks, currAttack)
	}
	return attacks, scanner.Err()
}

func getAttackText(currAttack *attack, replaceList [3]string, damage int, power int) string {
	currAtkText := strings.Replace(currAttack.text, "$P1", replaceList[1], 1)
	currAtkText = replaceList[0] + strings.Replace(currAtkText, "$P2", replaceList[2], 1)
	if damage == 0 && power != 0 {
		currAtkText = currAtkText[:len(currAtkText)-1] + ", but it misses!"
	}
	currAtkText += fmt.Sprintf(" It does %d damage.\n", damage)
	return currAtkText
}
