package appcore_config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config *Configurations

// Configurations wraps all the config variables required by the service
type Configurations struct {
	// Develop or production
	Mode string

	// Gin Mode
	GinIsReleaseMode bool

	IsProfiling bool

	// Jwt
	JwtSecretKey string

	// Stripe
	StripeSecretKey string

	// Observability
	ObserveIsActive     bool
	ObserveOTLPEndpoint string
	ObserveInsecureMode string

	// Database
	PostgresConnString string
	MongoConnString    string

	// Redis
	RedisUrl  string
	RedisPass string

	// Message broker (rabbitmq)
	RabbitmqUrl string

	// Storage
	MinioURL           string
	MinioSSL           bool
	MinioAccessKey     string
	MinioSecretKey     string
	MinioBucketName    string
	MinioECMBucketName string

	// ENV
	ENVIRONMENT string

	// SMTP
	EmailSTMPHost     string
	EmailSMTPPort     string
	EmailSMTPUser     string
	EmailSMTPPassword string
	EmailSMTPFrom     string
}

// NewConfigurations returns a new Configuration object
func InitConfigurations() {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")    // บอก viper ให้มองหาไฟล์ config ใน current directory
	viper.SetConfigName(".env") // บอก viper ว่าชื่อไฟล์ config คือ .env
	viper.SetConfigType("env")  // บอก viper ว่าประเภทของไฟล์ config คือ env (สำหรับ .env files)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, relying on system environment variables.")
		} else {
			log.Fatalf("Fatal error config file: %s \n", err)
		}
	}

	viper.SetDefault("MODE", "develop")
	viper.SetDefault("GIN_IS_RELEASE_MODE", false)
	viper.SetDefault("IS_PROFILING", false)
	viper.SetDefault("OBSERVE_IS_ACTIVE", false)
	viper.SetDefault("OBSERVE_OTLP_ENDPOINT", "localhost:4317")
	viper.SetDefault("OBSERVE_INSECURE_MODE", "false")

	viper.SetDefault("MINIO_URL", "localhost:9000")
	viper.SetDefault("MINIO_SSL", false)
	viper.SetDefault("MINIO_BUCKET_NAME", "go-rebuild-bucket")
	viper.SetDefault("MINIO_ECM_BUCKET_NAME", "miniobucketecm")

	Config = &Configurations{
		Mode:             viper.GetString("MODE"),
		JwtSecretKey:        viper.GetString("JWT_SECRET_KEY"),
		ENVIRONMENT:      viper.GetString("ENVIRONMENT"),
		GinIsReleaseMode: viper.GetBool("GIN_IS_RELEASE_MODE"),
		IsProfiling:      viper.GetBool("IS_PROFILING"),

		StripeSecretKey: viper.GetString("STRIPE_SECRET_KEY"),

		ObserveIsActive:     viper.GetBool("OBSERVE_IS_ACTIVE"),
		ObserveOTLPEndpoint: viper.GetString("OBSERVE_OTLP_ENDPOINT"),
		ObserveInsecureMode: viper.GetString("OBSERVE_INSECURE_MODE"),

		PostgresConnString: viper.GetString("POSTGRES_URL"),
		MongoConnString:    viper.GetString("MONGO_URL"),

		RedisUrl:  viper.GetString("REDIS_URL"),
		RedisPass: viper.GetString("REDIS_PASS"),

		RabbitmqUrl: viper.GetString("RABBITMQ_URL"),

		MinioURL:           viper.GetString("MINIO_URL"),
		MinioSSL:           viper.GetBool("MINIO_SSL"),
		MinioAccessKey:     viper.GetString("MINIO_ROOT_USER"),
		MinioSecretKey:     viper.GetString("MINIO_ROOT_PASSWORD"),
		MinioBucketName:    viper.GetString("MINIO_BUCKET_NAME"),
		MinioECMBucketName: viper.GetString("MINIO_ECM_BUCKET_NAME"),

		EmailSTMPHost:     viper.GetString("EMAIL_SMTP_HOST"),
		EmailSMTPPort:     viper.GetString("EMAIL_SMTP_PORT"),
		EmailSMTPUser:     viper.GetString("EMAIL_SMTP_USER"),
		EmailSMTPPassword: viper.GetString("EMAIL_SMTP_PASSWORD"),
		EmailSMTPFrom:     viper.GetString("EMAIL_SMTP_FROM"),
	}
}
