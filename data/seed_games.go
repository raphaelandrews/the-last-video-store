package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

func seedGames(s *store.Store) {
	games := []struct {
		Title       string
		Year        int
		Genre       string
		Platform    string
		Director    string
		Synopsis    string
		CopiesTotal int
		PlayPrice   float64
		RentalPrice float64
	}{
		{"Super Mario Bros. 3", 1988, "Platformer", "NES", "Nintendo", "Mario and Luigi embark on a quest across 8 worlds to rescue Princess Toadstool from Bowser and his Koopalings.", 5, 1.99, 2.99},
		{"The Legend of Zelda", 1986, "Action", "NES", "Nintendo", "Link explores the land of Hyrule to collect the Triforce pieces and rescue Princess Zelda from the evil Ganon.", 4, 1.99, 2.99},
		{"Metroid", 1986, "Action", "NES", "Nintendo", "Bounty hunter Samus Aran explores the planet Zebes to stop the Space Pirates from weaponizing the Metroids.", 3, 1.99, 2.99},
		{"Mega Man 2", 1988, "Action", "NES", "Capcom", "Mega Man battles through eight new Robot Masters to defeat Dr. Wily once again. Widely considered the series peak.", 4, 1.99, 2.99},
		{"Castlevania", 1986, "Action", "NES", "Konami", "Simon Belmont enters Dracula's castle armed with the Vampire Killer whip to end the Count's reign of terror.", 3, 1.99, 2.99},
		{"Contra", 1987, "Action", "NES", "Konami", "Two commandos battle alien forces in the jungle. Up, Up, Down, Down, Left, Right, Left, Right, B, A.", 5, 1.99, 2.49},
		{"Duck Hunt", 1984, "Sports", "NES", "Nintendo", "Take aim at ducks and clay pigeons with the NES Zapper. That laughing dog will never let you forget a miss.", 3, 1.49, 1.99},
		{"Tetris", 1989, "Puzzle", "NES", "Nintendo", "Arrange falling tetrominoes into complete lines. The perfect puzzle game that conquered the world.", 6, 1.49, 1.99},
		{"Punch-Out!!", 1987, "Sports", "NES", "Nintendo", "Little Mac fights his way through the World Video Boxing Association circuit. Glass Joe, Bald Bull, Mr. Dream.", 4, 1.49, 2.49},
		{"Super Mario World", 1990, "Platformer", "SNES", "Nintendo", "Mario and Luigi explore Dinosaur Land to rescue Princess Peach. Introduces Yoshi and 96 exits to find.", 5, 2.49, 3.49},
		{"The Legend of Zelda: A Link to the Past", 1991, "Action", "SNES", "Nintendo", "Link travels between the Light and Dark worlds of Hyrule to rescue seven maidens and defeat Agahnim.", 4, 2.49, 3.49},
		{"Super Metroid", 1994, "Action", "SNES", "Nintendo", "Samus returns to Zebes to rescue the last Metroid. A masterclass in atmosphere and exploration.", 3, 2.49, 3.49},
		{"Chrono Trigger", 1995, "RPG", "SNES", "Square", "A boy, a princess, and a robot travel through time to prevent a world-ending catastrophe. 13 endings.", 4, 2.99, 3.99},
		{"Final Fantasy III", 1994, "RPG", "SNES", "Square", "Terra and a band of rebels fight to save a dying world from the Empire. Magitek, espers, and Kefka's madness.", 3, 2.99, 3.99},
		{"Donkey Kong Country", 1994, "Platformer", "SNES", "Rare", "Donkey and Diddy Kong battle through Kremling-infested jungles to recover their stolen banana hoard.", 4, 2.49, 3.49},
		{"Street Fighter II Turbo", 1992, "Fighting", "SNES", "Capcom", "Twelve world warriors compete in the ultimate fighting tournament. Hadoken! Shoryuken! Sonic Boom!", 5, 2.49, 3.49},
		{"Super Mario Kart", 1992, "Racing", "SNES", "Nintendo", "Mario and friends race go-karts across colorful tracks. Red shells, banana peels, and fierce rivalries.", 5, 2.49, 3.49},
		{"EarthBound", 1994, "RPG", "SNES", "Nintendo", "Ness and his friends journey across Eagleland to stop the cosmic menace Giygas. Fuzzy pickles!", 2, 2.99, 3.99},
		{"Sonic the Hedgehog 2", 1992, "Platformer", "Genesis", "Sega", "Sonic and Tails race through zones to stop Dr. Robotnik's Death Egg. Introduces the Spin Dash.", 5, 2.49, 3.49},
		{"Streets of Rage 2", 1992, "Fighting", "Genesis", "Sega", "Axel, Blaze, Skate, and Max clear the streets of Mr. X's syndicate with devastating combos.", 4, 2.49, 3.49},
		{"Gunstar Heroes", 1993, "Action", "Genesis", "Treasure", "Free-form run-and-gun where weapon combos create unique attacks. Seven Force awaits at the end.", 3, 2.49, 3.49},
		{"Mortal Kombat II", 1994, "Fighting", "Genesis", "Midway", "Outworld's finest compete in the deadliest tournament. FINISH HIM! Scorpion, Sub-Zero, Kitana, and more.", 4, 2.49, 3.49},
		{"Final Fantasy VII", 1997, "RPG", "PS1", "Square", "Cloud Strife joins AVALANCHE to take down Shinra and the legendary Sephiroth. Materia, chocobos, and the Gold Saucer.", 4, 3.49, 4.49},
		{"Metal Gear Solid", 1998, "Action", "PS1", "Konami", "Solid Snake infiltrates Shadow Moses Island to stop FOXHOUND from launching a nuclear strike.", 3, 3.49, 4.49},
		{"Crash Bandicoot 2", 1997, "Platformer", "PS1", "Naughty Dog", "Crash collects crystals to stop Dr. Cortex from brainwashing the world. Wumpa fruit, TNT crates, and time trials.", 5, 2.99, 3.99},
		{"Resident Evil 2", 1998, "Horror", "PS1", "Capcom", "Leon and Claire survive the zombie outbreak in Raccoon City. Two scenarios, Mr. X, and plenty of ink ribbons.", 4, 2.99, 3.99},
		{"Tony Hawk's Pro Skater 2", 2000, "Sports", "PS1", "Activision", "The definitive skateboarding game. Manuals, reverts, and an iconic soundtrack. Goldfinger approved.", 5, 2.49, 3.49},
		{"Castlevania: SOTN", 1997, "Action", "PS1", "Konami", "Alucard explores Dracula's castle in this genre-defining RPG-platformer. What is a man? A miserable little pile of secrets.", 3, 3.49, 4.49},
		{"Tekken 3", 1998, "Fighting", "PS1", "Namco", "Twenty-three fighters compete in the King of Iron Fist Tournament 3. Jin Kazama, Eddy Gordo, and Gon.", 5, 2.49, 3.49},
		{"Super Mario 64", 1996, "Platformer", "N64", "Nintendo", "Mario jumps into paintings to collect Power Stars and rescue Peach from Bowser in the first 3D Mario adventure.", 4, 3.49, 4.49},
		{"The Legend of Zelda: Ocarina of Time", 1998, "Action", "N64", "Nintendo", "Link travels through time to stop Ganondorf. Master Sword, ocarina songs, and Hyrule at its finest.", 4, 3.49, 4.49},
		{"GoldenEye 007", 1997, "FPS", "N64", "Rare", "James Bond battles through 20 missions. Oddjob, proximity mines in Facility, and legendary split-screen multiplayer.", 5, 2.99, 3.99},
		{"Mario Kart 64", 1996, "Racing", "N64", "Nintendo", "Kart racing chaos with Mario and friends. Blue shells on Rainbow Road. Luigi's death stare.", 5, 2.99, 3.99},
		{"Super Smash Bros.", 1999, "Fighting", "N64", "Nintendo", "Nintendo's all-stars battle it out in platform-fighter mayhem. Pikachu vs Kirby vs Link vs Samus.", 5, 2.99, 3.99},
		{"Banjo-Kazooie", 1998, "Platformer", "N64", "Rare", "Bear and bird save Banjo's sister Tooty from the witch Gruntilda. Jiggies, Mumbo tokens, and Jinjos galore.", 4, 2.99, 3.99},
		{"Doom", 1993, "FPS", "PC", "id Software", "A space marine battles demons from Hell on Mars. The game that defined the first-person shooter genre.", 5, 2.49, 3.49},
		{"Half-Life", 1998, "FPS", "PC", "Valve", "Gordon Freeman fights through the Black Mesa Research Facility after a resonance cascade opens a portal to another dimension.", 4, 3.49, 4.49},
		{"StarCraft", 1998, "Strategy", "PC", "Blizzard", "Terran, Zerg, and Protoss clash in the Koprulu Sector. You must construct additional pylons.", 5, 2.99, 3.99},
		{"Age of Empires II", 1999, "Strategy", "PC", "Microsoft", "Build your civilization from the Dark Age to the Imperial Age. Wololo your enemies into submission.", 5, 2.99, 3.99},
		{"Diablo II", 2000, "RPG", "PC", "Blizzard", "Slay demons across five acts. Loot, level up, and chase that perfect Stone of Jordan. Stay a while and listen.", 4, 3.49, 4.49},
		{"The Sims", 2000, "Strategy", "PC", "Maxis", "Build homes, manage lives, and accidentally remove the pool ladder. The world's most popular life simulation.", 5, 2.99, 3.99},
		{"RollerCoaster Tycoon 2", 2002, "Strategy", "PC", "Infogrames", "Design and manage the ultimate theme park. Just keep building roller coasters and don't look at the balance sheet.", 4, 2.49, 3.49},
		{"Pac-Man", 1980, "Puzzle", "Arcade", "Namco", "Navigate mazes eating dots while avoiding ghosts. The highest-grossing arcade game of all time.", 3, 1.49, 1.99},
		{"Space Invaders", 1978, "Action", "Arcade", "Taito", "Defend Earth from descending alien waves. The game that launched a thousand arcades.", 3, 1.49, 1.99},
		{"Galaga", 1981, "Action", "Arcade", "Namco", "Fight waves of alien insects in deep space. Let them capture your ship — then free it for double the firepower.", 4, 1.49, 1.99},
		{"Street Fighter II", 1991, "Fighting", "Arcade", "Capcom", "The arcade classic that started the fighting game revolution. Ryu, Chun-Li, Guile, and the legendary combos.", 4, 2.49, 3.49},
		{"Donkey Kong", 1981, "Platformer", "Arcade", "Nintendo", "Jump over barrels to rescue Pauline from the mighty ape. Introduced both Mario and Donkey Kong to the world.", 3, 1.49, 1.99},
	}

	for _, g := range games {
		id := fmt.Sprintf("seed-movie-%s", sanitizeID(g.Title))
		movie := &models.Movie{
			ID:              id,
			Title:           g.Title,
			Year:            g.Year,
			Genre:           g.Genre,
			Format:          g.Platform + " Cartridge",
			Director:        g.Director,
			Synopsis:        g.Synopsis,
			CopiesTotal:     g.CopiesTotal,
			CopiesAvailable: g.CopiesTotal,
			Available:       true,
			RentalPrice:     g.RentalPrice,
			PlayPrice:       g.PlayPrice,
			Platform:        g.Platform,
			MediaType:       "game",
			Rating:          3.5 + rand.Float64()*1.5,
			RatingCount:     50 + rand.Intn(3000),
			CreatedAt:       time.Now().Unix(),
		}
		s.CreateMovie(movie)
	}
}
