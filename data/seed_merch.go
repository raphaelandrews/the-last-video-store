package main

import (
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

func seedMerch(s *store.Store) {
	items := []models.MerchItem{
		{ID: "merch-popcorn-bucket", Name: "Popcorn Bucket", Description: "Classic striped popcorn bucket. Unlimited refills for a month.", PointsCost: 50, Stock: 10},
		{ID: "merch-vhs-blank", Name: "Blank VHS Tape", Description: "Record your own movies. Maxell T-120 high grade.", PointsCost: 75, Stock: 5},
		{ID: "merch-poster", Name: "Movie Poster", Description: "Original theatrical poster from the 80s. Random title, mint condition.", PointsCost: 100, Stock: 3},
		{ID: "merch-tshirt", Name: "Store T-Shirt", Description: "The Last Video Store logo tee. Black, cotton, all sizes.", PointsCost: 150, Stock: 8},
		{ID: "merch-free-rental", Name: "Free Rental Coupon", Description: "One free rental on any movie, any format. No late fees.", PointsCost: 200, Stock: 4},
		{ID: "merch-screening", Name: "Private Screening", Description: "After-hours theater access. Bring 5 friends, 2 free rentals each.", PointsCost: 500, Stock: 1},
		{ID: "merch-membership-upgrade", Name: "Tier Upgrade", Description: "Permanent tier upgrade to the next level (max Gold). One-time use.", PointsCost: 1000, Stock: 2},

		{ID: "merch-pokemon-card", Name: "Pokemon TCG Booster", Description: "Vintage Jungle expansion booster pack. Chance of holographic Pikachu — or a Porygon cosplaying as a VHS tape.", PointsCost: 120, Stock: 6},
		{ID: "merch-matrix-pill", Name: "Red Pill / Blue Pill Set", Description: "Resin-cast pill keychain pair in a velvet pouch. Choose wisely — or collect both. Glows under blacklight.", PointsCost: 180, Stock: 4},
		{ID: "merch-matrix-coat", Name: "Neo's Trench Coat", Description: "Full-length black leather-look trench. Lined with Matrix digital rain pattern. One size, dramatic wind sold separately.", PointsCost: 800, Stock: 2},
		{ID: "merch-blade-runner", Name: "Origami Unicorn", Description: "Hand-folded metallic paper unicorn, just like Gaff leaves behind. Comes in a miniature evidence bag.", PointsCost: 90, Stock: 5},
		{ID: "merch-jurassic-amber", Name: "Jurassic Park Amber Cane", Description: "Polished resin cane top with a faux mosquito inclusion. Spared no expense. Does not actually contain dino DNA.", PointsCost: 350, Stock: 3},
		{ID: "merch-godfather-cat", Name: "Marlon Brando Cat Plush", Description: "Plush ginger tabby — the real star of the opening scene. Sits on your lap while you make offers they can't refuse.", PointsCost: 130, Stock: 4},
		{ID: "merch-shining-carpet", Name: "Overlook Carpet Coaster Set", Description: "Set of 4 hexagonal coasters with the iconic carpet pattern. All work and no play not included.", PointsCost: 60, Stock: 10},
		{ID: "merch-pulp-fiction", Name: "Big Kahuna Burger Box", Description: "Tasty burger-shaped tin lunchbox. Royale with cheese styling. That IS a tasty burger.", PointsCost: 140, Stock: 4},
		{ID: "merch-lotr-ring", Name: "One Ring Replica", Description: "Tungsten band with elvish inscription. Comes with a chain — you'll need it. Wearing it does not turn you invisible, unfortunately.", PointsCost: 250, Stock: 3},
		{ID: "merch-ghibli-soot", Name: "Soot Sprite Plushies", Description: "Set of 3 hand-sewn susuwatari from Spirited Away. Feed them konpeito (not included) and they'll carry coal for you.", PointsCost: 110, Stock: 6},
		{ID: "merch-back-future", Name: "Mr. Fusion Prop Replica", Description: "Desktop model of Doc Brown's Mr. Fusion. Banana peel and beer can come pre-loaded. 1.21 gigawatts of style.", PointsCost: 300, Stock: 2},
		{ID: "merch-alien-plush", Name: "Chestburster Plush", Description: "Surprisingly cute plush xenomorph hatchling. Squeeze it and it makes no sound — just stares into your soul.", PointsCost: 160, Stock: 5},
		{ID: "merch-inception-top", Name: "Totem Spinning Top", Description: "Brass spinning top in a felt-lined case. If it stops spinning, you're awake. If not — enjoy the ride.", PointsCost: 200, Stock: 4},
		{ID: "merch-fight-club-soap", Name: "Paper Street Soap Co. Bar", Description: "Handmade pink soap bar with Fight Club emboss. First rule: you do not talk about how good it smells.", PointsCost: 45, Stock: 12},
		{ID: "merch-indy-hat", Name: "Indiana Jones Fedora", Description: "Brown felt fedora, adventure-ready. Bullwhip not included — you'll have to earn that one in the temple.", PointsCost: 280, Stock: 3},
		{ID: "merch-hitchcock-birds", Name: "Crow Plush (Hitchcock Edition)", Description: "Surprisingly heavy plush crow — feels like a thousand of them are staring at you. Bodega Bay not included.", PointsCost: 100, Stock: 5},
		{ID: "merch-akira-pill", Name: "Akira Capsule Jacket Patch", Description: "Embroidered patch of the iconic pill capsule. Iron it onto your red jacket and yell TETSUO!", PointsCost: 70, Stock: 10},
		{ID: "merch-tarantino-feet", Name: "Tarantino Socks", Description: "Limited edition socks — each pair features a different Tarantino character's feet. Maya, is that a foot joke?", PointsCost: 55, Stock: 8},
		{ID: "merch-2001-monolith", Name: "Mini Monolith", Description: "Solid obsidian-black monolith paperweight. 1:4:9 proportions. Touch it and you might evolve.", PointsCost: 220, Stock: 4},
	}
	for i := range items {
		s.CreateMerchItem(&items[i])
	}
}
