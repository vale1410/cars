

n(1..1000). 

{ x(A,B) } :- n(A), n(B). 

y(B) :- n(B), not x(A,B), x(A,B+1). 
