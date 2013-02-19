
NAME=p
OUTPUT=$(NAME).pdf

all: show 

compile:
	pdflatex $(NAME).tex 

show: compile
	evince $(OUTPUT) 2>/dev/null

clean: 
	rm $(NAME).log
	rm $(NAME).aux
