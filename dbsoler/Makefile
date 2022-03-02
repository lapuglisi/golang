DBSOLER=dbsoler
DBSOLER_SRC=dbsoler.go
DBSOLER_DEPS=$(DBSOLER_SRC) coredb/*.go server/*.go

PLUGINS_SRC=plugins/*.go
PLUGINS=$(PLUGINS_SRC:%.go=%.so)

all: plugins dbsoler

dbsoler: $(DBSOLER_DEPS)
	@echo -en "\033[32mBuilding '$(DBSOLER)'... \033[0m"
	@go build -o $(DBSOLER) $(DBSOLER_SRC); \
	[[ $$? -eq 0 ]] && echo "OK!" || echo "ERROR!"

plugins: $(PLUGINS_SRC)
	@echo -en "\033[32mBuilding plugins... \033[0m"
	@go build -buildmode=plugin $(PLUGINS_SRC); \
	[[ $$? -eq 0 ]] && echo "OK!" || echo "ERROR!"

clean:
	@rm -f $(DBSOLER) *.so