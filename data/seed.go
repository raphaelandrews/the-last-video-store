package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
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
	fmt.Println("Seeded 40 movies and 8 users.")
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
}

type movieSeed struct {
	Title        string
	Year         int
	Genre        string
	Format       string
	Director     string
	Cast         []string
	Synopsis     string
	CopiesTotal  int
	IsNewRelease bool
}

func seedMovies(s *store.Store) {
	movies := []movieSeed{
		{"The Matrix", 1999, "SciFi", "DVD", "Wachowski", []string{"Keanu Reeves", "Laurence Fishburne", "Carrie-Anne Moss"}, "A computer hacker learns the true nature of reality.", 5, false},
		{"Pulp Fiction", 1994, "Drama", "VHS", "Tarantino", []string{"John Travolta", "Samuel L. Jackson", "Uma Thurman"}, "The lives of mob hitmen, a boxer, and a gangster intertwine.", 3, false},
		{"Fight Club", 1999, "Drama", "DVD", "Fincher", []string{"Brad Pitt", "Edward Norton", "Helena Bonham Carter"}, "An insomniac office worker forms an underground fight club.", 4, false},
		{"Jurassic Park", 1993, "Action", "Blu-ray", "Spielberg", []string{"Sam Neill", "Laura Dern", "Jeff Goldblum"}, "A theme park with genetically engineered dinosaurs turns deadly.", 5, false},
		{"Inception", 2010, "SciFi", "DVD", "Nolan", []string{"Leonardo DiCaprio", "Joseph Gordon-Levitt", "Elliot Page"}, "A thief steals corporate secrets through dream-sharing technology.", 4, true},
		{"The Godfather", 1972, "Drama", "VHS", "Coppola", []string{"Marlon Brando", "Al Pacino", "James Caan"}, "The aging patriarch of an organized crime dynasty transfers control.", 3, false},
		{"Goodfellas", 1990, "Drama", "DVD", "Scorsese", []string{"Robert De Niro", "Ray Liotta", "Joe Pesci"}, "The story of Henry Hill and his life in the mafia.", 4, false},
		{"Forrest Gump", 1994, "Drama", "VHS", "Zemeckis", []string{"Tom Hanks", "Robin Wright", "Gary Sinise"}, "The presidencies of Kennedy and Johnson through the eyes of an Alabama man.", 5, false},
		{"The Shawshank Redemption", 1994, "Drama", "DVD", "Darabont", []string{"Tim Robbins", "Morgan Freeman", "Bob Gunton"}, "Two imprisoned men bond over a number of years.", 4, false},
		{"Schindlers List", 1993, "Drama", "DVD", "Spielberg", []string{"Liam Neeson", "Ben Kingsley", "Ralph Fiennes"}, "Industrialist Oskar Schindler gradually becomes concerned for his Jewish workforce.", 3, false},
		{"The Dark Knight", 2008, "Action", "Blu-ray", "Nolan", []string{"Christian Bale", "Heath Ledger", "Aaron Eckhart"}, "Batman must accept one of the greatest psychological tests against the Joker.", 5, false},
		{"Back to the Future", 1985, "SciFi", "VHS", "Zemeckis", []string{"Michael J. Fox", "Christopher Lloyd", "Lea Thompson"}, "Marty McFly is sent back in time to 1955.", 4, false},
		{"Die Hard", 1988, "Action", "DVD", "McTiernan", []string{"Bruce Willis", "Alan Rickman", "Bonnie Bedelia"}, "NYPD officer John McClane tries to save hostages on Christmas Eve.", 4, false},
		{"Terminator 2", 1991, "SciFi", "Blu-ray", "Cameron", []string{"Arnold Schwarzenegger", "Linda Hamilton", "Edward Furlong"}, "A cyborg must protect John Connor from a more advanced Terminator.", 3, false},
		{"Toy Story", 1995, "Animation", "DVD", "Lasseter", []string{"Tom Hanks", "Tim Allen", "Don Rickles"}, "A cowboy doll is threatened when a new spaceman figure supplants him.", 5, false},
		{"The Silence of the Lambs", 1991, "Horror", "VHS", "Demme", []string{"Jodie Foster", "Anthony Hopkins", "Lawrence A. Bonney"}, "A young FBI cadet receives help from an incarcerated cannibal killer.", 3, false},
		{"Se7en", 1995, "Thriller", "DVD", "Fincher", []string{"Morgan Freeman", "Brad Pitt", "Kevin Spacey"}, "Two detectives hunt a serial killer using the seven deadly sins.", 4, false},
		{"Saving Private Ryan", 1998, "Action", "DVD", "Spielberg", []string{"Tom Hanks", "Matt Damon", "Tom Sizemore"}, "U.S. soldiers go behind enemy lines to retrieve a paratrooper.", 4, false},
		{"Gladiator", 2000, "Action", "Blu-ray", "Scott", []string{"Russell Crowe", "Joaquin Phoenix", "Connie Nielsen"}, "A former Roman General sets out to exact vengeance against the corrupt emperor.", 4, false},
		{"Memento", 2000, "Thriller", "DVD", "Nolan", []string{"Guy Pearce", "Carrie-Anne Moss", "Joe Pantoliano"}, "A man with short-term memory loss tracks down his wife's murderer.", 3, true},
		{"Kill Bill Vol.1", 2003, "Action", "DVD", "Tarantino", []string{"Uma Thurman", "Lucy Liu", "Vivica A. Fox"}, "A former assassin wakes from a coma and seeks revenge.", 4, false},
		{"Eternal Sunshine", 2004, "Romance", "DVD", "Gondry", []string{"Jim Carrey", "Kate Winslet", "Kirsten Dunst"}, "A couple undergo a procedure to erase each other from their memories.", 3, false},
		{"The Departed", 2006, "Thriller", "Blu-ray", "Scorsese", []string{"Leonardo DiCaprio", "Matt Damon", "Jack Nicholson"}, "An undercover cop and a mole in the police attempt to identify each other.", 4, false},
		{"No Country for Old Men", 2007, "Thriller", "DVD", "Coen", []string{"Tommy Lee Jones", "Javier Bardem", "Josh Brolin"}, "Violence ensues after a hunter stumbles upon a drug deal gone wrong.", 3, false},
		{"There Will Be Blood", 2007, "Drama", "VHS", "Anderson", []string{"Daniel Day-Lewis", "Paul Dano", "Ciaran Hinds"}, "A story of family, religion, hatred, oil and madness.", 2, false},
		{"WALL-E", 2008, "Animation", "DVD", "Stanton", []string{"Ben Burtt", "Elissa Knight", "Jeff Garlin"}, "A small waste-collecting robot inadvertently embarks on a space journey.", 5, false},
		{"Inglourious Basterds", 2009, "Action", "DVD", "Tarantino", []string{"Brad Pitt", "Diane Kruger", "Eli Roth"}, "A plan to assassinate Nazi leaders converges with furious revenge.", 4, false},
		{"Blade Runner", 1982, "SciFi", "VHS", "Scott", []string{"Harrison Ford", "Rutger Hauer", "Sean Young"}, "A blade runner must pursue and terminate four replicants.", 3, false},
		{"The Big Lebowski", 1998, "Comedy", "DVD", "Coen", []string{"Jeff Bridges", "John Goodman", "Julianne Moore"}, "The Dude is mistaken for a millionaire and seeks restitution.", 4, false},
		{"The Truman Show", 1998, "Comedy", "VHS", "Weir", []string{"Jim Carrey", "Laura Linney", "Noah Emmerich"}, "An insurance salesman discovers his life is a reality TV show.", 3, false},
		{"American Beauty", 1999, "Drama", "DVD", "Mendes", []string{"Kevin Spacey", "Annette Bening", "Thora Birch"}, "A depressed suburban father decides to turn his life around.", 3, false},
		{"Requiem for a Dream", 2000, "Drama", "DVD", "Aronofsky", []string{"Ellen Burstyn", "Jared Leto", "Jennifer Connelly"}, "The drug-induced utopias of four Coney Island people are shattered.", 2, false},
		{"Spirited Away", 2001, "Animation", "Blu-ray", "Miyazaki", []string{"Rumi Hiiragi", "Miyu Irino", "Mari Natsuki"}, "A girl wanders into a world ruled by gods and witches.", 4, false},
		{"City of God", 2002, "Drama", "DVD", "Meirelles", []string{"Alexandre Rodrigues", "Leandro Firmino", "Phellipe Haagensen"}, "In the slums of Rio, two kids paths diverge as one becomes a photographer.", 3, false},
		{"The Usual Suspects", 1995, "Thriller", "VHS", "Singer", []string{"Kevin Spacey", "Gabriel Byrne", "Chazz Palminteri"}, "A sole survivor tells of the twisty events leading up to a gun battle.", 3, false},
		{"American History X", 1998, "Drama", "DVD", "Kaye", []string{"Edward Norton", "Edward Furlong", "Beverly D'Angelo"}, "A former neo-nazi skinhead tries to prevent his brother from going down the same path.", 3, false},
		{"The Green Mile", 1999, "Drama", "VHS", "Darabont", []string{"Tom Hanks", "Michael Clarke Duncan", "David Morse"}, "The lives of guards on Death Row are affected by one of their charges.", 4, false},
		{"District 9", 2009, "SciFi", "DVD", "Blomkamp", []string{"Sharlto Copley", "Jason Cope", "Nathalie Boltt"}, "Extraterrestrials forced to live in slum-like conditions on Earth.", 3, false},
		{"Leon The Professional", 1994, "Action", "VHS", "Besson", []string{"Jean Reno", "Gary Oldman", "Natalie Portman"}, "A professional assassin reluctantly cares for a 12-year-old neighbor girl.", 3, false},
		{"The Prestige", 2006, "Thriller", "DVD", "Nolan", []string{"Hugh Jackman", "Christian Bale", "Michael Caine"}, "Two stage magicians engage in a battle to create the ultimate illusion.", 3, false},
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
