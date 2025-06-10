package appcore_config

// import (
// 	"github.com/spf13/viper"
// )

// var Config *Configurations

// // Configurations wraps all the config variables required by the service
// type Configurations struct {
// 	//develop or production
// 	Mode string

// 	//Gin Mode
// 	GinIsReleaseMode bool

// 	IsProfiling bool

// 	//observability
// 	ObserveIsActive     bool
// 	ObserveOTLPEndpoint string
// 	ObserveInsecureMode string

// 	//database
// 	PostgresConnString string
// 	MongoConnString    string

// 	//Redis
// 	RedisUrl  string
// 	RedisPass string

// 	// Message broker (rabbitmq)
// 	RabbitmqUrl  string

// 	//Storage
// 	MinioURL           string
// 	MinioSSL           bool
// 	MinioAccessKey     string
// 	MinioSecretKey     string
// 	MinioBucketName    string
// 	MinioECMBucketName string

// 	SecretKey string

// 	// ENV
// 	ENVIRONMENT string

// 	// STMP
// 	EmailSTMPHost     string
// 	EmailSMTPPort     string
// 	EmailSMTPUser     string
// 	EmailSMTPPassword string
// 	EmailSMTPFrom     string
// }

// // NewConfigurations returns a new Configuration object
// func InitConfigurations() {
// 	viper.AutomaticEnv()
// 	viper.SetDefault("MODE", )
// 	viper.SetDefault("GIN_IS_RELEASE_MODE", )
// 	viper.SetDefault("IS_PROFILING", )
// 	viper.SetDefault("OBSERVE_IS_ACTIVE", )
// 	viper.SetDefault("OBSERVE_OTLP_ENDPOINT", )
// 	viper.SetDefault("OBSERVE_INSECURE_MODE",  )
// 	viper.SetDefault("POSTGRES_URL", )
// 	viper.SetDefault("MONGO_URL", )
// 	viper.SetDefault("REDIS_URL", )
// 	viper.SetDefault("REDIS_PASS", )
// 	viper.SetDefault("RABBITMQ_URL", )
// 	viper.SetDefault("MINIO_URL", )
// 	viper.SetDefault("MINIO_SSL", )
// 	viper.SetDefault("MINIO_ACCESS_KEY", )
// 	viper.SetDefault("MINIO_SECRET_KEY", )
// 	viper.SetDefault("MINIO_BUCKET_NAME", )
// 	viper.SetDefault("MINIO_ECM_BUCKET_NAME", )
// 	viper.SetDefault("SECRET_KEY", )
// 	viper.SetDefault("ENVIRONMENT", )
// 	viper.SetDefault("EMAIL_SMTP_HOST", )
// 	viper.SetDefault("EMAIL_SMTP_PORT", )
// 	viper.SetDefault("EMAIL_SMTP_USER", )
// 	viper.SetDefault("EMAIL_SMTP_PASSWORD", )
// 	viper.SetDefault("EMAIL_SMTP_FROM", )

// 	Config = &Configurations{
// 		Mode:                viper.GetString("MODE"),
// 		GinIsReleaseMode:    viper.GetBool("GIN_IS_RELEASE_MODE"),
// 		IsProfiling:         viper.GetBool("IS_PROFILING"),
// 		ObserveIsActive:     viper.GetBool("OBSERVE_IS_ACTIVE"),
// 		ObserveOTLPEndpoint: viper.GetString("OBSERVE_OTLP_ENDPOINT"),
// 		ObserveInsecureMode: viper.GetString("OBSERVE_INSECURE_MODE"),
// 		PostgresConnString:  viper.GetString("POSTGRES_URL"),
// 		MongoConnString:     viper.GetString("MONGO_URL"),
// 		RedisUrl:            viper.GetString("REDIS_URL"),
// 		RedisPass:           viper.GetString("REDIS_PASS"),
// 		RabbitmqUrl:         viper.GetString("Rabbitmq_URL"),
// 		MinioURL:            viper.GetString("MINIO_URL"),
// 		MinioSSL:            viper.GetBool("MINIO_SSL"),
// 		MinioAccessKey:      viper.GetString("MINIO_ACCESS_KEY"),
// 		MinioSecretKey:      viper.GetString("MINIO_SECRET_KEY"),
// 		MinioBucketName:     viper.GetString("MINIO_BUCKET_NAME"),
// 		MinioECMBucketName:  viper.GetString("MINIO_ECM_BUCKET_NAME"),
// 		SecretKey:           viper.GetString("SECRET_KEY"),
// 		ENVIRONMENT:         viper.GetString("ENVIRONMENT"),
// 		EmailSTMPHost:       viper.GetString("EMAIL_SMTP_HOST"),
// 		EmailSMTPPort:       viper.GetString("EMAIL_SMTP_PORT"),
// 		EmailSMTPUser:       viper.GetString("EMAIL_SMTP_USER"),
// 		EmailSMTPPassword:   viper.GetString("EMAIL_SMTP_PASSWORD"),
// 		EmailSMTPFrom:       viper.GetString("EMAIL_SMTP_FROM"),
// 	}
// }
