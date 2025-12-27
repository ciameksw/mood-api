package postgres

import (
	"database/sql"
	"errors"
	"time"
)

type MoodType struct {
	ID          int
	Name        string
	Description string
}

// GetMoodTypes retrieves all mood types from the database
func (p *PostgresDB) GetMoodTypes() ([]MoodType, error) {
	var moodTypes []MoodType
	query := "SELECT id, name, description FROM mood_type"

	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mt MoodType
		if err := rows.Scan(&mt.ID, &mt.Name, &mt.Description); err != nil {
			return nil, err
		}
		moodTypes = append(moodTypes, mt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return moodTypes, nil
}

// AddMoodEntry inserts a new mood entry into the database
func (p *PostgresDB) AddMoodEntry(userId int, moodDate string, moodTypeID int, note string) (int, error) {
	var entryID int
	query := "INSERT INTO mood (user_id, mood_date, mood_type_id, note, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err := p.DB.QueryRow(query, userId, moodDate, moodTypeID, note, time.Now()).Scan(&entryID)
	if err != nil {
		return 0, err
	}

	return entryID, nil
}

// GetMoodEntryByDateAndUser retrieves a mood entry for a specific user on a specific date
func (p *PostgresDB) GetMoodEntryByDateAndUser(userId int, moodDate string) (*MoodEntry, error) {
	var me MoodEntry
	query := "SELECT id, user_id, mood_date, mood_type_id, note, created_at FROM mood WHERE user_id = $1 AND mood_date = $2"

	err := p.DB.QueryRow(query, userId, moodDate).Scan(&me.ID, &me.UserID, &me.MoodDate, &me.MoodTypeID, &me.Note, &me.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &me, nil
}

type MoodEntry struct {
	ID         int
	UserID     int
	MoodDate   string
	MoodTypeID int
	Note       string
	CreatedAt  time.Time
}

type GetInput struct {
	UserID    int
	StartDate string
	EndDate   string
}

// GetMoodEntries retrieves mood entries for a user within a date range
func (p *PostgresDB) GetMoodEntries(input GetInput) ([]MoodEntry, error) {
	var moodEntries []MoodEntry
	query := "SELECT id, user_id, mood_date, mood_type_id, note, created_at FROM mood WHERE user_id = $1 AND mood_date BETWEEN $2 AND $3"

	rows, err := p.DB.Query(query, input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var me MoodEntry
		if err := rows.Scan(&me.ID, &me.UserID, &me.MoodDate, &me.MoodTypeID, &me.Note, &me.CreatedAt); err != nil {
			return nil, err
		}
		moodEntries = append(moodEntries, me)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return moodEntries, nil
}

type MoodSummary struct {
	MoodTypeID int
	Count      int
	Percentage float64
}

// GetMoodSummary retrieves a summary of mood entries for a user within a date range
func (p *PostgresDB) GetMoodSummary(input GetInput) ([]MoodSummary, error) {
	var summary []MoodSummary
	query := `
        SELECT 
            mood_type_id, 
            COUNT(*) as count,
            ROUND(100.0 * COUNT(*) / SUM(COUNT(*)) OVER (), 2) as percentage
        FROM mood
        WHERE user_id = $1 AND mood_date BETWEEN $2 AND $3 
        GROUP BY mood_type_id
		ORDER BY count DESC
    `

	rows, err := p.DB.Query(query, input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ms MoodSummary
		if err := rows.Scan(&ms.MoodTypeID, &ms.Count, &ms.Percentage); err != nil {
			return nil, err
		}
		summary = append(summary, ms)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summary, nil
}

// UpdateMoodEntry updates an existing mood entry in the database
func (p *PostgresDB) UpdateMoodEntry(entryID int, moodTypeID int, note string) error {
	query := "UPDATE mood SET mood_type_id = $1, note = $2 WHERE id = $3"

	result, err := p.DB.Exec(query, moodTypeID, note, entryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}

	return nil
}

// DeleteMoodEntry deletes a mood entry from the database
func (p *PostgresDB) DeleteMoodEntry(entryID int) error {
	query := "DELETE FROM mood WHERE id = $1"

	result, err := p.DB.Exec(query, entryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows deleted")
	}

	return nil
}
