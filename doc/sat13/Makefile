sc2013.pdf: sc2013.tex sc2013.bib
	pdflatex sc2013.tex && \
	bibtex sc2013.aux && \
	pdflatex sc2013.tex && \
	pdflatex sc2013.tex

clean:
	rm -f sc2013.log sc2013.aux sc2013.bbl sc2013.blg sc2013.synctex.gz sc2013.pdf
