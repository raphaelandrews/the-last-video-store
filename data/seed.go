package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

var (
	genres    = []string{"Action", "Comedy", "Horror", "SciFi", "Drama", "Thriller", "Romance", "Animation"}
	directors = []string{"Wachowski", "Tarantino", "Fincher", "Spielberg", "Nolan", "Scorsese", "Coppola", "Kubrick", "Scott", "Cameron"}
)

func main() {
	cfgPath := "thelastvideostore.db"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	os.Remove(cfgPath)

	s, err := store.Open(cfgPath)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	seedUsers(s)
	seedMovies(s)
	seedMerch(s)
	fmt.Printf("Seeded %d movies, 8 users, and 7 merch items.\n", countMovies(s))
}

func seedUsers(s *store.Store) {
	entries := []struct {
		name, pass, sub string
		tier            bitmask.Permission
		banned          bool
		balance         float64
	}{
		{"bronze", "123", "bronze", bitmask.TierBronze, false, 50},
		{"silver", "123", "silver", bitmask.TierSilver, false, 50},
		{"gold", "123", "gold", bitmask.TierGold, false, 50},
		{"employee", "123", "gold", bitmask.TierEmployee, false, 50},
		{"supervisor", "123", "gold", bitmask.TierSupervisor, false, 50},
		{"manager", "123", "diamond", bitmask.TierManager, false, 100},
		{"owner", "123", "diamond", bitmask.TierOwner, false, 100},
		{"banned", "123", "wood", bitmask.TierBronze, true, 5},
	}

	for _, e := range entries {
		hash, _ := auth.HashPassword(e.pass)
		now := time.Now().Unix()
		tier := models.TierByName(e.sub)
		user := &models.User{
			ID:            fmt.Sprintf("seed-%s", e.name),
			Username:      e.name,
			PasswordHash:  hash,
			Tier:          e.tier,
			MaxRentals:    tier.MaxConcurrent,
			Banned:        e.banned,
			PopcornPoints: 250,
			FreeRentals:   tier.FreeRentals,
			Balance:       e.balance,
			Subscription:  e.sub,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		s.CreateUser(user)
	}

	_ = crypto.New()
}

func seedMovies(s *store.Store) {
	movies := []struct {
		Title        string
		Year         int
		Genre        string
		Format       string
		Director     string
		Cast         []string
		Synopsis     string
		CopiesTotal  int
		IsNewRelease bool
	}{
		{"The Matrix", 1999, "SciFi", "DVD", "Wachowski", []string{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss"}, "A computer hacker learns the true nature of reality and his role in the war against its controllers.", 5, false},
		{"Pulp Fiction", 1994, "Drama", "VHS", "Tarantino", []string{"John Travolta", "Samuel L. Jackson", "Uma Thurman"}, "The lives of two mob hitmen, a boxer, a gangster and his wife intertwine in four tales.", 3, false},
		{"Fight Club", 1999, "Drama", "DVD", "Fincher", []string{"Brad Pitt", "Edward Norton", "Helena Bonham Carter"}, "An insomniac office worker and a devil-may-care soap maker form an underground fight club.", 4, false},
		{"Jurassic Park", 1993, "Action", "Blu-ray", "Spielberg", []string{"Sam Neill", "Laura Dern", "Jeff Goldblum"}, "A pragmatic paleontologist touring an almost complete theme park on an island in Central America is tasked with protecting a couple of kids.", 5, false},
		{"Inception", 2010, "SciFi", "DVD", "Nolan", []string{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Elliot Page"}, "A thief who steals corporate secrets through the use of dream-sharing technology.", 4, true},
		{"The Godfather", 1972, "Drama", "VHS", "Coppola", []string{"Marlon Brando", "Al Pacino", "James Caan"}, "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant youngest son.", 3, false},
		{"Goodfellas", 1990, "Drama", "DVD", "Scorsese", []string{"Robert De Niro", "Ray Liotta", "Joe Pesci"}, "The story of Henry Hill and his life in the mafia, covering his relationship with his wife and his mob partners.", 4, false},
		{"Forrest Gump", 1994, "Drama", "VHS", "Zemeckis", []string{"Tom Hanks", "Robin Wright", "Gary Sinise"}, "The presidencies of Kennedy and Johnson through the eyes of an Alabama man with an IQ of 75.", 5, false},
		{"The Dark Knight", 2008, "Action", "Blu-ray", "Nolan", []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart"}, "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological tests.", 5, false},
		{"Back to the Future", 1985, "SciFi", "VHS", "Zemeckis", []string{"Michael J. Fox", "Christopher Lloyd", "Lea Thompson"}, "Marty McFly, a 17-year-old high school student, is accidentally sent 30 years into the past in a time-traveling DeLorean.", 4, false},
		{"Die Hard", 1988, "Action", "DVD", "McTiernan", []string{"Bruce Willis", "Alan Rickman", "Bonnie Bedelia"}, "NYPD officer John McClane tries to save his wife and several others taken hostage by terrorists.", 4, false},
		{"Terminator 2", 1991, "SciFi", "Blu-ray", "Cameron", []string{"Arnold Schwarzenegger", "Linda Hamilton", "Edward Furlong"}, "A cyborg must protect John Connor from a more advanced and powerful Terminator.", 3, false},
		{"Toy Story", 1995, "Animation", "DVD", "Lasseter", []string{"Tom Hanks", "Tim Allen", "Don Rickles"}, "A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room.", 5, false},
		{"The Silence of the Lambs", 1991, "Horror", "VHS", "Demme", []string{"Jodie Foster", "Anthony Hopkins", "Lawrence A. Bonney"}, "A young F.B.I. cadet must receive the help of an incarcerated and manipulative cannibal killer to help catch another serial killer.", 3, false},
		{"Saving Private Ryan", 1998, "Action", "DVD", "Spielberg", []string{"Tom Hanks", "Matt Damon", "Tom Sizemore"}, "Following the Normandy Landings, a group of U.S. soldiers go behind enemy lines to retrieve a paratrooper whose brothers have been killed.", 4, false},
		{"Gladiator", 2000, "Action", "Blu-ray", "Scott", []string{"Russell Crowe", "Joaquin Phoenix", "Connie Nielsen"}, "A former Roman General sets out to exact vengeance against the corrupt emperor who murdered his family and sent him into slavery.", 4, false},
		{"Gladiator Extended", 2005, "Action", "Blu-ray", "Scott", []string{"Russell Crowe", "Joaquin Phoenix", "Connie Nielsen"}, "Extended cut of the epic historical drama with additional scenes.", 2, true},
		{"Kill Bill Vol.1", 2003, "Action", "DVD", "Tarantino", []string{"Uma Thurman", "Lucy Liu", "Vivica A. Fox"}, "A former assassin wakes from a coma and seeks revenge against the team of assassins who betrayed her.", 4, false},
		{"The Departed", 2006, "Thriller", "Blu-ray", "Scorsese", []string{"Leonardo DiCaprio", "Matt Damon", "Jack Nicholson"}, "An undercover cop and a mole in the police attempt to identify each other while infiltrating an Irish gang.", 4, false},
		{"No Country for Old Men", 2007, "Thriller", "DVD", "Coen", []string{"Tommy Lee Jones", "Javier Bardem", "Josh Brolin"}, "Violence and mayhem ensue after a hunter stumbles upon a drug deal gone wrong and more than two million dollars in cash.", 2, false},
		{"There Will Be Blood", 2007, "Drama", "VHS", "Anderson", []string{"Daniel Day-Lewis", "Paul Dano", "Ciaran Hinds"}, "A story of family, religion, hatred, oil and madness, focusing on a turn-of-the-century prospector.", 2, false},
		{"Inglourious Basterds", 2009, "Action", "DVD", "Tarantino", []string{"Brad Pitt", "Diane Kruger", "Eli Roth"}, "In Nazi-occupied France during World War II, a plan to assassinate Nazi leaders by a group of Jewish U.S. soldiers coincides with a theatre owner's vengeful plans.", 4, false},
		{"Blade Runner", 1982, "SciFi", "VHS", "Scott", []string{"Harrison Ford", "Rutger Hauer", "Sean Young"}, "A blade runner must pursue and terminate four replicants who stole a ship in space and have returned to Earth to find their creator.", 3, false},
		{"American Beauty", 1999, "Drama", "DVD", "Mendes", []string{"Kevin Spacey", "Annette Bening", "Thora Birch"}, "A sexually frustrated suburban father has a mid-life crisis after becoming infatuated with his daughter's best friend.", 3, false},
		{"Spirited Away", 2001, "Animation", "Blu-ray", "Miyazaki", []string{"Rumi Hiiragi", "Miyu Irino", "Mari Natsuki"}, "A sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits.", 4, false},
		{"City of God", 2002, "Drama", "DVD", "Meirelles", []string{"Alexandre Rodrigues", "Leandro Firmino", "Phellipe Haagensen"}, "In the slums of Rio, two kids' paths diverge as one struggles to become a photographer and the other a kingpin.", 3, false},
		{"District 9", 2009, "SciFi", "DVD", "Blomkamp", []string{"Sharlto Copley", "Jason Cope", "Nathalie Boltt"}, "Violence ensues after an extraterrestrial race forced to live in slum-like conditions on Earth.", 3, false},
		{"2001: A Space Odyssey", 1968, "SciFi", "Blu-ray", "Kubrick", []string{"Keir Dullea", "Gary Lockwood", "William Sylvester"}, "Humanity finds a mysterious object buried beneath the lunar surface and sets off on a quest to find its origins.", 3, false},
		{"A Clockwork Orange", 1971, "Drama", "VHS", "Kubrick", []string{"Malcolm McDowell", "Patrick Magee", "Michael Bates"}, "In the future, a sadistic gang leader is imprisoned and volunteers for a conduct-aversion experiment.", 3, false},
		{"The Shining", 1980, "Horror", "DVD", "Kubrick", []string{"Jack Nicholson", "Shelley Duvall", "Danny Lloyd"}, "A family heads to an isolated hotel for the winter where a sinister presence influences the father into violence.", 4, false},
		{"Apocalypse Now", 1979, "Action", "Blu-ray", "Coppola", []string{"Martin Sheen", "Marlon Brando", "Robert Duvall"}, "A U.S. Army officer serving in Vietnam is tasked with assassinating a renegade Special Forces Colonel.", 3, false},
		{"Taxi Driver", 1976, "Drama", "VHS", "Scorsese", []string{"Robert De Niro", "Jodie Foster", "Cybill Shepherd"}, "A mentally unstable veteran works as a nighttime taxi driver in New York City.", 4, false},
		{"The Wolf of Wall Street", 2013, "Comedy", "Blu-ray", "Scorsese", []string{"Leonardo DiCaprio", "Jonah Hill", "Margot Robbie"}, "Based on the true story of Jordan Belfort, from his rise to a wealthy stock-broker living the high life to his fall involving crime, corruption and the federal government.", 4, false},
		{"Alien", 1979, "Horror", "DVD", "Scott", []string{"Sigourney Weaver", "Tom Skerritt", "John Hurt"}, "The crew of a commercial spacecraft encounter a deadly lifeform after investigating an unknown transmission.", 4, false},
		{"Aliens", 1986, "Action", "Blu-ray", "Cameron", []string{"Sigourney Weaver", "Michael Biehn", "Carrie Henn"}, "Fifty-seven years after surviving an apocalyptic attack aboard her space vessel, Ellen Ripley is called back to the planet LV-426.", 4, false},
		{"Raiders of the Lost Ark", 1981, "Action", "DVD", "Spielberg", []string{"Harrison Ford", "Karen Allen", "Paul Freeman"}, "In 1936, archaeologist Indiana Jones is hired by the U.S. government to locate the Ark of the Covenant.", 4, false},
		{"Jaws", 1975, "Thriller", "VHS", "Spielberg", []string{"Roy Scheider", "Robert Shaw", "Richard Dreyfuss"}, "When a killer shark unleashes chaos on a beach community off Cape Cod, it's up to a local sheriff, a marine biologist, and an old seafarer to hunt the beast down.", 5, false},
		{"E.T.", 1982, "SciFi", "VHS", "Spielberg", []string{"Henry Thomas", "Drew Barrymore", "Peter Coyote"}, "A troubled child summons the courage to help a friendly alien escape Earth and return to his home world.", 4, false},
		{"The Thing", 1982, "Horror", "DVD", "Carpenter", []string{"Kurt Russell", "Wilford Brimley", "Keith David"}, "A research team in Antarctica is hunted by a shape-shifting alien that assumes the appearance of its victims.", 3, false},
		{"Halloween", 1978, "Horror", "VHS", "Carpenter", []string{"Donald Pleasence", "Jamie Lee Curtis", "Nancy Kyes"}, "Fifteen years after murdering his sister on Halloween night 1963, Michael Myers escapes from a mental hospital and returns to the small town of Haddonfield.", 4, false},
		{"The Exorcist", 1973, "Horror", "DVD", "Friedkin", []string{"Ellen Burstyn", "Max von Sydow", "Linda Blair"}, "When a 12-year-old girl is possessed by a mysterious entity, her mother seeks the help of two priests.", 3, false},
		{"Reservoir Dogs", 1992, "Thriller", "VHS", "Tarantino", []string{"Harvey Keitel", "Tim Roth", "Michael Madsen"}, "When a simple jewelry heist goes horribly wrong, the surviving criminals suspect that one of them is a police informant.", 3, false},
		{"Braveheart", 1995, "Action", "Blu-ray", "Gibson", []string{"Mel Gibson", "Sophie Marceau", "Patrick McGoohan"}, "Scottish warrior William Wallace leads his countrymen in a rebellion to free his homeland from the tyranny of King Edward I.", 4, false},
		{"Casino", 1995, "Drama", "DVD", "Scorsese", []string{"Robert De Niro", "Sharon Stone", "Joe Pesci"}, "A tale of greed, deception, money, power, and murder occur between two mobster best friends and a trophy wife over a gambling empire.", 3, false},
		{"Interstellar", 2014, "SciFi", "Blu-ray", "Nolan", []string{"Matthew McConaughey", "Anne Hathaway", "Jessica Chastain"}, "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.", 5, false},
		{"The Avengers", 2012, "Action", "Blu-ray", "Whedon", []string{"Robert Downey Jr.", "Chris Evans", "Scarlett Johansson"}, "Earth's mightiest heroes must come together to stop Loki and his alien army from enslaving humanity.", 5, false},
		{"The Social Network", 2010, "Drama", "DVD", "Fincher", []string{"Jesse Eisenberg", "Andrew Garfield", "Justin Timberlake"}, "As Harvard student Mark Zuckerberg creates the social networking site that would become Facebook, he is sued by two brothers who claimed he stole their idea.", 4, false},
		{"Finding Nemo", 2003, "Animation", "DVD", "Stanton", []string{"Albert Brooks", "Ellen DeGeneres", "Alexander Gould"}, "After his son is captured, a timid clownfish sets out on a journey across the ocean to bring him home.", 5, false},
		{"The Incredibles", 2004, "Animation", "DVD", "Bird", []string{"Craig T. Nelson", "Holly Hunter", "Samuel L. Jackson"}, "While trying to lead a quiet suburban life, a family of undercover superheroes are forced into action to save the world.", 4, false},
		{"Shrek", 2001, "Animation", "DVD", "Adamson", []string{"Mike Myers", "Eddie Murphy", "Cameron Diaz"}, "A mean lord exiles fairytale creatures to the swamp of a grumpy ogre, who must go on a quest and rescue a princess.", 5, false},
		{"The Lion King", 1994, "Animation", "VHS", "Allers", []string{"Matthew Broderick", "Jeremy Irons", "James Earl Jones"}, "Lion prince Simba flees his kingdom after the murder of his father, only to learn the true meaning of responsibility and bravery.", 5, false},
		{"Predator", 1987, "Action", "DVD", "McTiernan", []string{"Arnold Schwarzenegger", "Carl Weathers", "Kevin Peter Hall"}, "A team of commandos on a mission in a Central American jungle find themselves hunted by an extraterrestrial warrior.", 3, false},
		{"Ghostbusters", 1984, "Comedy", "VHS", "Reitman", []string{"Bill Murray", "Dan Aykroyd", "Sigourney Weaver"}, "Three former parapsychology professors set up shop as a unique ghost removal service.", 4, false},
		{"Ferris Buellers Day Off", 1986, "Comedy", "VHS", "Hughes", []string{"Matthew Broderick", "Alan Ruck", "Mia Sara"}, "A charismatic high-school student convinces his friends to play hooky and spend one epic day in downtown Chicago.", 3, false},
		{"The Breakfast Club", 1985, "Comedy", "VHS", "Hughes", []string{"Emilio Estevez", "Judd Nelson", "Molly Ringwald"}, "Five high school students meet in Saturday detention and discover how they have a lot more in common than they thought.", 3, false},
		{"Home Alone", 1990, "Comedy", "VHS", "Columbus", []string{"Macaulay Culkin", "Joe Pesci", "Daniel Stern"}, "An eight-year-old troublemaker must protect his house from a pair of burglars when he is accidentally left home alone.", 4, false},
		{"Gremlins", 1984, "Comedy", "VHS", "Dante", []string{"Zach Galligan", "Phoebe Cates", "Hoyt Axton"}, "A young man inadvertently breaks three important rules concerning his new pet and unleashes a horde of malevolently mischievous monsters.", 3, false},
		{"Rocky", 1976, "Drama", "VHS", "Avildsen", []string{"Sylvester Stallone", "Talia Shire", "Burt Young"}, "A small-time Philadelphia boxer gets a supremely rare chance to fight the world heavyweight champion.", 4, false},
		{"The Fifth Element", 1997, "SciFi", "DVD", "Besson", []string{"Bruce Willis", "Milla Jovovich", "Gary Oldman"}, "In the colorful future, a cab driver unwittingly becomes the central figure in the search for a legendary cosmic weapon.", 3, false},
		{"Shaun of the Dead", 2004, "Comedy", "DVD", "Wright", []string{"Simon Pegg", "Nick Frost", "Kate Ashfield"}, "The uneventful life of a London electronics salesman is disrupted by a zombie apocalypse.", 3, false},
		{"Dead Poets Society", 1989, "Drama", "VHS", "Weir", []string{"Robin Williams", "Ethan Hawke", "Robert Sean Leonard"}, "Maverick teacher John Keating uses poetry to embolden his boarding school students to new heights of self-expression.", 3, false},
		{"A Beautiful Mind", 2001, "Drama", "DVD", "Howard", []string{"Russell Crowe", "Ed Harris", "Jennifer Connelly"}, "After John Nash, a brilliant but asocial mathematical genius, accepts secret work in cryptography, his life takes a turn for the nightmarish.", 3, false},
		{"Minority Report", 2002, "SciFi", "DVD", "Spielberg", []string{"Tom Cruise", "Colin Farrell", "Samantha Morton"}, "In a future where a special police unit is able to arrest murderers before they commit their crimes, an officer from that unit is himself accused of a future murder.", 3, false},
		{"Harry Potter and the Philosophers Stone", 2001, "Action", "Blu-ray", "Columbus", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "An orphaned boy enrolls in a school of wizardry, where he learns the truth about himself, his family, and the terrible evil that haunts the magical world.", 4, false},
		{"Harry Potter and the Chamber of Secrets", 2002, "Action", "Blu-ray", "Columbus", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "Harry ignores warnings not to return to Hogwarts, only to find the school plagued by a series of mysterious attacks.", 4, false},
		{"Harry Potter and the Prisoner of Azkaban", 2004, "Action", "Blu-ray", "Cuaron", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "Harry must confront soul-sucking Dementors, outsmart a werewolf and learn the truth about the escaped Sirius Black.", 4, false},
		{"Harry Potter and the Goblet of Fire", 2005, "Action", "Blu-ray", "Newell", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "Harry finds himself mysteriously selected as an under-aged competitor in a dangerous tournament between three schools.", 3, false},
		{"Harry Potter and the Order of the Phoenix", 2007, "Action", "Blu-ray", "Yates", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "With their warning about Lord Voldemort's return scoffed at, Harry and Dumbledore are targeted as the Ministry tightens its grip.", 3, false},
		{"Harry Potter and the Half-Blood Prince", 2009, "Action", "Blu-ray", "Yates", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "As Harry Potter begins his sixth year at Hogwarts, he discovers an old book marked as the property of the Half-Blood Prince.", 3, false},
		{"Harry Potter and the Deathly Hallows Part 1", 2010, "Action", "Blu-ray", "Yates", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "As Harry races against time and evil to destroy the Horcruxes, he uncovers the existence of three powerful objects.", 3, false},
		{"Harry Potter and the Deathly Hallows Part 2", 2011, "Action", "Blu-ray", "Yates", []string{"Daniel Radcliffe", "Emma Watson", "Rupert Grint"}, "Harry, Ron and Hermione search for Voldemort's remaining Horcruxes in their effort to destroy the Dark Lord.", 3, false},
		{"Predator 2", 1990, "Action", "DVD", "Hopkins", []string{"Danny Glover", "Gary Busey", "Kevin Peter Hall"}, "The Predator returns to Earth, this time to stake a claim on the war-torn streets of a dystopian Los Angeles.", 3, false},
		{"Star Wars: A New Hope", 1977, "SciFi", "Blu-ray", "Lucas", []string{"Mark Hamill", "Harrison Ford", "Carrie Fisher"}, "Luke Skywalker joins forces with a Jedi Knight, a cocky pilot, a Wookiee and two droids to save the galaxy.", 5, false},
		{"Star Wars: The Empire Strikes Back", 1980, "SciFi", "Blu-ray", "Kershner", []string{"Mark Hamill", "Harrison Ford", "Carrie Fisher"}, "After the Rebels are overpowered, Luke Skywalker begins Jedi training with Yoda while his friends are pursued by Darth Vader.", 5, false},
		{"Star Wars: Return of the Jedi", 1983, "SciFi", "Blu-ray", "Marquand", []string{"Mark Hamill", "Harrison Ford", "Carrie Fisher"}, "After rescuing Han Solo from Jabba the Hutt, the Rebels attempt to destroy the second Death Star.", 5, false},
		{"Star Wars: The Phantom Menace", 1999, "SciFi", "DVD", "Lucas", []string{"Liam Neeson", "Ewan McGregor", "Natalie Portman"}, "Two Jedi Knights escape a hostile blockade to find allies and come across a young boy who may bring balance to the Force.", 4, false},
		{"Star Wars: Attack of the Clones", 2002, "SciFi", "DVD", "Lucas", []string{"Ewan McGregor", "Natalie Portman", "Hayden Christensen"}, "Ten years after initially meeting, Anakin Skywalker shares a forbidden romance with Padmé Amidala.", 4, false},
		{"Star Wars: Revenge of the Sith", 2005, "SciFi", "DVD", "Lucas", []string{"Ewan McGregor", "Natalie Portman", "Hayden Christensen"}, "Anakin succumbs to the dark side, becoming Darth Vader, as the Jedi are purged and the Empire rises.", 4, false},
		{"Spirited Away", 2001, "Animation", "Blu-ray", "Miyazaki", []string{"Rumi Hiiragi", "Miyu Irino", "Mari Natsuki"}, "During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits.", 4, false},
		{"My Neighbor Totoro", 1988, "Animation", "DVD", "Miyazaki", []string{"Noriko Hidaka", "Chika Sakamoto", "Hitoshi Takagi"}, "When two girls move to the country, they befriend the magical creatures that inhabit the nearby forest.", 5, false},
		{"Princess Mononoke", 1997, "Animation", "DVD", "Miyazaki", []string{"Yoji Matsuda", "Yuriko Ishida", "Yuko Tanaka"}, "A prince infected with a lethal curse embarks on a journey to find a cure and lands in the middle of a battle between a mining town and forest gods.", 3, false},
		{"Howl's Moving Castle", 2004, "Animation", "Blu-ray", "Miyazaki", []string{"Chieko Baisho", "Takuya Kimura", "Akihiro Miwa"}, "When an unconfident young woman is cursed with an old body by a spiteful witch, her only chance of breaking the spell lies with a self-indulgent wizard.", 3, false},
		{"Kikis Delivery Service", 1989, "Animation", "DVD", "Miyazaki", []string{"Minami Takayama", "Rei Sakuma", "Kappei Yamaguchi"}, "A young witch, on her mandatory year of independent life, finds fitting into a new community difficult.", 4, false},
		{"Grave of the Fireflies", 1988, "Animation", "DVD", "Takahata", []string{"Tsutomu Tatsumi", "Ayano Shiraishi", "Yoshiko Shinohara"}, "A young boy and his little sister struggle to survive in Japan during World War II.", 2, false},
		{"Castle in the Sky", 1986, "Animation", "DVD", "Miyazaki", []string{"Mayumi Tanaka", "Keiko Yokozawa", "Kotoe Hatsui"}, "A young boy and a girl with a magic crystal must race against pirates and foreign agents for a floating castle.", 3, false},
		{"Nausicaa of the Valley of the Wind", 1984, "Animation", "DVD", "Miyazaki", []string{"Sumi Shimamoto", "Mahito Tsujimura", "Hisako Kyoda"}, "Warrior and pacifist Princess Nausicaa desperately struggles to prevent two warring nations from destroying themselves and their dying planet.", 3, false},
		{"Porco Rosso", 1992, "Animation", "DVD", "Miyazaki", []string{"Shuichiro Moriyama", "Tokiko Kato", "Akemi Okamura"}, "In 1930s Italy, a veteran World War I pilot is cursed to look like an anthropomorphic pig.", 3, false},
		{"Ponyo", 2008, "Animation", "Blu-ray", "Miyazaki", []string{"Tomoko Yamaguchi", "Kazushige Nagashima", "Yuki Amami"}, "A five-year-old boy develops a relationship with Ponyo, a young goldfish princess who longs to become human.", 4, false},
		{"The End of Evangelion", 1997, "Animation", "DVD", "Anno", []string{"Megumi Ogata", "Megumi Hayashibara", "Yuko Miyamura"}, "Concurrent theatrical ending to the TV series Neon Genesis Evangelion, depicting the apocalyptic final battle.", 2, false},
		{"Pokemon: The First Movie", 1998, "Animation", "VHS", "Yuyama", []string{"Veronica Taylor", "Rachael Lillis", "Eric Stuart"}, "Scientists genetically create a new Pokemon, Mewtwo, but the results are horrific and disastrous.", 3, false},
		{"Pokemon 2000", 1999, "Animation", "VHS", "Yuyama", []string{"Veronica Taylor", "Rachael Lillis", "Eric Stuart"}, "Ash must gather three spheres from the islands to save the world from a catastrophic weather disaster.", 3, false},
		{"Ice Age", 2002, "Animation", "DVD", "Wedge", []string{"Ray Romano", "John Leguizamo", "Denis Leary"}, "Set during the Ice Age, a sabertooth tiger, a sloth, and a wooly mammoth find a lost human infant, and they try to return him to his tribe.", 4, false},
		{"Ice Age: The Meltdown", 2006, "Animation", "DVD", "Saldanha", []string{"Ray Romano", "John Leguizamo", "Denis Leary"}, "Manny, Sid, and Diego discover that the ice age is coming to an end, and join everybody for a journey to higher ground.", 4, false},
		{"Ice Age: Dawn of the Dinosaurs", 2009, "Animation", "Blu-ray", "Saldanha", []string{"Ray Romano", "John Leguizamo", "Denis Leary"}, "When Sid's attempt to adopt three dinosaur eggs gets him abducted by their real mother, Manny and Diego set off to rescue him.", 4, false},
		{"Ghost in the Shell", 1995, "Animation", "DVD", "Oshii", []string{"Atsuko Tanaka", "Akio Otsuka", "Iemasa Kayumi"}, "A cyborg policewoman hunts a mysterious hacker called the Puppet Master.", 2, false},
		{"Lagaan", 2001, "Drama", "DVD", "Gowariker", []string{"Aamir Khan", "Gracy Singh", "Rachel Shelley"}, "The people of a small village in Victorian India stake their future on a game of cricket against their ruthless British rulers.", 2, false},
		{"3 Idiots", 2009, "Comedy", "Blu-ray", "Hirani", []string{"Aamir Khan", "Madhavan", "Sharman Joshi"}, "Two friends are searching for their long lost companion. They revisit their college days and recall the memories of their friend who inspired them.", 3, false},
		{"Rang De Basanti", 2006, "Drama", "DVD", "Mehra", []string{"Aamir Khan", "Soha Ali Khan", "Siddharth"}, "The story of six young Indians who assist an English woman to film a documentary on the freedom fighters from their past.", 2, false},
		{"American Pie", 1999, "Comedy", "VHS", "Weitz", []string{"Jason Biggs", "Chris Klein", "Thomas Ian Nicholas"}, "Four teenage boys enter a pact to lose their virginity by prom night.", 3, false},
		{"American Pie 2", 2001, "Comedy", "DVD", "Rogers", []string{"Jason Biggs", "Seann William Scott", "Thomas Ian Nicholas"}, "The whole gang reunites at a beach house for a wild summer of parties.", 3, false},
		{"American Wedding", 2003, "Comedy", "DVD", "Dylan", []string{"Jason Biggs", "Alyson Hannigan", "Seann William Scott"}, "Jim and Michelle are getting married, but Finch and Stifler threaten the upcoming nuptials.", 3, false},
		{"American Reunion", 2012, "Comedy", "DVD", "Hurwitz", []string{"Jason Biggs", "Alyson Hannigan", "Seann William Scott"}, "Jim, Michelle and their friends reunite in their hometown for their high school reunion.", 3, false},
		{"Munna Bhai MBBS", 2003, "Comedy", "DVD", "Hirani", []string{"Sanjay Dutt", "Arshad Warsi", "Gracy Singh"}, "A gangster sets out to fulfill his father's dream of becoming a doctor by enrolling in medical school.", 3, false},
		{"Lage Raho Munna Bhai", 2006, "Comedy", "DVD", "Hirani", []string{"Sanjay Dutt", "Arshad Warsi", "Vidya Balan"}, "Munna Bhai embarks on a journey with Mahatma Gandhi in order to fight against a corrupt property dealer.", 3, false},
		{"Mean Girls", 2004, "Comedy", "DVD", "Waters", []string{"Lindsay Lohan", "Rachel McAdams", "Tina Fey"}, "Cady Heron is a hit with a group known as The Plastics, but she soon learns how shallow they really are.", 4, false},
		{"Madagascar", 2005, "Animation", "DVD", "Darnell", []string{"Ben Stiller", "Chris Rock", "David Schwimmer"}, "A group of animals who have spent their lives in the Central Park Zoo find themselves stranded on Madagascar.", 4, false},
		{"10 Things I Hate About You", 1999, "Romance", "DVD", "Junger", []string{"Heath Ledger", "Julia Stiles", "Joseph Gordon-Levitt"}, "A high-school boy must find a date for the antisocial sister of the girl he wants to date.", 3, false},
		{"High School Musical", 2006, "Comedy", "DVD", "Ortega", []string{"Zac Efron", "Vanessa Hudgens", "Ashley Tisdale"}, "A popular high school athlete and an academically gifted girl get roles in the school musical, causing a rift among their cliques.", 4, false},
		{"High School Musical 2", 2007, "Comedy", "DVD", "Ortega", []string{"Zac Efron", "Vanessa Hudgens", "Ashley Tisdale"}, "School's out and the East High Wildcats are ready to enjoy summer at a country club.", 4, false},
		{"High School Musical 3", 2008, "Comedy", "DVD", "Ortega", []string{"Zac Efron", "Vanessa Hudgens", "Ashley Tisdale"}, "As seniors in high school, the Wildcats discover they are growing apart and decide to put on one last musical.", 3, false},
		{"Sleepy Hollow", 1999, "Horror", "DVD", "Burton", []string{"Johnny Depp", "Christina Ricci", "Miranda Richardson"}, "Ichabod Crane is sent to Sleepy Hollow to investigate the decapitations of three people, with the culprit being the legendary Headless Horseman.", 3, false},
		{"Pirates of the Caribbean: The Curse of the Black Pearl", 2003, "Action", "Blu-ray", "Verbinski", []string{"Johnny Depp", "Orlando Bloom", "Keira Knightley"}, "Blacksmith Will Turner teams up with eccentric pirate Jack Sparrow to save his love, the governor's daughter.", 4, false},
		{"Pirates of the Caribbean: Dead Mans Chest", 2006, "Action", "Blu-ray", "Verbinski", []string{"Johnny Depp", "Orlando Bloom", "Keira Knightley"}, "Jack Sparrow races to recover the heart of Davy Jones to avoid enslaving his soul to Jones' service.", 4, false},
		{"Pirates of the Caribbean: At Worlds End", 2007, "Action", "Blu-ray", "Verbinski", []string{"Johnny Depp", "Orlando Bloom", "Keira Knightley"}, "Captain Barbossa, Will Turner and Elizabeth Swann must sail off the edge of the map to rescue Captain Jack Sparrow.", 4, false},
		{"Kill Bill: Vol. 1", 2003, "Action", "DVD", "Tarantino", []string{"Uma Thurman", "Lucy Liu", "Vivica A. Fox"}, "After awakening from a four-year coma, a former assassin wreaks vengeance on the team of assassins who betrayed her.", 3, false},
		{"Kill Bill: Vol. 2", 2004, "Action", "DVD", "Tarantino", []string{"Uma Thurman", "David Carradine", "Michael Madsen"}, "The Bride continues her quest of vengeance against her former boss and lover Bill.", 3, false},
		{"Indiana Jones and the Raiders of the Lost Ark", 1981, "Action", "Blu-ray", "Spielberg", []string{"Harrison Ford", "Karen Allen", "Paul Freeman"}, "Archaeologist Indiana Jones races against Nazis to locate the legendary Ark of the Covenant.", 5, false},
		{"Indiana Jones and the Temple of Doom", 1984, "Action", "Blu-ray", "Spielberg", []string{"Harrison Ford", "Kate Capshaw", "Ke Huy Quan"}, "After arriving in India, Indiana Jones is asked by a village to find a mystical stone and rescue their children.", 4, false},
		{"Indiana Jones and the Last Crusade", 1989, "Action", "Blu-ray", "Spielberg", []string{"Harrison Ford", "Sean Connery", "Alison Doody"}, "Indiana Jones searches for his father, a Holy Grail scholar, who has been kidnapped by Nazis.", 4, false},
		{"Alien 3", 1992, "Horror", "DVD", "Fincher", []string{"Sigourney Weaver", "Charles S. Dutton", "Charles Dance"}, "After her last encounter, Ripley crash-lands on a maximum security prison planet and discovers an Alien was on board.", 3, false},
		{"Alien Resurrection", 1997, "Horror", "DVD", "Jeunet", []string{"Sigourney Weaver", "Winona Ryder", "Ron Perlman"}, "Two centuries after her death, Ripley is revived as a powerful human/alien hybrid clone.", 3, false},
		{"The Matrix Reloaded", 2003, "SciFi", "Blu-ray", "Wachowski", []string{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss"}, "Neo and the rebel leaders estimate that they have 72 hours until 250,000 probes discover Zion.", 4, false},
		{"The Matrix Revolutions", 2003, "SciFi", "Blu-ray", "Wachowski", []string{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss"}, "The human city of Zion defends itself against the massive invasion of the machines.", 4, false},
		{"The Godfather Part II", 1974, "Drama", "VHS", "Coppola", []string{"Al Pacino", "Robert De Niro", "Robert Duvall"}, "The early life and career of Vito Corleone in 1920s New York is portrayed while his son, Michael, expands the family crime syndicate.", 4, false},
		{"The Godfather Part III", 1990, "Drama", "DVD", "Coppola", []string{"Al Pacino", "Diane Keaton", "Andy Garcia"}, "Michael Corleone seeks to legitimize his family's business affairs and get out of the Mafia.", 3, false},
	}

	for _, m := range movies {
		id := fmt.Sprintf("seed-movie-%s", sanitizeID(m.Title))
		movie := &models.Movie{
			ID:              id,
			Title:           m.Title,
			Year:            m.Year,
			Genre:           m.Genre,
			Format:          m.Format,
			Director:        m.Director,
			Cast:            m.Cast,
			Synopsis:        m.Synopsis,
			CopiesTotal:     m.CopiesTotal,
			CopiesAvailable: m.CopiesTotal,
			Available:       true,
			IsNewRelease:    m.IsNewRelease,
			RentalPrice:     moviePrice(m.Year, m.Format, m.IsNewRelease),
			Rating:          3.0 + rand.Float64()*2.0,
			RatingCount:     100 + rand.Intn(5000),
			CreatedAt:       time.Now().Unix(),
		}
		s.CreateMovie(movie)
	}

	s.AddStaffPick("seed-movie-TheMatrix")
	s.AddStaffPick("seed-movie-TheDarkKnight")
	s.AddStaffPick("seed-movie-PulpFiction")

	seedSequels(s)
}

func seedSequels(s *store.Store) {
	links := map[string]string{
		"seed-movie-Aliens":                               "seed-movie-Alien",
		"seed-movie-Avengers":                             "seed-movie-IronMan",
		"seed-movie-HarryPotterandtheChamberofSecrets":    "seed-movie-HarryPotterandthePhilosophersStone",
		"seed-movie-HarryPotterandthePrisonerofAzkaban":   "seed-movie-HarryPotterandtheChamberofSecrets",
		"seed-movie-HarryPotterandtheGobletofFire":        "seed-movie-HarryPotterandthePrisonerofAzkaban",
		"seed-movie-HarryPotterandtheOrderofthePhoenix":   "seed-movie-HarryPotterandtheGobletofFire",
		"seed-movie-HarryPotterandtheHalfBloodPrince":     "seed-movie-HarryPotterandtheOrderofthePhoenix",
		"seed-movie-HarryPotterandtheDeathlyHallowsPart1": "seed-movie-HarryPotterandtheHalfBloodPrince",
		"seed-movie-HarryPotterandtheDeathlyHallowsPart2": "seed-movie-HarryPotterandtheDeathlyHallowsPart1",
		"seed-movie-Predator2":                            "seed-movie-Predator",
		"seed-movie-StarWarsTheEmpireStrikesBack":         "seed-movie-StarWarsANewHope",
		"seed-movie-StarWarsReturnoftheJedi":              "seed-movie-StarWarsTheEmpireStrikesBack",
		"seed-movie-StarWarsAttackoftheClones":            "seed-movie-StarWarsThePhantomMenace",
		"seed-movie-StarWarsRevengeoftheSith":             "seed-movie-StarWarsAttackoftheClones",
		"seed-movie-Pokemon2000":                          "seed-movie-PokemonTheFirstMovie",
		"seed-movie-IceAgeTheMeltdown":                    "seed-movie-IceAge",
		"seed-movie-IceAgeDawnoftheDinosaurs":             "seed-movie-IceAgeTheMeltdown",
		"seed-movie-AmericanPie2":                         "seed-movie-AmericanPie",
		"seed-movie-AmericanWedding":                      "seed-movie-AmericanPie2",
		"seed-movie-AmericanReunion":                      "seed-movie-AmericanWedding",
		"seed-movie-LageRahoMunnaBhai":                    "seed-movie-MunnaBhaiMBBS",
		"seed-movie-HighSchoolMusical2":                   "seed-movie-HighSchoolMusical",
		"seed-movie-HighSchoolMusical3":                   "seed-movie-HighSchoolMusical2",
		"seed-movie-PiratesoftheCaribbeanDeadMansChest":   "seed-movie-PiratesoftheCaribbeanTheCurseoftheBlackPearl",
		"seed-movie-PiratesoftheCaribbeanAtWorldsEnd":     "seed-movie-PiratesoftheCaribbeanDeadMansChest",
		"seed-movie-KillBillVol2":                         "seed-movie-KillBillVol1",
		"seed-movie-IndianaJonesandtheTempleofDoom":       "seed-movie-IndianaJonesandtheRaidersoftheLostArk",
		"seed-movie-IndianaJonesandtheLastCrusade":        "seed-movie-IndianaJonesandtheTempleofDoom",
		"seed-movie-Alien3":                               "seed-movie-Aliens",
		"seed-movie-AlienResurrection":                    "seed-movie-Alien3",
		"seed-movie-TheMatrixReloaded":                    "seed-movie-TheMatrix",
		"seed-movie-TheMatrixRevolutions":                 "seed-movie-TheMatrixReloaded",
		"seed-movie-TheGodfatherPartII":                   "seed-movie-TheGodfather",
		"seed-movie-TheGodfatherPartIII":                  "seed-movie-TheGodfatherPartII",
	}
	for movieID, prequelID := range links {
		m, err := s.GetMovieByID(movieID)
		if err == nil {
			m.SequelTo = prequelID
			s.UpdateMovie(m)
		}
	}
}

func countMovies(s *store.Store) int {
	movies, _, _ := s.ListMovies("", 0, 1000)
	return len(movies)
}

func moviePrice(year int, format string, isNew bool) float64 {
	base := 2.99
	if year >= 2020 || isNew {
		base = 5.99
	} else if year >= 2000 {
		base = 3.99
	}
	if format == "Blu-ray" {
		base += 1.00
	}
	return base
}

func sanitizeID(s string) string {
	out := make([]byte, 0, len(s))
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			out = append(out, byte(c))
		}
	}
	return string(out)
}

func seedMerch(s *store.Store) {
	items := []models.MerchItem{
		{ID: "merch-popcorn-bucket", Name: "Popcorn Bucket", Description: "Classic striped popcorn bucket. Unlimited refills for a month.", PointsCost: 50, Stock: 10},
		{ID: "merch-vhs-blank", Name: "Blank VHS Tape", Description: "Record your own movies. Maxell T-120 high grade.", PointsCost: 75, Stock: 5},
		{ID: "merch-poster", Name: "Movie Poster", Description: "Original theatrical poster from the 80s. Random title, mint condition.", PointsCost: 100, Stock: 3},
		{ID: "merch-tshirt", Name: "Store T-Shirt", Description: "The Last Video Store logo tee. Black, cotton, all sizes.", PointsCost: 150, Stock: 8},
		{ID: "merch-free-rental", Name: "Free Rental Coupon", Description: "One free rental on any movie, any format. No late fees.", PointsCost: 200, Stock: 4},
		{ID: "merch-screening", Name: "Private Screening", Description: "After-hours theater access. Bring 5 friends, 2 free rentals each.", PointsCost: 500, Stock: 1},
		{ID: "merch-membership-upgrade", Name: "Tier Upgrade", Description: "Permanent tier upgrade to the next level (max Gold). One-time use.", PointsCost: 1000, Stock: 2},
	}
	for i := range items {
		s.CreateMerchItem(&items[i])
	}
}
