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

var genres = []string{"Action", "Comedy", "Horror", "SciFi", "Drama", "Thriller", "Romance", "Animation"}
var directors = []string{"Wachowski", "Tarantino", "Fincher", "Spielberg", "Nolan", "Scorsese", "Coppola", "Kubrick", "Scott", "Cameron"}

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
	fmt.Printf("Seeded %d movies and 8 users.\n", countMovies(s))
}

func seedUsers(s *store.Store) {
	entries := []struct {
		name, pass string
		tier       bitmask.Permission
		banned     bool
	}{
		{"bronze", "password1", bitmask.TierBronze, false},
		{"silver", "password2", bitmask.TierSilver, false},
		{"gold", "password3", bitmask.TierGold, false},
		{"employee", "password4", bitmask.TierEmployee, false},
		{"supervisor", "password8", bitmask.TierSupervisor, false},
		{"manager", "password5", bitmask.TierManager, false},
		{"owner", "password6", bitmask.TierOwner, false},
		{"banned", "password7", bitmask.TierBronze, true},
	}

	for _, e := range entries {
		hash, _ := auth.HashPassword(e.pass)
		now := time.Now().Unix()
		user := &models.User{
			ID:           fmt.Sprintf("seed-%s", e.name),
			Username:     e.name,
			PasswordHash: hash,
			Tier:         e.tier,
			MaxRentals:   bitmask.MaxRentalsForTier(e.tier),
			Banned:       e.banned,
			CreatedAt:    now,
			UpdatedAt:    now,
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

func sanitizeID(s string) string {
	out := make([]byte, 0, len(s))
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			out = append(out, byte(c))
		}
	}
	return string(out)
}
