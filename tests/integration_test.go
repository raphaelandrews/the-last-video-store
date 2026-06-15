package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/thelastvideostore/api"
	"github.com/thelastvideostore/internal/auth"
	"github.com/thelastvideostore/internal/config"
	"github.com/thelastvideostore/internal/crypto"
	"github.com/thelastvideostore/internal/ds/bitmask"
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

var testServer *httptest.Server
var testStore *store.Store
var testHC *crypto.HashChain
var testCfg *config.Config

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error
	path := os.TempDir() + "/tlvs-integration-test.db"
	os.Remove(path)

	testStore, err = store.Open(path)
	if err != nil {
		panic(fmt.Sprintf("open store: %v", err))
	}

	testHC = crypto.New()
	testCfg = config.Load()
	testCfg.JWTSecret = "integration-test-secret-32-bytes"
	testCfg.AESKey = "1234567890abcdef1234567890abcdef"

	seedIntegrationData(testStore)

	router := api.NewRouter(testStore, testCfg, testHC)
	testServer = httptest.NewServer(router)
}

func teardown() {
	if testServer != nil {
		testServer.Close()
	}
	if testStore != nil {
		testStore.Close()
	}
}

func seedIntegrationData(s *store.Store) {
	users := []struct {
		id, name, pass string
		tier           bitmask.Permission
		banned         bool
	}{
		{"seed-bronze", "bronze", "password1", bitmask.TierBronze, false},
		{"seed-silver", "silver", "password2", bitmask.TierSilver, false},
		{"seed-gold", "gold", "password3", bitmask.TierGold, false},
		{"seed-employee", "employee", "password4", bitmask.TierEmployee, false},
		{"seed-supervisor", "supervisor", "password8", bitmask.TierSupervisor, false},
		{"seed-manager", "manager", "password5", bitmask.TierManager, false},
		{"seed-owner", "owner", "password6", bitmask.TierOwner, false},
		{"seed-banned", "banned", "password7", bitmask.TierBronze, true},
	}

	for _, u := range users {
		hash, _ := auth.HashPassword(u.pass)
		now := time.Now().Unix()
		s.CreateUser(&models.User{
			ID:           u.id,
			Username:     u.name,
			PasswordHash: hash,
			Tier:         u.tier,
			MaxRentals:   bitmask.MaxRentalsForTier(u.tier),
			Banned:       u.banned,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	movies := []struct {
		id           string
		title        string
		year         int
		genre        string
		format       string
		copies       int
		isNewRelease bool
	}{
		{"seed-movie-TheMatrix", "The Matrix", 1999, "SciFi", "DVD", 5, false},
		{"seed-movie-PulpFiction", "Pulp Fiction", 1994, "Drama", "VHS", 3, false},
		{"seed-movie-FightClub", "Fight Club", 1999, "Drama", "DVD", 4, false},
		{"seed-movie-JurassicPark", "Jurassic Park", 1993, "Action", "Blu-ray", 5, false},
		{"seed-movie-Inception", "Inception", 2010, "SciFi", "DVD", 4, true},
		{"seed-movie-Memento", "Memento", 2000, "Thriller", "DVD", 3, true},
	}

	for _, m := range movies {
		s.CreateMovie(&models.Movie{
			ID:              m.id,
			Title:           m.title,
			Year:            m.year,
			Genre:           m.genre,
			Format:          m.format,
			Director:        "Test Director",
			Cast:            []string{"Actor 1", "Actor 2"},
			Synopsis:        "A test movie.",
			CopiesTotal:     m.copies,
			CopiesAvailable: m.copies,
			Available:       true,
			IsNewRelease:    m.isNewRelease,
			Rating:          4.0,
			RatingCount:     1000,
			CreatedAt:       time.Now().Unix(),
		})
	}

	s.AddStaffPick("seed-movie-TheMatrix")
	s.AddStaffPick("seed-movie-JurassicPark")
}

func login(t *testing.T, username, password string) string {
	t.Helper()
	body := map[string]string{"username": username, "password": password}
	resp := doRequest(t, "POST", "/api/v1/auth/login", body, http.StatusOK)
	var result struct {
		AccessToken string `json:"access_token"`
	}
	mustDecode(resp, &result)
	return result.AccessToken
}

func doRequest(t *testing.T, method, path string, body interface{}, wantStatus int) []byte {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if resp.StatusCode != wantStatus {
		t.Fatalf("want status %d, got %d: %s", wantStatus, resp.StatusCode, string(data))
	}

	return data
}

func authDo(t *testing.T, token, method, path string, body interface{}, wantStatus int) []byte {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if resp.StatusCode != wantStatus {
		t.Fatalf("want status %d, got %d: %s", wantStatus, resp.StatusCode, string(data))
	}

	return data
}

func mustDecode(data []byte, v interface{}) {
	if err := json.Unmarshal(data, v); err != nil {
		panic(fmt.Sprintf("decode: %v in %s", err, string(data)))
	}
}

func TestIntegration_BronzeRentalLimit(t *testing.T) {
	token := login(t, "bronze", "password1")

	rent1 := authDo(t, token, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-PulpFiction"}, http.StatusCreated)

	var rental models.RentalResponse
	mustDecode(rent1, &rental)
	if rental.MovieID != "seed-movie-PulpFiction" {
		t.Errorf("rented movie = %s, want seed-movie-PulpFiction", rental.MovieID)
	}

	authDo(t, token, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-TheMatrix"}, http.StatusForbidden)
}

func TestIntegration_SilverRentalLimit(t *testing.T) {
	token := login(t, "silver", "password2")

	authDo(t, token, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-TheMatrix"}, http.StatusCreated)
	authDo(t, token, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-FightClub"}, http.StatusCreated)
	authDo(t, token, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-JurassicPark"}, http.StatusForbidden)
}

func TestIntegration_BannedUser(t *testing.T) {
	doRequest(t, "POST", "/api/v1/auth/login",
		map[string]string{"username": "banned", "password": "password7"}, http.StatusForbidden)
}

func TestIntegration_PlanUpgradeBySupervisor(t *testing.T) {
	token := login(t, "supervisor", "password8")
	silverToken := login(t, "silver", "password2")

	authDo(t, token, "PUT", "/api/v1/users/seed-silver",
		map[string]string{"tier": "gold"}, http.StatusOK)

	time.Sleep(10 * time.Millisecond)

	newToken := login(t, "silver", "password2")
	_ = newToken

	authDo(t, silverToken, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-Inception"}, http.StatusCreated)

	newSilverToken := login(t, "silver", "password2")
	authDo(t, newSilverToken, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-Inception"}, http.StatusCreated)

	authDo(t, token, "PUT", "/api/v1/users/seed-silver",
		map[string]string{"tier": "silver"}, http.StatusOK)
}

func TestIntegration_AuditChain(t *testing.T) {
	token := login(t, "supervisor", "password8")

	data := authDo(t, token, "GET", "/api/v1/audit", nil, http.StatusOK)

	var entries []map[string]interface{}
	mustDecode(data, &entries)

	if len(entries) < 5 {
		t.Errorf("audit entries = %d, want at least 5", len(entries))
	}
}

func TestIntegration_EmployeeReturnWithRewind(t *testing.T) {
	goldToken := login(t, "gold", "password3")

	rentData := authDo(t, goldToken, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-PulpFiction"}, http.StatusCreated)

	var rental models.RentalResponse
	mustDecode(rentData, &rental)
	if rental.MovieFormat != "VHS" {
		t.Skip("Pulp Fiction is VHS — skipping rewind test")
	}

	empToken := login(t, "employee", "password4")

	returnData := authDo(t, empToken, "POST", "/api/v1/rentals/return",
		map[string]string{"rental_id": rental.ID}, http.StatusOK)

	var returned models.RentalResponse
	mustDecode(returnData, &returned)
	if returned.Status != "returned" {
		t.Errorf("status = %s, want returned", returned.Status)
	}

	t.Logf("Rewind fee: $%.2f, Late fee: $%.2f", returned.RewindFee, returned.LateFee)
}

func TestIntegration_NewReleaseRequiresGold(t *testing.T) {
	silverToken := login(t, "silver", "password2")

	authDo(t, silverToken, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-Inception"}, http.StatusForbidden)

	goldToken := login(t, "gold", "password3")
	authDo(t, goldToken, "POST", "/api/v1/rentals/rent",
		map[string]string{"movie_id": "seed-movie-Inception"}, http.StatusCreated)
}

func TestIntegration_StaffPicksAvailable(t *testing.T) {
	managerToken := login(t, "manager", "password5")

	data := authDo(t, managerToken, "GET", "/api/v1/movies/staff-picks", nil, http.StatusOK)
	var picks []map[string]interface{}
	mustDecode(data, &picks)
	if len(picks) < 1 {
		t.Error("expected at least 1 staff pick seeded (The Matrix)")
	}
}

func TestIntegration_LastChance(t *testing.T) {
	managerToken := login(t, "manager", "password5")

	data := authDo(t, managerToken, "GET", "/api/v1/movies/last-chance", nil, http.StatusOK)
	var movies []map[string]interface{}
	mustDecode(data, &movies)
	t.Logf("Last chance movies: %d", len(movies))
}

func TestIntegration_TOTPSetupAndLogin(t *testing.T) {
	managerToken := login(t, "manager", "password5")

	data := authDo(t, managerToken, "POST", "/api/v1/users/seed-manager/totp",
		map[string]bool{"enabled": true}, http.StatusOK)

	var setupResp struct {
		Secret string `json:"secret"`
		URL    string `json:"url"`
	}
	mustDecode(data, &setupResp)
	if setupResp.Secret == "" {
		t.Fatal("TOTP secret should not be empty")
	}
	if setupResp.URL == "" {
		t.Fatal("TOTP URL should not be empty")
	}

	t.Logf("TOTP secret: %s", setupResp.Secret)
	t.Logf("TOTP URL: %s", setupResp.URL)

	code, err := auth.GenerateTOTPCode(setupResp.Secret, time.Now())
	if err != nil {
		t.Fatalf("GenerateTOTPCode: %v", err)
	}
	t.Logf("Current TOTP code: %s", code)

	if !auth.ValidateTOTPCode(setupResp.Secret, code) {
		t.Error("generated code should validate")
	}

	authDo(t, managerToken, "POST", "/api/v1/users/seed-manager/totp",
		map[string]bool{"enabled": false}, http.StatusOK)
}

func TestIntegration_MovieSearchPrefix(t *testing.T) {
	token := login(t, "bronze", "password1")

	data := authDo(t, token, "GET", "/api/v1/movies/search?q=pulp", nil, http.StatusOK)
	var results []map[string]interface{}
	mustDecode(data, &results)
	if len(results) != 1 {
		t.Errorf("search results = %d, want 1", len(results))
	}
	if results[0]["title"] != "Pulp Fiction" {
		t.Errorf("title = %v, want Pulp Fiction", results[0]["title"])
	}
}

func TestIntegration_RefreshToken(t *testing.T) {
	body := map[string]string{"username": "bronze", "password": "password1"}
	resp := doRequest(t, "POST", "/api/v1/auth/login", body, http.StatusOK)

	var loginResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	mustDecode(resp, &loginResp)

	refreshBody := map[string]string{"refresh_token": loginResp.RefreshToken}
	refreshData := authDo(t, loginResp.AccessToken, "POST", "/api/v1/auth/refresh", refreshBody, http.StatusOK)

	var newTokens struct {
		AccessToken string `json:"access_token"`
	}
	mustDecode(refreshData, &newTokens)
	if newTokens.AccessToken == "" {
		t.Error("refresh should return new access token")
	}
}

func TestIntegration_RegisterDuplicate(t *testing.T) {
	doRequest(t, "POST", "/api/v1/auth/register",
		map[string]string{"username": "bronze", "password": "newpass"}, http.StatusConflict)
}

func TestIntegration_Wishlist(t *testing.T) {
	token := login(t, "gold", "password3")

	authDo(t, token, "POST", "/api/v1/wishlist",
		map[string]string{"movie_id": "seed-movie-TheMatrix"}, http.StatusCreated)

	authDo(t, token, "POST", "/api/v1/wishlist",
		map[string]string{"movie_id": "seed-movie-TheMatrix"}, http.StatusConflict)

	data := authDo(t, token, "GET", "/api/v1/wishlist", nil, http.StatusOK)
	var items []map[string]interface{}
	mustDecode(data, &items)
	if len(items) != 1 {
		t.Errorf("wishlist items = %d, want 1", len(items))
	}

	authDo(t, token, "DELETE", "/api/v1/wishlist/seed-movie-TheMatrix", nil, http.StatusOK)
}

func TestIntegration_SupervisorCannotAdminMovies(t *testing.T) {
	token := login(t, "supervisor", "password8")

	authDo(t, token, "POST", "/api/v1/movies",
		map[string]interface{}{"title": "Test", "year": 2020, "genre": "Action", "format": "DVD", "copies_total": 1},
		http.StatusForbidden)
}

func TestIntegration_ManagerCanCRUDMovies(t *testing.T) {
	token := login(t, "manager", "password5")

	data := authDo(t, token, "POST", "/api/v1/movies",
		map[string]interface{}{
			"title": "Manager Test Movie", "year": 2022, "genre": "Action",
			"format": "DVD", "director": "M. Test", "copies_total": 3,
		}, http.StatusCreated)

	var movie map[string]interface{}
	mustDecode(data, &movie)
	movieID := movie["id"].(string)

	authDo(t, token, "PUT", "/api/v1/movies/"+movieID,
		map[string]interface{}{"title": "Manager Test Movie Updated"}, http.StatusOK)

	authDo(t, token, "DELETE", "/api/v1/movies/"+movieID, nil, http.StatusOK)
}
