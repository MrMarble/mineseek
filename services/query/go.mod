module github.com/mrmarble/mineseek/query

go 1.18

require (
	github.com/labstack/echo/v4 v4.7.2
	github.com/mrmarble/mineseek/libs/database v0.0.0-20220526135802-0ac7ae60944f
	github.com/mrmarble/mineseek/libs/minecraft v0.0.0-20220526171628-a88414ad5f0a
	github.com/mrmarble/mineseek/libs/queue v0.0.0-00010101000000-000000000000
	github.com/rs/zerolog v1.26.1
)

require (
	github.com/adjust/rmq/v4 v4.0.5 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.3.2 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/labstack/gommon v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/xrjr/mcutils v1.2.1 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.9.1 // indirect
	go.opentelemetry.io/otel v0.13.0 // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220524220425-1d687d428aca // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/mrmarble/mineseek/libs/minecraft => ../../libs/minecraft

replace github.com/mrmarble/mineseek/libs/database => ../../libs/database

replace github.com/mrmarble/mineseek/libs/queue => ../../libs/queue
