package main

type Config struct {
	Port           int
	AllowedOrigins []string
	World          WorldConfig
}

type WorldConfig struct {
	Width     int
	Height    int
	MaxPlayer int
	Food      int
}
