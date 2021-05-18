package mysql

import (
	"database/sql"
	"errors"
	"strconv"

	"quenten.nl/pepebot/models"
)


type BadwordModel struct {
	DB *sql.DB
}

func (m *BadwordModel) GetWord(word string, serverID string) (*models.Badword, error) {
	stmt := `SELECT id, word, serverid, lastsaid FROM badwords
	WHERE word = ? AND serverid = ?`

	row := m.DB.QueryRow(stmt, word, serverID)

	bw := &models.Badword{}

	err := row.Scan(&bw.ID, &bw.Word, &bw.ServerID, &bw.LastSaid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, models.ErrNoRecord
		}else{
			return nil, err
		}
	}

	return bw, nil
}

func (m *BadwordModel) AllWords() (map[string][]string, error) {
	stmt := `SELECT word, serverid FROM badwords`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	type tmp struct{
		word string
		serverid string
	}
	tmp2 := []*tmp{}

	for rows.Next() {
		t := &tmp{}
		err = rows.Scan(&t.word, &t.serverid)
		if err != nil {
			return nil, err
		}

		tmp2 = append(tmp2, t)
	}

	if err = rows.Err(); err != nil{
		return nil, err
	}

	finaltmp := make(map[string][]string)
	for i := 0; i < len(tmp2); i++ {
		
		finaltmp[tmp2[i].serverid] = append(finaltmp[tmp2[i].serverid], tmp2[i].word)
	}

	return finaltmp, nil
}

func (m *BadwordModel) InsertNewWord(word string, serverid string) (int, error) {
	stmt := `INSERT INTO badwords (word, serverid, lastsaid)
	VALUES (?, ?, UTC_TIMESTAMP())`

	id1, err := strconv.Atoi(serverid)
	if err != nil {
		return 0, err
	}
	result, err := m.DB.Exec(stmt, word, id1)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BadwordModel) RemoveWord(word string, serverid string) (error) {
	allwords, err := m.AllWords()
	if err != nil {
		return err
	}
	found := false
	for i := 0; i < len(allwords[serverid]); i++ {
		if allwords[serverid][i] == word {
			found = true
			break
		}
	}
	if !found {
		return errors.New("that word doesn't exist")
	} 

	stmt := `DELETE FROM badwords WHERE word = ? AND serverid = ?`

	result, err := m.DB.Exec(stmt, word, serverid)
	if err != nil {
		return err
	}

	if r, _ := result.RowsAffected(); r == 0 || r > 1 {
		return errors.New("an unknown error occured")
	}

	return nil
}

func (m *BadwordModel) UpdateLastSaid(word string, serverid string) (int, error) {
	stmt := `SELECT id FROM badwords WHERE
	word = ? AND serverid = ?`

	row := m.DB.QueryRow(stmt, word, serverid)
	
	id := 0
	row.Scan(&id)

	stmt = `UPDATE badwords SET lastsaid = UTC_TIMESTAMP() WHERE id = ?`

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	id2, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id2), nil
}