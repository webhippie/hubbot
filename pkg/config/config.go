package config

import (
	"flag"
	"os"
	"strconv"
)

type GitHubConfig struct {
	WebhookSecret string
}

type DroneConfig struct {
	Server string
	Token  string
}

type Config struct {
	GitHub    GitHubConfig
	Drone     DroneConfig
	DebugMode bool
}

// New returns a new Config struct
func New() *Config {

	// Read command line flags
	var hub_webhook_secret = flag.String("hub_webhook_secret", "", "Github webhook secret")
	var drone_server = flag.String("drone_server", "https://cloud.drone.io/", "Drone server")
	var drone_token = flag.String("drone_token", "", "Drone token")
	var debug_mode = flag.Bool("debug_mode", false, "Enables debug outputs")

	flag.Parse()

	// Overwrite with env variables if existing
	if value, exists := os.LookupEnv("HUB_WEBHOOK_SECRET"); exists {
		*hub_webhook_secret = value
	}

	if value, exists := os.LookupEnv("DRONE_SERVER"); exists {
		*drone_server = value
	}

	if value, exists := os.LookupEnv("DRONE_TOKEN"); exists {
		*drone_token = value
	}

	if value, exists := os.LookupEnv("DEBUG_MODE"); exists {
		if val, err := strconv.ParseBool(value); err == nil {
			*debug_mode = val
		}
	}

	return &Config{
		GitHub: GitHubConfig{
			WebhookSecret: *hub_webhook_secret,
		},
		Drone: DroneConfig{
			Server: *drone_server,
			Token:  *drone_token,
		},
		DebugMode: *debug_mode,
	}
}
