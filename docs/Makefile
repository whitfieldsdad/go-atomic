render:
	for dot_file in $(wildcard *.dot); do \
		png_file=$${dot_file%.dot}.png; \
		dot -Tpng -Gdpi=300 $$dot_file -o $$png_file; \
	done
