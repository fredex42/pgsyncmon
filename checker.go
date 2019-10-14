package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func sslModeFromBool(ssl bool) string {
	if ssl == true {
		return "require"
	} else {
		return "disable"
	}
}

func TestRecoveryStatus(user string, ssl bool) (*RecoveryStatus, error) {
	connStr := fmt.Sprintf("user=%s sslmode=%s", user, sslModeFromBool(ssl))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Could not connect to database: %s", err)
		return nil, err
	}

	rows, queryErr := db.Query("select pg_last_xlog_receive_location() \"receive_location\", pg_last_xlog_replay_location() \"replay_location\", pg_is_in_recovery() \"recovery_status\";")
	if queryErr != nil {
		log.Printf("Could not query database: %s", queryErr)
		return nil, queryErr
	}

	builder := NewPostgresXlogLocationBuilder()
	var lastReceiveStr string
	var lastReplayStr string
	var recoveryStatus bool

	rows.Next()
	scanErr := rows.Scan(&lastReceiveStr, &lastReplayStr, &recoveryStatus)
	if scanErr != nil {
		log.Printf("Could not scan returned value from db: %s", scanErr)
		return nil, scanErr
	}

	lastReceiver, convErr := builder.BuildFromString(lastReceiveStr)
	if convErr != nil {
		log.Printf("Could not check last receive: %s", convErr)
		return nil, convErr
	}

	lastReplay, convErr := builder.BuildFromString(lastReplayStr)
	if convErr != nil {
		log.Printf("Could not check last replay: %s", convErr)
		return nil, convErr
	}

	rtn := RecoveryStatus{
		LastXlogReceive: *lastReceiver,
		LastXlogReplay:  *lastReplay,
		IsInRecovery:    recoveryStatus,
	}

	return &rtn, nil
}
