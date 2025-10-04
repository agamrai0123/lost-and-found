package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func NewItemService() (*ItemService, error) {
	if err := readConfiguration(); err != nil {
		return nil, err
	}

	mode := gin.DebugMode
	gin.SetMode(mode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	s := &ItemService{
		storePath: "data/posts.json",
		imageDir:  "uploads", // directory to store uploaded images
		router:    router,
		logger:    getLogger(AppConfig.ServiceName),
		posts:     make(map[string]*Post),
		basePath:  "/api/v1",
	}

	// load existing posts if file exists
	if err := s.loadFromFile(); err != nil {
		// If the file doesn't exist, that's fine — we'll create on first save.
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("loading posts from file: %w", err)
		}
		s.logger.Warn().Msgf("posts file %q does not exist, starting fresh", s.storePath)
	} else {
		s.logger.Info().Msgf("loaded %d posts from %s", len(s.posts), s.storePath)
	}

	// register routes
	home := s.router.Group(s.basePath)
	s.routes(home)

	return s, nil
}

func (s *ItemService) Start(ctx context.Context) error {
	if s == nil {
		return errors.New("service is nil")
	}

	addr := ":8080"
	if AppConfig.ServerPort != "" {
		addr = fmt.Sprintf(":%s", AppConfig.ServerPort)
	}

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	// run server
	errCh := make(chan error, 1)
	go func() {
		s.logger.Info().Msgf("starting server on %s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// wait for context done or listen error
	select {
	case <-ctx.Done():
		s.logger.Info().Msg("shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		// ListenAndServe returned an error
		return err
	}
}
func (s *ItemService) Shutdown(ctx context.Context) error {
	if s == nil || s.httpServer == nil {
		return nil
	}
	s.logger.Info().Msg("attempting graceful shutdown")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("graceful shutdown failed, forcing close")
		return err
	}
	s.logger.Info().Msg("server stopped")
	return nil
}
