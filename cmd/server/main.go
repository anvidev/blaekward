package main

import "go.uber.org/zap"

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	s := &server{
		clients: make(map[string]*client),
		logger:  logger,
	}

	s.start()
}
