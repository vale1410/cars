

car(1..300). 

{ succ(A,B) } :- car(A;B), A < B. 

:- car(A), 3 { succ(A,_), succ(_,A) }. 

in(A) :- succ(_,A). 
out(A) :- succ(A,_). 

middle(A) :- in(A), out(A). 

left_edge(A) :- car(A), not middle(A), 2 { succ(A,_) }. 
right_edge(A) :- car(A), not middle(A), 2 { succ(_,A) }. 

edge(A) :- left_edge(A).
edge(A) :- right_edge(A).

one(A) :- car(A), not middle(A), not edge(A). 

:- not 2 { one(_) } 2. 

first(A) :- one(A;B), A < B. 

reached(A) :- first(A). 
reached(B) :- succ(A,B), reached(A). 
reached(B) :- succ(B,A), reached(A). 

:- car(A), not reached(A). 

left(A) :- in(A), first(A). 
right(A) :- out(A), first(A).

left(A) :- left(B), succ(A,B), middle(A). 
right(B) :- right(A), succ(A,B), middle(B). 

%left(C) :- right(A), succ(A,B), succ(C,B), A!=C. 
%right(C) :- left(A), succ(B,A), succ(B,C), A!=C. 
