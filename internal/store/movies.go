package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/thelastvideostore/internal/models"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) CreateMovie(movie *models.Movie) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket(bucketMovies)
		gb := tx.Bucket(bucketMoviesByGenre)
		tb := tx.Bucket(bucketMoviesByTitle)

		data, err := encode(movie)
		if err != nil {
			return err
		}

		if err := mb.Put([]byte(movie.ID), data); err != nil {
			return err
		}

		genreKey := fmt.Sprintf("%s:%s", movie.Genre, movie.ID)
		if err := gb.Put([]byte(genreKey), []byte(movie.ID)); err != nil {
			return err
		}

		titleKey := fmt.Sprintf("%s:%s", strings.ToLower(movie.Title), movie.ID)
		if err := tb.Put([]byte(titleKey), []byte(movie.ID)); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) GetMovieByID(id string) (*models.Movie, error) {
	var movie models.Movie
	err := s.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketMovies).Get([]byte(id))
		if data == nil {
			return fmt.Errorf("movie not found: %s", id)
		}
		return json.Unmarshal(data, &movie)
	})
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *Store) UpdateMovie(movie *models.Movie) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket(bucketMovies)
		gb := tx.Bucket(bucketMoviesByGenre)
		tb := tx.Bucket(bucketMoviesByTitle)

		oldData := mb.Get([]byte(movie.ID))
		if oldData == nil {
			return fmt.Errorf("movie not found: %s", movie.ID)
		}

		var old models.Movie
		if err := json.Unmarshal(oldData, &old); err != nil {
			return err
		}

		oldGenreKey := fmt.Sprintf("%s:%s", old.Genre, old.ID)
		gb.Delete([]byte(oldGenreKey))

		oldTitleKey := fmt.Sprintf("%s:%s", strings.ToLower(old.Title), old.ID)
		tb.Delete([]byte(oldTitleKey))

		data, err := encode(movie)
		if err != nil {
			return err
		}

		if err := mb.Put([]byte(movie.ID), data); err != nil {
			return err
		}

		newGenreKey := fmt.Sprintf("%s:%s", movie.Genre, movie.ID)
		if err := gb.Put([]byte(newGenreKey), []byte(movie.ID)); err != nil {
			return err
		}

		newTitleKey := fmt.Sprintf("%s:%s", strings.ToLower(movie.Title), movie.ID)
		if err := tb.Put([]byte(newTitleKey), []byte(movie.ID)); err != nil {
			return err
		}

		return nil
	})
}

func (s *Store) DeleteMovie(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		mb := tx.Bucket(bucketMovies)
		gb := tx.Bucket(bucketMoviesByGenre)
		tb := tx.Bucket(bucketMoviesByTitle)

		data := mb.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("movie not found: %s", id)
		}

		var movie models.Movie
		if err := json.Unmarshal(data, &movie); err != nil {
			return err
		}

		genreKey := fmt.Sprintf("%s:%s", movie.Genre, movie.ID)
		gb.Delete([]byte(genreKey))

		titleKey := fmt.Sprintf("%s:%s", strings.ToLower(movie.Title), movie.ID)
		tb.Delete([]byte(titleKey))

		return mb.Delete([]byte(id))
	})
}

func (s *Store) ListMovies(genre string, offset, limit int) ([]*models.Movie, int, error) {
	return s.ListMoviesFiltered(genre, "", offset, limit)
}

func (s *Store) ListMoviesFiltered(genre, mediaType string, offset, limit int) ([]*models.Movie, int, error) {
	var movies []*models.Movie
	total := 0

	err := s.db.View(func(tx *bolt.Tx) error {
		if genre != "" {
			gb := tx.Bucket(bucketMoviesByGenre)
			c := gb.Cursor()
			prefix := genre + ":"
			count := 0
			skipped := 0
			for k, v := c.Seek([]byte(prefix)); k != nil && strings.HasPrefix(string(k), prefix); k, v = c.Next() {
				total++
				if skipped < offset {
					skipped++
					continue
				}
				if count >= limit {
					continue
				}
				movieData := tx.Bucket(bucketMovies).Get(v)
				if movieData == nil {
					continue
				}
				var movie models.Movie
				if err := json.Unmarshal(movieData, &movie); err != nil {
					continue
				}
				movies = append(movies, &movie)
				count++
			}
		} else {
			mb := tx.Bucket(bucketMovies)
			c := mb.Cursor()
			count := 0
			skipped := 0
			for k, v := c.First(); k != nil; k, v = c.Next() {
				total++
				if skipped < offset {
					skipped++
					continue
				}
				if count >= limit {
					continue
				}
				var movie models.Movie
				if err := json.Unmarshal(v, &movie); err != nil {
					continue
				}
				if mediaType != "" && movie.MediaType != mediaType {
					total--
					continue
				}
				movies = append(movies, &movie)
				count++
			}
		}
		return nil
	})

	return movies, total, err
}

func (s *Store) SearchMoviesByPrefix(prefix string, limit int) ([]*models.Movie, error) {
	var movies []*models.Movie
	prefix = strings.ToLower(prefix)

	err := s.db.View(func(tx *bolt.Tx) error {
		tb := tx.Bucket(bucketMoviesByTitle)
		mb := tx.Bucket(bucketMovies)
		c := tb.Cursor()

		count := 0
		for k, v := c.Seek([]byte(prefix)); k != nil && strings.HasPrefix(string(k), prefix); k, v = c.Next() {
			if count >= limit {
				break
			}
			movieData := mb.Get(v)
			if movieData == nil {
				continue
			}
			var movie models.Movie
			if err := json.Unmarshal(movieData, &movie); err != nil {
				continue
			}
			movies = append(movies, &movie)
			count++
		}
		return nil
	})

	return movies, err
}

func (s *Store) GetNewReleases() ([]*models.Movie, error) {
	var movies []*models.Movie
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMovies)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var movie models.Movie
			if err := json.Unmarshal(v, &movie); err != nil {
				continue
			}
			if movie.IsNewRelease {
				movies = append(movies, &movie)
			}
		}
		return nil
	})
	return movies, err
}

func (s *Store) GetLastChanceMovies() ([]*models.Movie, error) {
	var movies []*models.Movie
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketMovies)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var movie models.Movie
			if err := json.Unmarshal(v, &movie); err != nil {
				continue
			}
			if movie.IsLastChance() {
				movies = append(movies, &movie)
			}
		}
		return nil
	})
	return movies, err
}
