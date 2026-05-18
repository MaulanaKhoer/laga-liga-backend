package handlers

import (
	"sort"

	"github.com/gofiber/fiber/v2"

	"laga-liga-backend/internal/config"
	"laga-liga-backend/internal/models"
)

type StandingRow struct {
	Rank      int    `json:"rank"`
	TeamID    uint   `json:"team_id"`
	TeamName  string `json:"team_name"`
	Played    int    `json:"played"`
	Won       int    `json:"won"`
	Draw      int    `json:"draw"`
	Lost      int    `json:"lost"`
	GF        int    `json:"gf"` // Goals For (memasang gol)
	GA        int    `json:"ga"` // Goals Against (kebobolan)
	GD        int    `json:"gd"` // Goal Difference (selisih gol)
	Points    int    `json:"points"`
}

type ScorerRow struct {
	PlayerID   uint   `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	Goals      int    `json:"goals"`
}

func CreateTournament(c *fiber.Ctx) error {
	var input models.Tournament
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
	}

	if err := config.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan ke database"})
	}

	config.DB.Preload("Status").Preload("Sport").First(&input, input.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Turnamen berhasil dibuat!",
		"data":    input,
	})
}

func GetTournaments(c *fiber.Ctx) error {
	var tournaments []models.Tournament

	config.DB.Preload("Status").Preload("Sport").Find(&tournaments)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": tournaments})
}

func GetTournamentByID(c *fiber.Ctx) error {
	var tournament models.Tournament
	id := c.Params("id")

	if err := config.DB.Preload("Status").Preload("Sport").First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": tournament})
}

func UpdateTournament(c *fiber.Ctx) error {
	var tournament models.Tournament
	id := c.Params("id")

	if err := config.DB.First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	var input models.Tournament
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak sesuai!"})
	}

	config.DB.Model(&tournament).Updates(input)
	config.DB.Preload("Status").Preload("Sport").First(&tournament, id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Data berhasil diupdate!", "data": tournament})
}

func DeleteTournament(c *fiber.Ctx) error {
	var tournament models.Tournament
	id := c.Params("id")

	if err := config.DB.First(&tournament, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	config.DB.Delete(&tournament)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Turnamen berhasil dihapus!"})
}

func GetTournamentStandings(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	var tournament models.Tournament
	if err := config.DB.Preload("Teams").First(&tournament, tournamentID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Turnamen tidak ditemukan!"})
	}

	var matches []models.Match
	config.DB.Where("tournament_id = ? AND status = ?", tournamentID, "finished").Find(&matches)

	standingsMap := make(map[uint]*StandingRow)
	for _, team := range tournament.Teams {
		standingsMap[team.ID] = &StandingRow{
			TeamID:   team.ID,
			TeamName: team.Name,
		}
	}

	for _, m := range matches {
		homeRow := standingsMap[m.HomeTeamID]
		awayRow := standingsMap[m.AwayTeamID]

		if homeRow == nil || awayRow == nil {
			continue
		}

		homeRow.Played++
		awayRow.Played++

		hScore := 0
		if m.HomeScore != nil {
			hScore = *m.HomeScore
		}
		aScore := 0
		if m.AwayScore != nil {
			aScore = *m.AwayScore
		}

		homeRow.GF += hScore
		homeRow.GA += aScore
		awayRow.GF += aScore
		awayRow.GA += hScore

		if hScore > aScore {
			homeRow.Won++
			homeRow.Points += 3
			awayRow.Lost++
		} else if hScore < aScore {
			awayRow.Won++
			awayRow.Points += 3
			homeRow.Lost++
		} else {
			homeRow.Draw++
			awayRow.Draw++
			homeRow.Points += 1
			awayRow.Points += 1
		}
	}

	var standings []StandingRow
	for _, row := range standingsMap {
		row.GD = row.GF - row.GA
		standings = append(standings, *row)
	}

	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		if standings[i].GD != standings[j].GD {
			return standings[i].GD > standings[j].GD
		}
		return standings[i].GF > standings[j].GF
	})

	for i := range standings {
		standings[i].Rank = i + 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tournament": fiber.Map{"id": tournament.ID, "name": tournament.Name},
		"data":       standings,
	})
}

func GetTopScorers(c *fiber.Ctx) error {
	tournamentID := c.Params("id")

	var scorers []ScorerRow
	
	err := config.DB.Table("match_events").
		Select("players.id as player_id, players.name as player_name, teams.name as team_name, count(match_events.id) as goals").
		Joins("JOIN players ON players.id = match_events.player_id").
		Joins("JOIN teams ON teams.id = match_events.team_id").
		Joins("JOIN matches ON matches.id = match_events.match_id").
		Where("matches.tournament_id = ? AND match_events.type = ?", tournamentID, "goal").
		Group("players.id, players.name, teams.name").
		Order("goals DESC").
		Scan(&scorers).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data top scorer"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tournament_id": tournamentID,
		"data":          scorers,
	})
}
