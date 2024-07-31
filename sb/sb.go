package sb

import (
	"fmt"
	"os"

	"github.com/nedpals/supabase-go"
)

var Client *supabase.Client

// Init initializes the supabase client using the SUPABASE_URL and SUPABASE_KEY environment variables.
// It returns an error if either of the environment variables is not set.
func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	sbKey := os.Getenv("SUPABASE_KEY")

	if sbHost == "" {
		return fmt.Errorf("SUPABASE_URL environment variable is not set")
	}

	if sbKey == "" {
		return fmt.Errorf("SUPABASE_KEY environment variable is not set")
	}

	Client = supabase.CreateClient(sbHost, sbKey)
	return nil
}
