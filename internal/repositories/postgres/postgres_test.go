package postgres

import (
	"database/sql"
	"effective-mobile-test-task/internal/lib/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type PostgresTestSuite struct {
	suite.Suite
	db  *sql.DB
	cfg *config.TestConfig
}

func (p *PostgresTestSuite) SetupSuite() {
	envPath := os.Getenv("TEST_ENV_PATH")

	err := godotenv.Load(envPath)
	require.NoError(p.T(), err)

	p.cfg, err = config.LoadTest()
	require.NoError(p.T(), err)
	require.NotNil(p.T(), p.cfg)

	p.db, err = New(p.cfg.PG)
	require.NoError(p.T(), err)
	require.NotNil(p.T(), p.db)

	require.NoError(p.T(), MigrateTop(p.db, p.cfg.MigrationsPath))
}

func (p *PostgresTestSuite) TearDownTest() {
	require.NoError(p.T(), DropMigrations(p.db, p.cfg.MigrationsPath))
	require.NoError(p.T(), MigrateTop(p.db, p.cfg.MigrationsPath))
}

func (p *PostgresTestSuite) TearDownSuite() {
	require.NoError(p.T(), DropMigrations(p.db, p.cfg.MigrationsPath))
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
