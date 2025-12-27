package postgres

import (
	"github.com/lib/pq"
)

type MoodSummaryEntry struct {
	MoodTypeID int     `json:"moodTypeId" validate:"required"`
	Percentage float64 `json:"percentage" validate:"required"`
}

type MoodAdviceMapping struct {
	AdviceTypeID int
	MoodTypeID   int
	Priority     int
}

func (pg *PostgresDB) GetAdviceTypeIDByMoodSummary(moodSummary []MoodSummaryEntry) (int, error) {
	moodTypeIDs := extractMoodTypeIDs(moodSummary)

	query := `
		SELECT advice_type_id, mood_type_id, priority
		FROM public.mood_advice_type_mapping
		WHERE mood_type_id = ANY($1);
	`

	rows, err := pg.DB.Query(query, pq.Array(moodTypeIDs))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	mappings := []MoodAdviceMapping{}
	for rows.Next() {
		var mapping MoodAdviceMapping
		err := rows.Scan(&mapping.AdviceTypeID, &mapping.MoodTypeID, &mapping.Priority)
		if err != nil {
			return 0, err
		}
		mappings = append(mappings, mapping)
	}

	percentageMap := buildPercentageMap(moodSummary)
	scoresMap := calculateAdviceScores(mappings, percentageMap)
	adviceTypeID := findHighestScoredAdviceType(scoresMap)

	return adviceTypeID, nil
}

func extractMoodTypeIDs(moodSummary []MoodSummaryEntry) []int {
	ids := make([]int, len(moodSummary))
	for i, entry := range moodSummary {
		ids[i] = entry.MoodTypeID
	}
	return ids
}

func buildPercentageMap(moodSummary []MoodSummaryEntry) map[int]float64 {
	percentageMap := make(map[int]float64, len(moodSummary))
	for _, entry := range moodSummary {
		percentageMap[entry.MoodTypeID] = entry.Percentage
	}
	return percentageMap
}

func calculateAdviceScores(mappings []MoodAdviceMapping, percentageMap map[int]float64) map[int]float64 {
	scoresMap := make(map[int]float64)
	for _, mapping := range mappings {
		percentage, exists := percentageMap[mapping.MoodTypeID]
		if !exists {
			continue
		}
		score := percentage * (1.0 / float64(mapping.Priority))
		scoresMap[mapping.AdviceTypeID] += score
	}
	return scoresMap
}

func findHighestScoredAdviceType(scoresMap map[int]float64) int {
	var maxScore float64
	var adviceTypeID int
	for atID, score := range scoresMap {
		if score > maxScore {
			maxScore = score
			adviceTypeID = atID
		}
	}
	return adviceTypeID
}

func (pg *PostgresDB) SelectRandomAdviceByAdviceTypeID(adviceTypeID int) (int, string, string, error) {
	var id int
	var title, content string
	query := `
		SELECT id, title, content
		FROM public.advice
		WHERE advice_type_id = $1
		ORDER BY RANDOM()
		LIMIT 1;
	`

	err := pg.DB.QueryRow(query, adviceTypeID).Scan(&id, &title, &content)
	if err != nil {
		return 0, "", "", err
	}

	return id, title, content, nil
}
