PRE=pre
OUTPUT=$(NAME).pdf

all: pre

pre:
	pandoc --template beamer.template -i $(PRE).md -o $(PRE).tex
	perl -pe 's/\\section/\\end{frame}\n\\begin{frame}[fragile]\n\\frametitle/' -i $(PRE).tex
	perl -pe 's/\\end{document}/\\end{frame}\n\\end{document}/' -i $(PRE).tex
	sed -i 's/\[<+->\]//g' $(PRE).tex
	pdflatex $(PRE).tex
#	bibtex $(PRE)
#	pdflatex $(PRE).tex
#	pdflatex $(PRE).tex
	evince $(PRE).pdf 2>/dev/null &

clean: 
	rm -f pre.tex
	rm -f *.aux
	rm -f *.log
	rm -f *.html
	rm -f *.nav
	rm -f *.out
	rm -f *.snm
	rm -f *.toc
	rm -f *.vrb
