% Simple Counter Model: 

% This model uses two variable: 
% is_car(Position,Type),
% is_option(Position,Option,Window,Counter)
% The constraints on the options are modelled with counters 
% The constraints on the cars are modelled with cardinality constraints
% Triggering is done by is_option/2

% Definitions

class(C) :- class_count(C,_). 
pos(1..N) :- N = #sum[class_count(_,K) = K ].
next(D,DD) :- pos(D;DD), DD = D+1. 

window(1 .. N) :- N = #max[ option_window(_,W) = W ].
window(1,D,D) :- pos(D).
window(N,D1,D2) :- window(N;N-1), window(N-1,D1,D2).
window(N,DN,D1) :- window(N;N-1), window(N-1,DN,D2),next(D1,D2). 
              
% Model 

1 { is_car(P,C) : class(C) } 1 :- pos(P). 

:- class_count(C,S), not S { is_car(P,C) } S. 

is_option(P,O) :- is_car(P,C), class_option(C,O).

% start counter with window
is_option_counter(P,O,1,1) :-
    is_option(P,O).

is_option_counter(P2,O,W+1,M) :-
    is_option_counter(P1,O,W,M),
    next(P1,P2),
    option_window(O,WM),
    W < WM. 

% increase counter if option

is_option_counter(P2,O,W+1,M+1) :-
    is_option_counter(P1,O,W,M), 
    option_window(O,WMax),
    option_max(O,MMax),
    W < WMax, M <= MMax,
    next(P1,P2),
    is_option(P2,O).

:-  is_option_counter(P,O,W,M+1),
    option_max(O,M),
    option_window(O,W).

#hide.
#show is_car/2. 
#show is_option/2. 
#show is_option_counter/4. 
