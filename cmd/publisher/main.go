package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/icrowley/fake"
	"github.com/nats-io/stan.go"
	"github.com/underthetreee/L0/config"
	"github.com/underthetreee/L0/internal/model"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	pub, err := stan.Connect(cfg.Nats.Cluster, "orders-pub", stan.NatsURL(cfg.Nats.URL))
	if err != nil {
		return err
	}

	for {
		order, err := seedOrder()
		if err != nil {
			return err
		}

		orderBytes, err := json.Marshal(order)
		if err != nil {
			return err
		}

		log.Println("sending order...")
		if err := pub.Publish("orders", orderBytes); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func seedOrder() (model.Order, error) {
	order := model.Order{
		UID:         fake.Password(12, 12, false, true, false),
		TrackNumber: fake.Model(),
		Entry:       fake.Sentence(),
		Delivery: model.Delivery{
			Name:    fake.FullName(),
			Phone:   fake.Phone(),
			Zip:     fake.Zip(),
			City:    fake.City(),
			Address: fake.StreetAddress(),
			Region:  fake.State(),
			Email:   fake.EmailAddress(),
		},
		Payment: model.Payment{
			Transaction:  fake.CharactersN(10),
			RequestID:    fake.CharactersN(8),
			Currency:     fake.Currency(),
			Provider:     fake.Word(),
			Amount:       toInt(fake.Digits()),
			PaymentDT:    int(time.Now().Unix()),
			Bank:         fake.Company(),
			DeliveryCost: toInt(fake.Digits()),
			GoodsTotal:   toInt(fake.Digits()),
			CustomFee:    toInt(fake.Digits()),
		},
		Items: []model.Item{
			{
				ChrtID:      toInt(fake.Digits()),
				TrackNumber: fake.CharactersN(10),
				Price:       toInt(fake.Digits()),
				RID:         fake.Model(),
				Name:        fake.ProductName(),
				Sale:        toInt(fake.Digits()),
				Size:        fake.Word(),
				TotalPrice:  toInt(fake.Digits()),
				NMID:        toInt(fake.Digits()),
				Brand:       fake.Brand(),
				Status:      toInt(fake.Digits()),
			},
			{
				ChrtID:      toInt(fake.Digits()),
				TrackNumber: fake.CharactersN(10),
				Price:       toInt(fake.Digits()),
				RID:         fake.Model(),
				Name:        fake.ProductName(),
				Sale:        toInt(fake.Digits()),
				Size:        fake.Word(),
				TotalPrice:  toInt(fake.Digits()),
				NMID:        toInt(fake.Digits()),
				Brand:       fake.Brand(),
				Status:      toInt(fake.Digits()),
			},
			{
				ChrtID:      toInt(fake.Digits()),
				TrackNumber: fake.CharactersN(10),
				Price:       toInt(fake.Digits()),
				RID:         fake.Model(),
				Name:        fake.ProductName(),
				Sale:        toInt(fake.Digits()),
				Size:        fake.Word(),
				TotalPrice:  toInt(fake.Digits()),
				NMID:        toInt(fake.Digits()),
				Brand:       fake.Brand(),
				Status:      toInt(fake.Digits()),
			},
		},
		Locale:            "ru",
		InternalSignature: fake.Model(),
		CustomID:          fake.CharactersN(6),
		DeliveryService:   fake.Word(),
		ShardKey:          fake.Model(),
		SMID:              toInt(fake.Digits()),
		DateCreated:       time.Now(),
		OofShard:          fake.Word(),
	}
	return order, nil
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
