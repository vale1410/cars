% Simple Model with counters on cars:

% This model has just one variable: is_car(Position,Type)
% The constraints regarding the option entirely cardinatlity
% The constraints regarding the cars is done by counters
% constraints on basis of is_car. 

% Definitions

class(C) :- class_count(C,_). 
first(1).
last(N) :- N = #sum[class_count(_,K) = K ].
pos(F..L) :- first(F), last(L). 
next(D,DD) :- pos(D;DD), DD = D+1. 

window(1 .. N) :- N = #max[ option_window(_,W) = W ].
              
% Model 

1 { is_car(P,C) : class(C) } 1 :- pos(P). 

:- pos(P), option_window(O,W), option_max(O,M), 
   M+1 { is_car(PP,C) : class_option(C,O) : PP = P..P+W-1}.

is_class_counter(P,C,1) :- 
    is_car(P,C).

is_class_counter(P2,C,M) :-
    is_class_counter(P1,C,M),
    class_count(C,MM),
    M <= MM+1,
    next(P1,P2).

is_class_counter(P2,C,M+1) :-
    is_class_counter(P1,C,M),
    class_count(C,MM),
    M <= MM,
    next(P1,P2),
    is_car(P2,C).

:-  class_count(C,M),
    last(P),
    is_class_counter(P,C,M+1).

:-  class_count(C,M),
    last(P),
    not is_class_counter(P,C,M).


#hide.
#show is_car/2. 
