% Complex model
% only by successor

% Definitions

car((C,K)) :- K = 1..M, class_count(C,M). 

size(N) :- N = #count{car(_)}.

% Model 

dont(A,B) :- 
    class_option(A,O), 
    option_max(O,1), 
    option_window(O,W), 
    W > 1,
    class_option(B,O). 


poss((A,K),(A,K+1)) :- 
    class_count(A,M),
    K = 1..M,
    K+1 <= M,
    not dont(A,A). 


poss((A,K),(B,L)) :- 
    class_count(A,M),
    class_count(B,N),
    A != B,
    K = 1..M,
    L = 1..N,
    not dont(A,B).

0 { succ(A,B) : poss(A,B) } 1 
    :- car(A). 

is_pred(B) :- succ(_,B). 
is_succ(A) :- succ(A,_). 

:- car(C), 2 {succ(_,C)} . 

is_first((A,1)) :- car((A,1)), not is_pred((A,1)). 
is_last((A,N)) :- class_count(A,N), not is_succ((A,N)). 

:- not 1 { is_first(_) } 1. 

is_reachable(A) :- is_first(A).  
is_reachable(B) :- is_reachable(A), succ(A,B). 
:- not is_reachable(C), car(C). 

% Options
is_option_counter((A,K),O,W-1,M-1) :- 
    car((A,K)), 
    class_option(A,O), 
    option_window(O,W),
    option_max(O,M), 
    (W,M) > (2,1). 
    

% pass if car does not have option
is_option_counter((B,L),O,W-1,M) :- 
    is_option_counter(A,O,W,M),
    succ(A,(B,L)),
    not class_option(B,O), 
    M >= 0,
    W > 0. 

% pass if car does have option
is_option_counter((B,L),O,W-1,M-1) :- 
    is_option_counter(A,O,W,M),
    succ(A,(B,L)),
    class_option(B,O),
    W > 0, M >= 0. 

:-  is_option_counter((A,_),O,0..W-2,-1),
    class_option(A,O), 
    option_window(O,W). 


%% break symmetries
%
%ord((A,1),B) :- succ((A,1),B). 
%ord((A,K),(C,L)) :- ord((A,K),B), succ(B,(C,L)), A != C. 
%ord((A,K),(A,K+1)) :- class_count(A,M), K = 1..M-1.
%
%:- succ((A,K),B), ord(B,(A,L)), L <= K. 

#hide.
#show is_first/1. 
#show succ/2. 
