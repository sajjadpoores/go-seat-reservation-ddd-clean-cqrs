package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	storage "github.com/sajjadpoores/go-seat-reservation-ddd-clean-cqrs/internal/adapter/storage/in-memory"
	"github.com/sajjadpoores/go-seat-reservation-ddd-clean-cqrs/internal/domain/event"
)

func main() {
	eventRepo := storage.NewInMemoryRepository()
	ctx := context.Background()

	concert := event.CreateNewEvent("Eminem concert", 2)

	if err := eventRepo.Save(ctx, concert); err != nil {
		panic(err)
	}

	loadedEvent, err := eventRepo.FindByID(ctx, concert.Id)
	if err != nil {
		panic(err)
	}

	customerID := event.CustomerID(uuid.New())
	hold, err := loadedEvent.RequestHold(customerID, 1, 2*time.Minute, time.Now())
	_, err3 := loadedEvent.RequestHold(customerID, 1, 2*time.Minute, time.Now())

	if err != nil {
		fmt.Println("Could not book:", err)
	} else {
		fmt.Println("Hold successful!")
	}

	if err := eventRepo.Save(ctx, loadedEvent); err != nil {
		panic(err)
	}

	if err := loadedEvent.ConfirmHold(hold.Id, time.Now()); err != nil {
		fmt.Println("Could not confirm hold:", err)
	} else {
		fmt.Println("Hold confirmed!")
	}

	if err := eventRepo.Save(ctx, loadedEvent); err != nil {
		panic(err)
	}

	eventRepo.PrintAllEvents(ctx)
}
