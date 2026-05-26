package database

import "math/rand"

func Seed() {
	users := []User{
		{
			Name:     "John Doe",
			Email:    "john@doe.com",
			Password: "P@ssw0rd",
		},
		{
			Name:     "Jane Doe",
			Email:    "jane@doe.com",
			Password: "P@ssw0rd",
		},
		{
			Name:     "Johan Doe",
			Email:    "johan@doe.com",
			Password: "P@ssw0rd",
		},
	}

	for _, user := range users {
		GetInstance().Create(&user)
		for range rand.Intn(5) + 3 {
			trade := Trade{
				UserID:      user.ID,
				Cost:        int64(rand.Intn(1000)),
				SoldItems:   nil,
				BoughtItems: nil,
			}
			GetInstance().Create(&trade)

			var soldItems []TradeItem
			var boughtItems []TradeItem
			for range rand.Intn(10) + 20 {
				item := TradeItem{
					TradeID:  trade.ID,
					ItemID:   uint(rand.Uint32()),
					Quantity: 1,
				}
				if rand.NormFloat64() <= 8.0 {
					soldItems = append(soldItems, item)
				} else {
					boughtItems = append(boughtItems, item)
				}
				GetInstance().Create(&item)
			}

			trade.SoldItems = soldItems
			trade.BoughtItems = boughtItems
			GetInstance().Save(&trade)
			user.Trades = append(user.Trades, trade)
			GetInstance().Save(&user)
		}
	}
}
