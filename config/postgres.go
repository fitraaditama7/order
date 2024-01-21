package config

type PostgresConfig struct {
	User                      string `mapstructure:"POSTGRES_USER"`
	DBName                    string `mapstructure:"POSTGRES_DB_NAME"`
	Password                  string `mapstructure:"POSTGRES_PASSWORD"`
	Host                      string `mapstructure:"POSTGRES_HOST"`
	Port                      string `mapstructure:"POSTGRES_PORT"`
	MigrationFilePath         string `mapstructure:"POSTGRES_MIGRATION_FILE_PATH"`
	MigrationOnStartupEnabled bool   `mapstructure:"POSTGRES_MIGRATION_ON_STARTUP_ENABLED"`
}
