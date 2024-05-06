package performance

import (
	"database-tester/mysqldb"
	"database-tester/types"
	"fmt"
	"math"
	"os"
	"time"
)

type StorageManager interface {
	InsertNote(note types.Note) (id string, err error)
	GetNote(id string) (note types.Note, err error)
	GetNotes(userID string) (notes []types.Note, err error)
	PatchNote(note types.Note) (id string, err error)
	DeleteNote(note types.Note) (id string, err error)

	InsertUser(user types.User) (id string, err error)
	GetUser(id string) (user types.User, err error)
	PatchUser(user types.User) (id string, err error)
	GetUserStats(userID string) (stats types.UserStats, err error)
	GetUserModifiedNotesStats(userID string) (stats types.ModifiedNotesStats, err error)
	DeleteUser(user types.User) (id string, err error)
}

type PerformanceSuite struct {
	StorageManager StorageManager
	logFile        *os.File
}

//func RunPerformanceTest() {
//
//	fmt.Println("Creating DynamoDB manager")
//	manager, err := dynamodb.NewDynamoManager()
//	if err != nil {
//		fmt.Println("Error creating DynamoDB manager:", err)
//		return
//	}
//
//	file, err := os.OpenFile("performance/data/dynamo.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		fmt.Println("Error opening performance.log:", err)
//		return
//	}
//	defer file.Close()
//
//	fmt.Println("Reading data from files")
//	users, notes, updatedNotes, err := readDataFromFiles()
//	if err != nil {
//		fmt.Println("Error reading files:", err)
//		return
//	}
//
//	fmt.Println("Starting performance test")
//	ps := PerformanceSuite{
//		StorageManager: manager,
//		logFile:        file,
//	}
//	ps.measureInsertUserPerformance(users)
//	ps.measureInsertNotePerformance(notes)
//	ps.measureGetUserPerformance(users)
//	ps.measureGetNotePerformance(notes)
//	ps.measureGetUserNotesPerformance(users)
//	ps.measureGetUserStatsPerformance(users)
//	ps.measurePatchUserPerformance(users)
//	ps.measurePatchNotePerformance(updatedNotes)
//	ps.getUserModifiedNotesStatsPerformance(users)
//	ps.measureDeleteNotePerformance(notes)
//	ps.measureDeleteUserPerformance(users)
//
//}

func RunMySQLPerformanceTest() {
	fmt.Println("Creating MySQLDB manager")
	manager, err := mysqldb.NewMySQLManager()
	if err != nil {
		fmt.Println("Error creating MySQLDB manager:", err)
		return
	}

	file, err := os.OpenFile("performance/data/mysql.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening mysql.log:", err)
		return
	}
	defer file.Close()

	fmt.Println("Reading data from files")
	users, notes, updatedNotes, err := readDataFromFiles()
	if err != nil {
		fmt.Println("Error reading files:", err)
		return
	}
	fmt.Println("Starting MySQL performance test")
	ps := PerformanceSuite{
		StorageManager: manager,
		logFile:        file,
	}
	ps.measureInsertUserPerformance(users)
	ps.measureInsertNotePerformance(notes)
	ps.measureGetUserPerformance(users)
	ps.measureGetNotePerformance(notes)
	ps.measureGetUserNotesPerformance(users)
	ps.measureGetUserStatsPerformance(users)
	ps.measurePatchUserPerformance(users)
	ps.measurePatchNotePerformance(updatedNotes)
	ps.getUserModifiedNotesStatsPerformance(users)
	ps.measureDeleteNotePerformance(notes)
	ps.measureDeleteUserPerformance(users)
}

func (ps PerformanceSuite) measureInsertUserPerformance(users []types.User) {

	start := time.Now()
	for i, user := range users {
		_, err := ps.StorageManager.InsertUser(user)
		if err != nil {
			fmt.Println("Error inserting user:", err)
			return
		}
		ps.logPerformanceRecord("Insert", "User", i, len(users), time.Since(start))
	}
}

func (ps PerformanceSuite) measureInsertNotePerformance(notes []types.Note) {

	start := time.Now()
	for i, note := range notes {
		_, err := ps.StorageManager.InsertNote(note)
		if err != nil {
			fmt.Println("Error inserting note:", err)
			return
		}

		ps.logPerformanceRecord("Insert", "Note", i, len(notes), time.Since(start))
	}
}

func (ps PerformanceSuite) measureGetUserPerformance(users []types.User) {

	start := time.Now()
	for i, u := range users {
		_, err := ps.StorageManager.GetUser(u.ID)
		if err != nil {
			fmt.Println("Error getting user:", err)
			return
		}

		ps.logPerformanceRecord("Get", "User", i, len(users), time.Since(start))
	}
}

func (ps PerformanceSuite) measureGetNotePerformance(notes []types.Note) {

	start := time.Now()
	for i, n := range notes {
		_, err := ps.StorageManager.GetNote(n.ID)
		if err != nil {
			fmt.Println("Error getting note:", err)
			return
		}

		ps.logPerformanceRecord("Get", "Note", i, len(notes), time.Since(start))
	}
}

func (ps PerformanceSuite) measureGetUserNotesPerformance(users []types.User) {

	start := time.Now()
	notesPerUser := 1000
	for i, u := range users {
		_, err := ps.StorageManager.GetNotes(u.ID)
		if err != nil {
			fmt.Println("Error getting user notes:", err)
			return
		}

		ps.logPerformanceRecord("Get", "User Notes", i, len(users)*notesPerUser, time.Since(start))
	}
}

func (ps PerformanceSuite) measureGetUserStatsPerformance(users []types.User) {

	start := time.Now()
	for i, u := range users {
		_, err := ps.StorageManager.GetUserStats(u.ID)
		if err != nil {
			fmt.Println("Error getting user stats:", err)
			return
		}

		ps.logPerformanceRecord("Get", "User Stats", i, len(users), time.Since(start))
	}
}

func (ps PerformanceSuite) measurePatchUserPerformance(users []types.User) {

	start := time.Now()
	for i, user := range users {
		_, err := ps.StorageManager.PatchUser(user)
		if err != nil {
			fmt.Println("Error patching user:", err)
			return
		}

		ps.logPerformanceRecord("Patch", "User", i, len(users), time.Since(start))

	}
}

func (ps PerformanceSuite) measurePatchNotePerformance(notes []types.Note) {

	start := time.Now()
	for i, note := range notes {
		_, err := ps.StorageManager.PatchNote(note)
		if err != nil {
			fmt.Println("Error patching note:", err)
			return
		}

		ps.logPerformanceRecord("Patch", "Note", i, len(notes), time.Since(start))
	}
}

func (ps PerformanceSuite) getUserModifiedNotesStatsPerformance(users []types.User) {

	start := time.Now()
	for i, u := range users {
		_, err := ps.StorageManager.GetUserModifiedNotesStats(u.ID)
		if err != nil {
			fmt.Println("Error getting user modified notes stats:", err)
			return
		}
		ps.logPerformanceRecord("Get", "User Modified Notes Stats", i, len(users), time.Since(start))
	}
}

func (ps PerformanceSuite) measureDeleteUserPerformance(users []types.User) {

	start := time.Now()
	for i, user := range users {
		_, err := ps.StorageManager.DeleteUser(user)
		if err != nil {
			fmt.Println("Error deleting user:", err)
			return
		}

		ps.logPerformanceRecord("Delete", "User", i, len(users), time.Since(start))
	}
}

func (ps PerformanceSuite) measureDeleteNotePerformance(notes []types.Note) {

	start := time.Now()
	for i, note := range notes {
		_, err := ps.StorageManager.DeleteNote(note)
		if err != nil {
			fmt.Println("Error deleting note:", err)
			return
		}

		ps.logPerformanceRecord("Delete", "Note", i, len(notes), time.Since(start))
	}
}

func (ps PerformanceSuite) logPerformanceRecord(method, dataType string, iteration, expectedNumberOfIterations int, elapsed time.Duration) {

	if iteration%100 == 0 {
		fmt.Println("\nMethod:", method, "Data Type:", dataType, "Iteration:", iteration, "Elapsed:", elapsed)
		fmt.Println("Percent done:", float32(iteration)/float32(expectedNumberOfIterations)*100)
		fmt.Printf("Time elapsed: %.0f seconds\n", elapsed.Seconds())
		remaining := float64(expectedNumberOfIterations-iteration) * elapsed.Seconds() / float64(iteration)
		fmt.Printf("Time remaining: %.0f seconds\n", remaining)
	}

	iteration++
	if math.Log10(float64(iteration)) == math.Floor(math.Log10(float64(iteration))) {
		ps.logFile.WriteString(fmt.Sprintf("Method: %v, Data Type: %v, Iteration: %v, Elapsed: %v,\n", method, dataType, iteration, elapsed))
	}
}

//test batch read of 10k records
