% Simplest Model: 

% This model has just one variable: is_car(Position,Type)
% The constraints regarding the option entirely cardinatlity
% constraints on basis of is_car. 

% Definitions

class(C) :- class_count(C,_). 
pos(1..N) :- N = #sum[class_count(_,K) = K ].
next(D,DD) :- pos(D;DD), DD = D+1. 

window(1 .. N) :- N = #max[ option_window(_,W) = W ].
              
% Model 

1 { is_car(P,C) : class(C) } 1 :- pos(P). 

:- class_count(C,S), not S { is_car(P,C) } S. 

:- pos(P), option_window(O,W), option_max(O,M), 
   M+1 { is_car(PP,C) : class_option(C,O) : PP = P..P+W-1}.

#hide.
#show is_car/2. 
