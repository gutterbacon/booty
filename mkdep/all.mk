dep-all: dep-all-devtime dep-all-deploytime

## install all development tools
dep-all-devtime:
	@echo
	@echo -- dep-all-devtime  -- 
	@echo
	$(MAKE) dwn-print
	#$(MAKE) proto-dep:

## install all deployment tools
dep-all-deploytime:
	@echo
	@echo -- dep-all-deploytime -- 
	@echo
	$(MAKE) dwn-print
	$(MAKE) gor-dep