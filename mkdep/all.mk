dep-all: dep-all-print dep-all-install

## print all tools
dep-all-print:
	@echo
	@echo -- dep-all-print  -- 
	@echo
	$(MAKE) gor-print
	#$(MAKE) grafana-print

## install all tools
dep-all-install:
	@echo
	@echo -- dep-all-install -- 
	@echo
	$(MAKE) dwn-print
	$(MAKE) gor-dep
	#$(MAKE) grafana-dep