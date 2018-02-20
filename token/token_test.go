package token

import (
	"log"
	"testing"
)

// TestToken tests some implementations of selecting records with a token
// each record should be returned to client at least once (dups are ok)
func TestToken(t *testing.T) {
	k := NewToken()
	t.Run("select", func(t *testing.T) {
		log.Printf("SIMPLE SELECTION TEST")
		k.createDB()
		defer k.destroyDB()

		records := 10
		// go make some records
		for i := 0; i < records; i++ {
			k.insertData("", 0)
		}
		// go select some records
		var token string
		counter := make(map[int]int)
		for {
			data, newToken := k.selectData(token)
			for _, d := range data {
				counter[d.ID]++
			}
			log.Printf("%s", toJSON(data, token, newToken))
			if newToken == "" {
				log.Printf("DONE! no more records")
				break
			}
			token = newToken
		}
		if len(counter) != records {
			t.Error("not enough or duplicate records")
		}
	})
	t.Run("select with updated records", func(t *testing.T) {
		log.Printf("SELECTION WITH UPDATES TEST")
		k.createDB()
		defer k.destroyDB()

		records := 5
		// go make some records
		for i := 0; i < records; i++ {
			k.insertData("", 0)
		}
		// go select some records
		var token string
		counter := make(map[int]int)
		for {
			data, newToken := k.selectData(token)
			for _, d := range data {
				counter[d.ID]++
			}
			log.Printf("%s", toJSON(data, token, newToken))
			if newToken == "" {
				log.Printf("DONE! no more records")
				break
			}
			// save the last token
			token = newToken
		}
		// go update a record
		updatedRecord := 3
		updatedData := "new hotness"
		log.Printf("UPDATING RECORD %d with %s", updatedRecord, updatedData)
		k.updateData(updatedRecord, updatedData)
		// now select again
		for {
			data, newToken := k.selectData(token)
			for _, d := range data {
				counter[d.ID]++
			}
			log.Printf("%s", toJSON(data, token, newToken))
			if newToken == "" {
				log.Printf("DONE! no more records")
				break
			}
			token = newToken
		}
		if len(counter) != records {
			t.Error("not enough or duplicate records")
		}
		for i := 0; i < records; i++ {
			// +1 to i due to 1 index of postgres
			val, ok := counter[i+1]
			if !ok {
				t.Errorf("expected %d in counter, but not found", i+1)
			}
			if val == 0 {
				t.Errorf("got a 0 value for %d, but should be at least 1", i+1)
			}
		}
	})
	t.Run("select with new records", func(t *testing.T) {
		log.Printf("SELECTION WITH NEW RECORDS TEST")
		k.createDB()
		defer k.destroyDB()

		records := 5
		// go make some records
		for i := 0; i < records; i++ {
			k.insertData("", 0)
		}
		// go select some records
		var token string
		counter := make(map[int]int)
		for {
			data, newToken := k.selectData(token)
			for _, d := range data {
				counter[d.ID]++
			}
			if newToken == "" {
				log.Printf("DONE! no more records")
				break
			}
			token = newToken
		}
		log.Printf("making %d more records", records)
		// go make some records
		for i := 0; i < records; i++ {
			k.insertData("", 0)
		}
		// now select again
		for {
			data, newToken := k.selectData(token)
			for _, d := range data {
				counter[d.ID]++
			}
			log.Printf("%s", toJSON(data, token, newToken))
			if newToken == "" {
				log.Printf("DONE! no more records")
				break
			}
			token = newToken
		}
		if len(counter) != records*2 {
			t.Error("not enough or duplicate records")
		}
		for i := 0; i < records*2; i++ {
			// +1 to i due to 1 index of postgres
			val, ok := counter[i+1]
			if !ok {
				t.Errorf("expected %d in counter, but not found", i+1)
			}
			if val == 0 {
				t.Errorf("got a 0 value for %d, but should be at least 1", i+1)
			}
		}
	})
}
