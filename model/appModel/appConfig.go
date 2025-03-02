package appModel

import "time"

type AppConfig struct {
	App struct {
		Name string
		Host string
		Port string
	}
	Mysql struct {
		Dsn             string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime string
	}
	Redis struct {
		TokenDB      int
		CaptchaDB    int
		MsgDB        int
		MsgHistoryDB int
		Host         string
		Port         string
		Password     string
		Timeout      time.Duration
		Pool         struct {
			MaxActive int
			MaxIdle   int
			MinIdle   int
			MaxWait   time.Duration
		}
	}
	Mongo struct {
		Host     string
		Port     string
		Username string
		Password string
	}
	Es struct {
		Host string
		Port string
	}
	Kafka struct {
		Network   string
		Host      string
		Port      string
		NoteLikes struct {
			Topic      string
			Partitions int
		}
		NoteComments struct {
			Topic      string
			Partitions int
		}
	}
	Oss struct {
		AvatarBucket   string
		NotePicsBucket string
		StyleBucket    string
		EndPoint       string
		Region         string
	}
	Captcha struct {
		EndPoint string
		SceneId  string
	}
}
