package models

// User represents a user in the leaderboard
type User struct {
    ID       int64  `db:"id" json:"id"`
    Username string `db:"username" json:"username"`
    Rating   int    `db:"rating" json:"rating"`
    Version  int64  `db:"version" json:"version"` 
}

type RankedUser struct {
    Rank     int64  `json:"rank"`
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Rating   int    `json:"rating"`
}

type RatingUpdate struct {
    UserID  int64 `json:"user_id"`
    Rating  int   `json:"rating"`
    Version int64 `json:"version"`
}

// We are using version for conflict resolution in rating updates
// It prevents older updates from overwriting newer ones in async scenarios.