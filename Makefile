.PHONY: tag-next tag-patch tag-minor tag-major

tag-next:
	@if [ -z "$(TYPE)" ]; then \
		echo "TYPE is required (patch|minor|major)"; \
		echo "example: make tag-next TYPE=patch"; \
		exit 2; \
	fi
	@./scripts/tag-next-version.sh "$(TYPE)"

tag-patch:
	@./scripts/tag-next-version.sh patch

tag-minor:
	@./scripts/tag-next-version.sh minor

tag-major:
	@./scripts/tag-next-version.sh major
