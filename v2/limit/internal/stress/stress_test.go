package stress

import (
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/stretchr/testify/require"
)

func TestStress(t *testing.T) {
	type config struct {
		CPUFactor  int           `env:"CQOS_STRESS_DATA_CPU_FACTOR"`
		DataAmount int           `env:"CQOS_STRESS_DATA_AMOUNT"`
		Duration   time.Duration `env:"CQOS_STRESS_DURATION"`
	}

	cfg := config{}

	err := env.Parse(&cfg)
	require.NoError(t, err)

	stress, err := New(cfg.CPUFactor, cfg.DataAmount)
	require.NoError(t, err)

	defer stress.Stop()

	time.Sleep(cfg.Duration)
}

func BenchmarkStress(b *testing.B) {
	stress, err := New(0, 0)
	require.NoError(b, err)

	defer stress.Stop()

	time.Sleep(10 * time.Second)
}
