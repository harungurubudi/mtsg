package errorhandler

// Config holds error handler configuration
type Config struct {
	Environment string `json:"environment"`
	ShowDetails bool   `json:"show_details"`
	LogErrors   bool   `json:"log_errors"`
}

// DefaultConfig returns default error handler configuration
func DefaultConfig() *Config {
	return &Config{
		Environment: "development",
		ShowDetails: true,
		LogErrors:   true,
	}
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
