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
		{"Memento", 2000, "Thriller", "DVD", "Nolan", []string{"Guy Pearce", "Carrie-Anne Moss", "Joe Pantoliano"}, "A man with short-term memory loss attempts to track down his wife's murderer.", 3, true},
		{"The Godfather", 1972, "Drama", "VHS", "Coppola", []string{"Marlon Brando", "Al Pacino", "James Caan"}, "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant youngest son.", 3, false},
		{"Goodfellas", 1990, "Drama", "DVD", "Scorsese", []string{"Robert De Niro", "Ray Liotta", "Joe Pesci"}, "The story of Henry Hill and his life in the mafia, covering his relationship with his wife and his mob partners.", 4, false},
		{"Forrest Gump", 1994, "Drama", "VHS", "Zemeckis", []string{"Tom Hanks", "Robin Wright", "Gary Sinise"}, "The presidencies of Kennedy and Johnson through the eyes of an Alabama man with an IQ of 75.", 5, false},
		{"The Shawshank Redemption", 1994, "Drama", "DVD", "Darabont", []string{"Tim Robbins", "Morgan Freeman", "Bob Gunton"}, "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.", 4, false},
		{"Schindlers List", 1993, "Drama", "DVD", "Spielberg", []string{"Liam Neeson", "Ben Kingsley", "Ralph Fiennes"}, "In German-occupied Poland during World War II, industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce.", 3, false},
		{"The Dark Knight", 2008, "Action", "Blu-ray", "Nolan", []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart"}, "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological tests.", 5, false},
		{"Back to the Future", 1985, "SciFi", "VHS", "Zemeckis", []string{"Michael J. Fox", "Christopher Lloyd", "Lea Thompson"}, "Marty McFly, a 17-year-old high school student, is accidentally sent 30 years into the past in a time-traveling DeLorean.", 4, false},
		{"Die Hard", 1988, "Action", "DVD", "McTiernan", []string{"Bruce Willis", "Alan Rickman", "Bonnie Bedelia"}, "NYPD officer John McClane tries to save his wife and several others taken hostage by terrorists.", 4, false},
		{"Terminator 2", 1991, "SciFi", "Blu-ray", "Cameron", []string{"Arnold Schwarzenegger", "Linda Hamilton", "Edward Furlong"}, "A cyborg must protect John Connor from a more advanced and powerful Terminator.", 3, false},
		{"Toy Story", 1995, "Animation", "DVD", "Lasseter", []string{"Tom Hanks", "Tim Allen", "Don Rickles"}, "A cowboy doll is profoundly threatened and jealous when a new spaceman figure supplants him as top toy in a boy's room.", 5, false},
		{"The Silence of the Lambs", 1991, "Horror", "VHS", "Demme", []string{"Jodie Foster", "Anthony Hopkins", "Lawrence A. Bonney"}, "A young F.B.I. cadet must receive the help of an incarcerated and manipulative cannibal killer to help catch another serial killer.", 3, false},
		{"Se7en", 1995, "Thriller", "DVD", "Fincher", []string{"Morgan Freeman", "Brad Pitt", "Kevin Spacey"}, "Two detectives, a rookie and a veteran, hunt a serial killer who uses the seven deadly sins as his motives.", 4, false},
		{"Saving Private Ryan", 1998, "Action", "DVD", "Spielberg", []string{"Tom Hanks", "Matt Damon", "Tom Sizemore"}, "Following the Normandy Landings, a group of U.S. soldiers go behind enemy lines to retrieve a paratrooper whose brothers have been killed.", 4, false},
		{"Gladiator", 2000, "Action", "Blu-ray", "Scott", []string{"Russell Crowe", "Joaquin Phoenix", "Connie Nielsen"}, "A former Roman General sets out to exact vengeance against the corrupt emperor who murdered his family and sent him into slavery.", 4, false},
		{"Gladiator Extended", 2005, "Action", "Blu-ray", "Scott", []string{"Russell Crowe", "Joaquin Phoenix", "Connie Nielsen"}, "Extended cut of the epic historical drama with additional scenes.", 2, true},
		{"Kill Bill Vol.1", 2003, "Action", "DVD", "Tarantino", []string{"Uma Thurman", "Lucy Liu", "Vivica A. Fox"}, "A former assassin wakes from a coma and seeks revenge against the team of assassins who betrayed her.", 4, false},
		{"Eternal Sunshine", 2004, "Romance", "DVD", "Gondry", []string{"Jim Carrey", "Kate Winslet", "Kirsten Dunst"}, "A couple undergo a medical procedure to have each other erased from their memories.", 3, false},
		{"The Departed", 2006, "Thriller", "Blu-ray", "Scorsese", []string{"Leonardo DiCaprio", "Matt Damon", "Jack Nicholson"}, "An undercover cop and a mole in the police attempt to identify each other while infiltrating an Irish gang.", 4, false},
		{"No Country for Old Men", 2007, "Thriller", "DVD", "Coen", []string{"Tommy Lee Jones", "Javier Bardem", "Josh Brolin"}, "Violence and mayhem ensue after a hunter stumbles upon a drug deal gone wrong and more than two million dollars in cash.", 2, false},
		{"There Will Be Blood", 2007, "Drama", "VHS", "Anderson", []string{"Daniel Day-Lewis", "Paul Dano", "Ciaran Hinds"}, "A story of family, religion, hatred, oil and madness, focusing on a turn-of-the-century prospector.", 2, false},
		{"WALL-E", 2008, "Animation", "DVD", "Stanton", []string{"Ben Burtt", "Elissa Knight", "Jeff Garlin"}, "In the distant future, a small waste-collecting robot inadvertently embarks on a space journey that will ultimately decide the fate of mankind.", 5, false},
		{"Inglourious Basterds", 2009, "Action", "DVD", "Tarantino", []string{"Brad Pitt", "Diane Kruger", "Eli Roth"}, "In Nazi-occupied France during World War II, a plan to assassinate Nazi leaders by a group of Jewish U.S. soldiers coincides with a theatre owner's vengeful plans.", 4, false},
		{"Blade Runner", 1982, "SciFi", "VHS", "Scott", []string{"Harrison Ford", "Rutger Hauer", "Sean Young"}, "A blade runner must pursue and terminate four replicants who stole a ship in space and have returned to Earth to find their creator.", 3, false},
		{"The Big Lebowski", 1998, "Comedy", "DVD", "Coen", []string{"Jeff Bridges", "John Goodman", "Julianne Moore"}, "The Dude, mistaken for a millionaire with the same name, seeks restitution for his ruined rug and gets drawn into a kidnapping plot.", 4, false},
		{"The Truman Show", 1998, "Comedy", "VHS", "Weir", []string{"Jim Carrey", "Laura Linney", "Noah Emmerich"}, "An insurance salesman discovers his whole life is actually a reality TV show.", 3, false},
		{"American Beauty", 1999, "Drama", "DVD", "Mendes", []string{"Kevin Spacey", "Annette Bening", "Thora Birch"}, "A sexually frustrated suburban father has a mid-life crisis after becoming infatuated with his daughter's best friend.", 3, false},
		{"Requiem for a Dream", 2000, "Drama", "DVD", "Aronofsky", []string{"Ellen Burstyn", "Jared Leto", "Jennifer Connelly"}, "The drug-induced utopias of four Coney Island people are shattered.", 2, false},
		{"Spirited Away", 2001, "Animation", "Blu-ray", "Miyazaki", []string{"Rumi Hiiragi", "Miyu Irino", "Mari Natsuki"}, "A sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits.", 4, false},
		{"City of God", 2002, "Drama", "DVD", "Meirelles", []string{"Alexandre Rodrigues", "Leandro Firmino", "Phellipe Haagensen"}, "In the slums of Rio, two kids' paths diverge as one struggles to become a photographer and the other a kingpin.", 3, false},
		{"The Usual Suspects", 1995, "Thriller", "VHS", "Singer", []string{"Kevin Spacey", "Gabriel Byrne", "Chazz Palminteri"}, "A sole survivor tells of the twisty events leading up to a horrific gun battle on a boat.", 3, false},
		{"American History X", 1998, "Drama", "DVD", "Kaye", []string{"Edward Norton", "Edward Furlong", "Beverly D'Angelo"}, "A former neo-nazi skinhead tries to prevent his younger brother from going down the same wrong path.", 3, false},
		{"The Green Mile", 1999, "Drama", "VHS", "Darabont", []string{"Tom Hanks", "Michael Clarke Duncan", "David Morse"}, "The lives of guards on Death Row are affected by one of their charges: a black man accused of child murder and rape.", 4, false},
		{"District 9", 2009, "SciFi", "DVD", "Blomkamp", []string{"Sharlto Copley", "Jason Cope", "Nathalie Boltt"}, "Violence ensues after an extraterrestrial race forced to live in slum-like conditions on Earth.", 3, false},
		{"Leon The Professional", 1994, "Action", "VHS", "Besson", []string{"Jean Reno", "Gary Oldman", "Natalie Portman"}, "A professional assassin reluctantly takes care of a 12-year-old neighbor girl after her parents are murdered.", 3, false},
		{"The Prestige", 2006, "Thriller", "DVD", "Nolan", []string{"Hugh Jackman", "Christian Bale", "Michael Caine"}, "Two stage magicians engage in competitive one-upmanship in an attempt to create the ultimate stage illusion.", 3, false},
		{"Casablanca", 1942, "Romance", "VHS", "Curtiz", []string{"Humphrey Bogart", "Ingrid Bergman", "Paul Henreid"}, "A cynical expatriate American cafe owner struggles to decide whether or not to help his former lover and her fugitive husband.", 3, false},
		{"Psycho", 1960, "Horror", "VHS", "Hitchcock", []string{"Anthony Perkins", "Janet Leigh", "Vera Miles"}, "A Phoenix secretary embezzles $40,000 from her employer's client, goes on the run and checks into a remote motel.", 4, false},
		{"Rear Window", 1954, "Thriller", "DVD", "Hitchcock", []string{"James Stewart", "Grace Kelly", "Wendell Corey"}, "A wheelchair-bound photographer spies on his neighbors and becomes convinced one of them has committed murder.", 3, false},
		{"Vertigo", 1958, "Thriller", "DVD", "Hitchcock", []string{"James Stewart", "Kim Novak", "Barbara Bel Geddes"}, "A former San Francisco police detective juggles wrestling with his personal demons and becoming obsessed with a hauntingly beautiful woman.", 3, false},
		{"Citizen Kane", 1941, "Drama", "VHS", "Welles", []string{"Orson Welles", "Joseph Cotten", "Dorothy Comingore"}, "Following the death of a publishing tycoon, news reporters scramble to discover the meaning of his final utterance.", 3, false},
		{"Lawrence of Arabia", 1962, "Drama", "Blu-ray", "Lean", []string{"Peter O'Toole", "Alec Guinness", "Anthony Quinn"}, "The story of T.E. Lawrence, the English officer who successfully united and led the diverse, often warring, Arab tribes during World War I.", 2, false},
		{"2001: A Space Odyssey", 1968, "SciFi", "Blu-ray", "Kubrick", []string{"Keir Dullea", "Gary Lockwood", "William Sylvester"}, "Humanity finds a mysterious object buried beneath the lunar surface and sets off on a quest to find its origins.", 3, false},
		{"A Clockwork Orange", 1971, "Drama", "VHS", "Kubrick", []string{"Malcolm McDowell", "Patrick Magee", "Michael Bates"}, "In the future, a sadistic gang leader is imprisoned and volunteers for a conduct-aversion experiment.", 3, false},
		{"The Shining", 1980, "Horror", "DVD", "Kubrick", []string{"Jack Nicholson", "Shelley Duvall", "Danny Lloyd"}, "A family heads to an isolated hotel for the winter where a sinister presence influences the father into violence.", 4, false},
		{"Full Metal Jacket", 1987, "Action", "DVD", "Kubrick", []string{"Matthew Modine", "R. Lee Ermey", "Vincent D'Onofrio"}, "A pragmatic U.S. Marine observes the dehumanizing effects the Vietnam War has on his fellow recruits.", 3, false},
		{"Apocalypse Now", 1979, "Action", "Blu-ray", "Coppola", []string{"Martin Sheen", "Marlon Brando", "Robert Duvall"}, "A U.S. Army officer serving in Vietnam is tasked with assassinating a renegade Special Forces Colonel.", 3, false},
		{"Taxi Driver", 1976, "Drama", "VHS", "Scorsese", []string{"Robert De Niro", "Jodie Foster", "Cybill Shepherd"}, "A mentally unstable veteran works as a nighttime taxi driver in New York City.", 4, false},
		{"Raging Bull", 1980, "Drama", "DVD", "Scorsese", []string{"Robert De Niro", "Joe Pesci", "Cathy Moriarty"}, "The life of boxer Jake LaMotta, whose violence and temper that led him to the top in the ring destroyed his life outside of it.", 3, false},
		{"The Wolf of Wall Street", 2013, "Comedy", "Blu-ray", "Scorsese", []string{"Leonardo DiCaprio", "Jonah Hill", "Margot Robbie"}, "Based on the true story of Jordan Belfort, from his rise to a wealthy stock-broker living the high life to his fall involving crime, corruption and the federal government.", 4, false},
		{"Alien", 1979, "Horror", "DVD", "Scott", []string{"Sigourney Weaver", "Tom Skerritt", "John Hurt"}, "The crew of a commercial spacecraft encounter a deadly lifeform after investigating an unknown transmission.", 4, false},
		{"Aliens", 1986, "Action", "Blu-ray", "Cameron", []string{"Sigourney Weaver", "Michael Biehn", "Carrie Henn"}, "Fifty-seven years after surviving an apocalyptic attack aboard her space vessel, Ellen Ripley is called back to the planet LV-426.", 4, false},
		{"Raiders of the Lost Ark", 1981, "Action", "DVD", "Spielberg", []string{"Harrison Ford", "Karen Allen", "Paul Freeman"}, "In 1936, archaeologist Indiana Jones is hired by the U.S. government to locate the Ark of the Covenant.", 4, false},
		{"Jaws", 1975, "Thriller", "VHS", "Spielberg", []string{"Roy Scheider", "Robert Shaw", "Richard Dreyfuss"}, "When a killer shark unleashes chaos on a beach community off Cape Cod, it's up to a local sheriff, a marine biologist, and an old seafarer to hunt the beast down.", 5, false},
		{"E.T.", 1982, "SciFi", "VHS", "Spielberg", []string{"Henry Thomas", "Drew Barrymore", "Peter Coyote"}, "A troubled child summons the courage to help a friendly alien escape Earth and return to his home world.", 4, false},
		{"The Thing", 1982, "Horror", "DVD", "Carpenter", []string{"Kurt Russell", "Wilford Brimley", "Keith David"}, "A research team in Antarctica is hunted by a shape-shifting alien that assumes the appearance of its victims.", 3, false},
		{"Halloween", 1978, "Horror", "VHS", "Carpenter", []string{"Donald Pleasence", "Jamie Lee Curtis", "Nancy Kyes"}, "Fifteen years after murdering his sister on Halloween night 1963, Michael Myers escapes from a mental hospital and returns to the small town of Haddonfield.", 4, false},
		{"The Exorcist", 1973, "Horror", "DVD", "Friedkin", []string{"Ellen Burstyn", "Max von Sydow", "Linda Blair"}, "When a 12-year-old girl is possessed by a mysterious entity, her mother seeks the help of two priests.", 3, false},
		{"One Flew Over the Cuckoo's Nest", 1975, "Drama", "VHS", "Forman", []string{"Jack Nicholson", "Louise Fletcher", "Michael Berryman"}, "A criminal pleads insanity and is admitted to a mental institution, where he rebels against the oppressive nurse.", 3, false},
		{"Chinatown", 1974, "Thriller", "DVD", "Polanski", []string{"Jack Nicholson", "Faye Dunaway", "John Huston"}, "A private detective hired to expose an adulterer in 1930s Los Angeles finds himself caught up in a web of deceit, corruption, and murder.", 2, false},
		{"The Deer Hunter", 1978, "Drama", "DVD", "Cimino", []string{"Robert De Niro", "Christopher Walken", "John Savage"}, "An in-depth examination of the ways in which the Vietnam War impacts and disrupts the lives of several friends in a small steel mill town.", 2, false},
		{"Unforgiven", 1992, "Action", "DVD", "Eastwood", []string{"Clint Eastwood", "Gene Hackman", "Morgan Freeman"}, "Retired Old West gunslinger William Munny reluctantly takes on one last job, with the help of his old partner Ned Logan.", 3, false},
		{"Reservoir Dogs", 1992, "Thriller", "VHS", "Tarantino", []string{"Harvey Keitel", "Tim Roth", "Michael Madsen"}, "When a simple jewelry heist goes horribly wrong, the surviving criminals suspect that one of them is a police informant.", 3, false},
		{"Braveheart", 1995, "Action", "Blu-ray", "Gibson", []string{"Mel Gibson", "Sophie Marceau", "Patrick McGoohan"}, "Scottish warrior William Wallace leads his countrymen in a rebellion to free his homeland from the tyranny of King Edward I.", 4, false},
		{"Heat", 1995, "Action", "DVD", "Mann", []string{"Al Pacino", "Robert De Niro", "Val Kilmer"}, "A group of high-end professional thieves start to feel the heat from the LAPD when they unknowingly leave a clue at their latest heist.", 3, false},
		{"Casino", 1995, "Drama", "DVD", "Scorsese", []string{"Robert De Niro", "Sharon Stone", "Joe Pesci"}, "A tale of greed, deception, money, power, and murder occur between two mobster best friends and a trophy wife over a gambling empire.", 3, false},
		{"Donnie Darko", 2001, "SciFi", "DVD", "Kelly", []string{"Jake Gyllenhaal", "Jena Malone", "Mary McDonnell"}, "After narrowly escaping a bizarre accident, a troubled teenager is plagued by visions of a man in a large rabbit suit.", 2, false},
		{"The Grand Budapest Hotel", 2014, "Comedy", "Blu-ray", "Anderson", []string{"Ralph Fiennes", "F. Murray Abraham", "Mathieu Amalric"}, "A writer encounters the owner of an aging high-class hotel, who tells him of his early years serving as a lobby boy.", 3, false},
		{"Mad Max: Fury Road", 2015, "Action", "Blu-ray", "Miller", []string{"Tom Hardy", "Charlize Theron", "Nicholas Hoult"}, "In a post-apocalyptic wasteland, a woman rebels against a tyrannical ruler in search for her homeland.", 5, true},
		{"Parasite", 2019, "Thriller", "Blu-ray", "Bong", []string{"Kang-ho Song", "Sun-kyun Lee", "Yeo-jeong Jo"}, "Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.", 4, false},
		{"Interstellar", 2014, "SciFi", "Blu-ray", "Nolan", []string{"Matthew McConaughey", "Anne Hathaway", "Jessica Chastain"}, "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.", 5, false},
		{"The Avengers", 2012, "Action", "Blu-ray", "Whedon", []string{"Robert Downey Jr.", "Chris Evans", "Scarlett Johansson"}, "Earth's mightiest heroes must come together to stop Loki and his alien army from enslaving humanity.", 5, false},
		{"Get Out", 2017, "Horror", "DVD", "Peele", []string{"Daniel Kaluuya", "Allison Williams", "Bradley Whitford"}, "A young African-American visits his white girlfriend's parents for the weekend, where his simmering uneasiness about their reception of him eventually reaches a boiling point.", 3, false},
		{"Moonlight", 2016, "Drama", "DVD", "Jenkins", []string{"Mahershala Ali", "Naomie Harris", "Trevante Rhodes"}, "A young African-American man grapples with his identity and sexuality while experiencing the everyday struggles of childhood, adolescence, and burgeoning adulthood.", 2, false},
		{"Trainspotting", 1996, "Drama", "DVD", "Boyle", []string{"Ewan McGregor", "Ewen Bremner", "Jonny Lee Miller"}, "Renton, deeply immersed in the Edinburgh drug scene, tries to clean up and get out, despite the allure of the drugs.", 3, false},
		{"Snatch", 2000, "Comedy", "DVD", "Ritchie", []string{"Jason Statham", "Brad Pitt", "Stephen Graham"}, "Unscrupulous boxing promoters, violent bookmakers, a Russian gangster, and incompetent amateur robbers compete to track down a priceless stolen diamond.", 3, false},
		{"Zodiac", 2007, "Thriller", "DVD", "Fincher", []string{"Jake Gyllenhaal", "Robert Downey Jr.", "Mark Ruffalo"}, "A San Francisco cartoonist becomes an amateur detective obsessed with tracking down the Zodiac Killer.", 3, false},
		{"Gone Girl", 2014, "Thriller", "Blu-ray", "Fincher", []string{"Ben Affleck", "Rosamund Pike", "Neil Patrick Harris"}, "With his wife's disappearance having become the focus of an intense media circus, a man sees the spotlight turned on him.", 3, false},
		{"The Social Network", 2010, "Drama", "DVD", "Fincher", []string{"Jesse Eisenberg", "Andrew Garfield", "Justin Timberlake"}, "As Harvard student Mark Zuckerberg creates the social networking site that would become Facebook, he is sued by two brothers who claimed he stole their idea.", 4, false},
		{"Dunkirk", 2017, "Action", "Blu-ray", "Nolan", []string{"Fionn Whitehead", "Barry Keoghan", "Mark Rylance"}, "Allied soldiers from Belgium, the British Commonwealth and Empire, and France are surrounded by the German Army and evacuated during a fierce battle.", 4, false},
		{"1917", 2019, "Action", "Blu-ray", "Mendes", []string{"Dean-Charles Chapman", "George MacKay", "Daniel Mays"}, "Two young British soldiers during World War I are given an impossible mission: deliver a message that will stop 1,600 men from walking into a trap.", 3, false},
		{"Joker", 2019, "Drama", "Blu-ray", "Phillips", []string{"Joaquin Phoenix", "Robert De Niro", "Zazie Beetz"}, "In Gotham City, mentally troubled comedian Arthur Fleck is disregarded and mistreated by society.", 4, false},
		{"Dune", 2021, "SciFi", "Blu-ray", "Villeneuve", []string{"Timothee Chalamet", "Rebecca Ferguson", "Zendaya"}, "A noble family becomes embroiled in a war for control over the galaxy's most valuable asset.", 5, true},
		{"Everything Everywhere All at Once", 2022, "SciFi", "Blu-ray", "Kwan", []string{"Michelle Yeoh", "Stephanie Hsu", "Ke Huy Quan"}, "An aging Chinese immigrant is swept up in an insane adventure, where she alone can save the world by exploring other universes.", 4, true},
		{"Whiplash", 2014, "Drama", "DVD", "Chazelle", []string{"Miles Teller", "J.K. Simmons", "Melissa Benoist"}, "A promising young drummer enrolls at a cut-throat music conservatory where his dreams of greatness are mentored by an instructor who will stop at nothing.", 3, false},
		{"La La Land", 2016, "Romance", "Blu-ray", "Chazelle", []string{"Ryan Gosling", "Emma Stone", "John Legend"}, "While navigating their careers in Los Angeles, a pianist and an actress fall in love while attempting to reconcile their aspirations for the future.", 3, false},
		{"The Revenant", 2015, "Action", "Blu-ray", "Iñárritu", []string{"Leonardo DiCaprio", "Tom Hardy", "Will Poulter"}, "A frontiersman on a fur trading expedition in the 1820s fights for survival after being mauled by a bear and left for dead.", 3, false},
		{"Nightcrawler", 2014, "Thriller", "DVD", "Gilroy", []string{"Jake Gyllenhaal", "Rene Russo", "Bill Paxton"}, "When Louis Bloom, a driven man desperate for work, muscles into the world of L.A. crime journalism, he blurs the line between observer and participant.", 2, false},
		{"Ex Machina", 2015, "SciFi", "DVD", "Garland", []string{"Domhnall Gleeson", "Alicia Vikander", "Oscar Isaac"}, "A young programmer is selected to participate in a ground-breaking experiment in synthetic intelligence.", 3, false},
		{"Arrival", 2016, "SciFi", "Blu-ray", "Villeneuve", []string{"Amy Adams", "Jeremy Renner", "Forest Whitaker"}, "A linguist works with the military to communicate with alien lifeforms after twelve mysterious spacecraft appear around the world.", 3, false},
		{"Logan", 2017, "Action", "Blu-ray", "Mangold", []string{"Hugh Jackman", "Patrick Stewart", "Dafne Keen"}, "In a future where mutants are nearly extinct, an elderly and weary Logan leads a quiet life. But when a young mutant arrives, he must protect her.", 4, false},
		{"Spider-Man: Into the Spider-Verse", 2018, "Animation", "Blu-ray", "Persichetti", []string{"Shameik Moore", "Jake Johnson", "Hailee Steinfeld"}, "Teen Miles Morales becomes the Spider-Man of his universe and must join with five counterparts from other dimensions to stop a threat.", 4, false},
		{"Coco", 2017, "Animation", "DVD", "Unkrich", []string{"Anthony Gonzalez", "Gael Garcia Bernal", "Benjamin Bratt"}, "Aspiring musician Miguel enters the Land of the Dead to find his great-great-grandfather, a legendary singer.", 4, false},
		{"Up", 2009, "Animation", "DVD", "Docter", []string{"Edward Asner", "Jordan Nagai", "Christopher Plummer"}, "Seventy-eight year old Carl Fredricksen travels to Paradise Falls, unintentionally taking a young stowaway.", 5, false},
		{"Finding Nemo", 2003, "Animation", "DVD", "Stanton", []string{"Albert Brooks", "Ellen DeGeneres", "Alexander Gould"}, "After his son is captured, a timid clownfish sets out on a journey across the ocean to bring him home.", 5, false},
		{"The Incredibles", 2004, "Animation", "DVD", "Bird", []string{"Craig T. Nelson", "Holly Hunter", "Samuel L. Jackson"}, "While trying to lead a quiet suburban life, a family of undercover superheroes are forced into action to save the world.", 4, false},
		{"Shrek", 2001, "Animation", "DVD", "Adamson", []string{"Mike Myers", "Eddie Murphy", "Cameron Diaz"}, "A mean lord exiles fairytale creatures to the swamp of a grumpy ogre, who must go on a quest and rescue a princess.", 5, false},
		{"The Lion King", 1994, "Animation", "VHS", "Allers", []string{"Matthew Broderick", "Jeremy Irons", "James Earl Jones"}, "Lion prince Simba flees his kingdom after the murder of his father, only to learn the true meaning of responsibility and bravery.", 5, false},
		{"Predator", 1987, "Action", "DVD", "McTiernan", []string{"Arnold Schwarzenegger", "Carl Weathers", "Kevin Peter Hall"}, "A team of commandos on a mission in a Central American jungle find themselves hunted by an extraterrestrial warrior.", 3, false},
		{"RoboCop", 1987, "SciFi", "VHS", "Verhoeven", []string{"Peter Weller", "Nancy Allen", "Dan O'Herlihy"}, "In a dystopic and crime-ridden Detroit, a terminally wounded cop returns to the force as a powerful cyborg.", 3, false},
		{"Total Recall", 1990, "SciFi", "DVD", "Verhoeven", []string{"Arnold Schwarzenegger", "Sharon Stone", "Michael Ironside"}, "When a man goes for virtual vacation memories of Mars, an unexpected series of events forces him to go to the planet for real.", 3, false},
		{"Groundhog Day", 1993, "Comedy", "DVD", "Ramis", []string{"Bill Murray", "Andie MacDowell", "Chris Elliott"}, "A narcissistic weatherman finds himself living the same day over and over again.", 4, false},
		{"Ghostbusters", 1984, "Comedy", "VHS", "Reitman", []string{"Bill Murray", "Dan Aykroyd", "Sigourney Weaver"}, "Three former parapsychology professors set up shop as a unique ghost removal service.", 4, false},
		{"Ferris Buellers Day Off", 1986, "Comedy", "VHS", "Hughes", []string{"Matthew Broderick", "Alan Ruck", "Mia Sara"}, "A charismatic high-school student convinces his friends to play hooky and spend one epic day in downtown Chicago.", 3, false},
		{"The Breakfast Club", 1985, "Comedy", "VHS", "Hughes", []string{"Emilio Estevez", "Judd Nelson", "Molly Ringwald"}, "Five high school students meet in Saturday detention and discover how they have a lot more in common than they thought.", 3, false},
		{"Home Alone", 1990, "Comedy", "VHS", "Columbus", []string{"Macaulay Culkin", "Joe Pesci", "Daniel Stern"}, "An eight-year-old troublemaker must protect his house from a pair of burglars when he is accidentally left home alone.", 4, false},
		{"Gremlins", 1984, "Comedy", "VHS", "Dante", []string{"Zach Galligan", "Phoebe Cates", "Hoyt Axton"}, "A young man inadvertently breaks three important rules concerning his new pet and unleashes a horde of malevolently mischievous monsters.", 3, false},
		{"Rocky", 1976, "Drama", "VHS", "Avildsen", []string{"Sylvester Stallone", "Talia Shire", "Burt Young"}, "A small-time Philadelphia boxer gets a supremely rare chance to fight the world heavyweight champion.", 4, false},
		{"The Fifth Element", 1997, "SciFi", "DVD", "Besson", []string{"Bruce Willis", "Milla Jovovich", "Gary Oldman"}, "In the colorful future, a cab driver unwittingly becomes the central figure in the search for a legendary cosmic weapon.", 3, false},
		{"Children of Men", 2006, "SciFi", "DVD", "Cuaron", []string{"Clive Owen", "Julianne Moore", "Michael Caine"}, "In 2027, in a chaotic world in which women have become somehow infertile, a former activist agrees to help transport a miraculously pregnant woman.", 2, false},
		{"Oldboy", 2003, "Thriller", "DVD", "Park", []string{"Min-sik Choi", "Ji-tae Yoo", "Hye-jeong Kang"}, "After being kidnapped and imprisoned for fifteen years, Oh Dae-Su is released, only to find that he must find his captor in five days.", 2, false},
		{"12 Angry Men", 1957, "Drama", "VHS", "Lumet", []string{"Henry Fonda", "Lee J. Cobb", "Martin Balsam"}, "The jury in a New York City murder trial is frustrated by a single member whose skeptical caution forces them to more carefully consider the evidence.", 2, false},
		{"Eraserhead", 1977, "Horror", "VHS", "Lynch", []string{"Jack Nance", "Charlotte Stewart", "Allen Joseph"}, "Henry Spencer tries to survive his industrial environment, his angry girlfriend, and the unbearable screams of his newly born mutant child.", 1, false},
		{"Mulholland Drive", 2001, "Thriller", "DVD", "Lynch", []string{"Naomi Watts", "Laura Harring", "Justin Theroux"}, "After a car wreck on the winding Mulholland Drive renders a woman amnesiac, she and a perky Hollywood-hopeful search for clues and answers.", 2, false},
		{"Blue Velvet", 1986, "Thriller", "VHS", "Lynch", []string{"Kyle MacLachlan", "Isabella Rossellini", "Dennis Hopper"}, "The discovery of a severed human ear found in a field leads a young man on an investigation related to a beautiful nightclub singer.", 2, false},
		{"The Grand Illusion", 1937, "Drama", "VHS", "Renoir", []string{"Jean Gabin", "Dita Parlo", "Pierre Fresnay"}, "During WWI, two French soldiers are captured and imprisoned in a German POW camp.", 1, false},
		{"Rashomon", 1950, "Drama", "DVD", "Kurosawa", []string{"Toshiro Mifune", "Machiko Kyo", "Masayuki Mori"}, "The rape of a bride and the murder of her samurai husband are recalled from the perspectives of a bandit, the bride, the samurai's ghost and a woodcutter.", 2, false},
		{"Seven Samurai", 1954, "Action", "DVD", "Kurosawa", []string{"Toshiro Mifune", "Takashi Shimura", "Keiko Tsushima"}, "A poor village under attack by bandits recruits seven unemployed samurai to help them defend themselves.", 2, false},
		{"Yojimbo", 1961, "Action", "DVD", "Kurosawa", []string{"Toshiro Mifune", "Tatsuya Nakadai", "Yoko Tsukasa"}, "A crafty ronin comes to a town divided by two criminal gangs and decides to play them against each other.", 2, false},
		{"Amelie", 2001, "Romance", "DVD", "Jeunet", []string{"Audrey Tautou", "Mathieu Kassovitz", "Rufus"}, "Amelie, an innocent and naive girl in Paris, with her own sense of justice, decides to help those around her.", 3, false},
		{"Pan's Labyrinth", 2006, "Drama", "DVD", "Del Toro", []string{"Ivana Baquero", "Sergi Lopez", "Maribel Verdu"}, "In the Falangist Spain of 1944, the bookish young stepdaughter of a sadistic army officer escapes into an eerie but captivating fantasy world.", 2, false},
		{"The Shape of Water", 2017, "Romance", "Blu-ray", "Del Toro", []string{"Sally Hawkins", "Michael Shannon", "Richard Jenkins"}, "A lonely janitor forms a unique relationship with an amphibious creature that is being held in captivity.", 3, false},
		{"Drive", 2011, "Action", "Blu-ray", "Refn", []string{"Ryan Gosling", "Carey Mulligan", "Bryan Cranston"}, "A mysterious Hollywood action film stuntman gets in trouble with gangsters when he tries to help his neighbor's husband rob a pawn shop.", 2, false},
		{"Baby Driver", 2017, "Action", "Blu-ray", "Wright", []string{"Ansel Elgort", "Kevin Spacey", "Lily James"}, "After being coerced into working for a crime boss, a young getaway driver finds himself taking part in a heist doomed to fail.", 3, false},
		{"Hot Fuzz", 2007, "Comedy", "DVD", "Wright", []string{"Simon Pegg", "Nick Frost", "Jim Broadbent"}, "A skilled London police officer is transferred to a small town with a dark secret.", 3, false},
		{"Shaun of the Dead", 2004, "Comedy", "DVD", "Wright", []string{"Simon Pegg", "Nick Frost", "Kate Ashfield"}, "The uneventful life of a London electronics salesman is disrupted by a zombie apocalypse.", 3, false},
		{"The Worlds End", 2013, "Comedy", "Blu-ray", "Wright", []string{"Simon Pegg", "Nick Frost", "Paddy Considine"}, "Five friends who reunite in an attempt to top their epic pub crawl from 20 years earlier unwittingly become humankind's only hope for survival.", 2, false},
		{"Good Will Hunting", 1997, "Drama", "VHS", "Van Sant", []string{"Matt Damon", "Robin Williams", "Ben Affleck"}, "Will Hunting, a janitor at M.I.T., has a gift for mathematics, but needs help from a psychologist to find direction in his life.", 4, false},
		{"Dead Poets Society", 1989, "Drama", "VHS", "Weir", []string{"Robin Williams", "Ethan Hawke", "Robert Sean Leonard"}, "Maverick teacher John Keating uses poetry to embolden his boarding school students to new heights of self-expression.", 3, false},
		{"A Beautiful Mind", 2001, "Drama", "DVD", "Howard", []string{"Russell Crowe", "Ed Harris", "Jennifer Connelly"}, "After John Nash, a brilliant but asocial mathematical genius, accepts secret work in cryptography, his life takes a turn for the nightmarish.", 3, false},
		{"The Martian", 2015, "SciFi", "Blu-ray", "Scott", []string{"Matt Damon", "Jessica Chastain", "Kristen Wiig"}, "An astronaut becomes stranded on Mars after his team assume him dead, and must rely on his ingenuity to find a way to signal to Earth that he is alive.", 4, false},
		{"Gravity", 2013, "SciFi", "Blu-ray", "Cuaron", []string{"Sandra Bullock", "George Clooney", "Ed Harris"}, "Two astronauts work together to survive after an accident leaves them stranded in space.", 3, false},
		{"Sunshine", 2007, "SciFi", "DVD", "Boyle", []string{"Cillian Murphy", "Rose Byrne", "Chris Evans"}, "A team of international astronauts are sent on a dangerous mission to reignite the dying Sun.", 2, false},
		{"Edge of Tomorrow", 2014, "SciFi", "Blu-ray", "Liman", []string{"Tom Cruise", "Emily Blunt", "Brendan Gleeson"}, "A soldier fighting aliens gets to relive the same day over and over again, the day restarting every time he dies.", 3, false},
		{"Looper", 2012, "SciFi", "Blu-ray", "Johnson", []string{"Joseph Gordon-Levitt", "Bruce Willis", "Emily Blunt"}, "In 2074, the mob sends victims back in time to get killed by Loopers. Joe is getting rich until his future self is sent back.", 3, false},
		{"Minority Report", 2002, "SciFi", "DVD", "Spielberg", []string{"Tom Cruise", "Colin Farrell", "Samantha Morton"}, "In a future where a special police unit is able to arrest murderers before they commit their crimes, an officer from that unit is himself accused of a future murder.", 3, false},
		{"Annihilation", 2018, "SciFi", "Blu-ray", "Garland", []string{"Natalie Portman", "Jennifer Jason Leigh", "Tessa Thompson"}, "A biologist signs up for a dangerous, secret expedition into a mysterious zone where the laws of nature don't apply.", 2, false},
		{"Room", 2015, "Drama", "DVD", "Abrahamson", []string{"Brie Larson", "Jacob Tremblay", "Sean Bridgers"}, "Held captive for years in an enclosed space, a woman and her young son finally gain their freedom.", 2, false},
		{"The Florida Project", 2017, "Drama", "DVD", "Baker", []string{"Brooklynn Prince", "Willem Dafoe", "Bria Vinaite"}, "Set over one summer, the film follows precocious six-year-old Moonee as she courts mischief and adventure with her ragtag playmates.", 2, false},
		{"Boyhood", 2014, "Drama", "DVD", "Linklater", []string{"Ellar Coltrane", "Patricia Arquette", "Ethan Hawke"}, "The life of Mason, from early childhood to his arrival at college.", 2, false},
		{"Birdman", 2014, "Comedy", "Blu-ray", "Iñárritu", []string{"Michael Keaton", "Zach Galifianakis", "Edward Norton"}, "A washed-up superhero actor attempts to revive his fading career by writing, directing, and starring in a Broadway production.", 3, false},
		{"The Artist", 2011, "Drama", "DVD", "Hazanavicius", []string{"Jean Dujardin", "Berenice Bejo", "John Goodman"}, "An egomaniacal film star develops a relationship with a young dancer against the backdrop of Hollywood's silent to sound transition.", 2, false},
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
