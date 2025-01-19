package utils

import (
	"Diplomski/database"
	"fmt"
	"time"
	"github.com/go-co-op/gocron"
)

func ArchiveWorkOrders() error {
    query := "EXEC Arhivaraj" 

    _, err := database.DB.Exec(query)
    if err != nil {
        return err
    }
    return nil
}

func StartScheduler() {
    scheduler := gocron.NewScheduler(time.UTC)

    
    scheduler.Every(1).Day().At("00:00").Do(func() {
        err := ArchiveWorkOrders()
        if err != nil {
            fmt.Println("Archive error", err)
        } else {
            fmt.Println("Work order archived successfully.")
        }
    })

    
    scheduler.StartBlocking()
}
