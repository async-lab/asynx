PROJECT_NAME := asynx

GOZERO_MODULES := custodian

define GENERATE_BY_GOCTL
	goctl api go -api ./api/$(1).api -dir ./$(1)
	goctl api swagger -api ./api/$(1).api -dir ./swagger
endef



# ------------------------------------------------------------------

.PHONY: all generate

all: generate

generate:
	$(foreach module,$(GOZERO_MODULES),$(call GENERATE_BY_GOCTL,$(module)))