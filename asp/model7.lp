% Simple Model with counters on cars:

% This model has variables: is_car(Position,Type) and
% is_option(Position,Option)
% The constraints regarding the cars is done by counters
% The constraints regarding the option is done by counters
% constraints on option is 

% doing some tweeking on it

% Definitions

class(C) :- class_count(C,_). 
first(1).
last(N) :- N = #sum[class_count(_,K) = K ].
pos(F..L) :- first(F), last(L). 
next(D,DD) :- pos(D;DD), DD = D+1. 

window(1 .. N) :- N = #max[ option_window(_,W) = W ].
              
% Model 

1 { is_car(P,C) : class(C) } 1 :- pos(P). 

% restriction on classes: 

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

%:-  class_count(C,M),
%    last(P),
%    is_class_counter(P,C,M+1).
%
%:-  class_count(C,M),
%    last(P),
%    not is_class_counter(P,C,M).

% restriction on options:

is_option(P,O) :- is_car(P,C), class_option(C,O).

% special case when option_max is 1

:-  is_option(P,O), 
    option_max(O,1), 
    option_window(O,W), 
    PP = P-W+1..P+W-1,
    pos(PP), 
    P != PP, 
    is_option(PP,O). 

% start counter with window where option_max > 1 

is_option_counter(P,O,1,1) :-
    option_max(O,M),
    M > 1, 
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
    W < WMax, M < MMax,
    next(P1,P2),
    is_option(P2,O).

:-  is_option_counter(P1,O,W,M),
    next(P1,P2),
    is_option(P2,O),
    option_max(O,M),
    option_window(O,WMax), 
    W < WMax. 

% symmetry breaking

:- first(F), last(L), is_car(F,A), is_car(L,B), 
   A > B. 

% redundant constraints on options

:-  option_max(O,_), 
    N = #sum [class_count(A,C) = C : class_option(A,O)], 
    not N { is_option(_,O) } N. 


%#hide.
%#show is_car/2. 
%#show is_option/2. 
%#show class/1. 
%#show next/2. 
%#show first/1. 
%#show option/1. 
                                                                                              
