
NAME=p
OUTPUT=$(NAME).pdf

all: show 

compile:
	pdflatex -synctex=1 $(NAME).tex 

show: compile
	evince $(OUTPUT) 2>/dev/null &

clean: 
	rm $(NAME).log
	rm $(NAME).aux
