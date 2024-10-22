MOCK_OUT=pkg/mock/pkg_gen.go
MOCK_SRC=./pkg/domain/interfaces
MOCK_INTERFACES=CloudStorageClient BigQueryClient PubSubClient PolicyClient CacheService NotifyService

all: mock

mock: $(MOCK_OUT)

$(MOCK_OUT): $(MOCK_SRC)/*
	go run github.com/matryer/moq@latest -pkg mock -out $(MOCK_OUT) $(MOCK_SRC) $(MOCK_INTERFACES)

clean:
	rm -f $(MOCK_OUT)
