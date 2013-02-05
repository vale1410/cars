% revolutionary idea
%
% Probably the GAC encoding of the AtMostSequcenConstraint

%Input: 

class(C) :- class_count(C,_). 
last(N) :- N = #sum[class_count(_,K) = K ].
pos(0..L) :- last(L). 

% mapping class to options
% class_option(Car,O).

option(O) :- class_option(_,O). 

%%Idea and Model

counter(o,O,Q,M,N) :- 
    option_window(O,Q),
    option_max(O,M),
    N = #sum [class_count(C,T) = T : class_option(C,O)].

ord(C,Q1,M1,Q2,M2) :- 
    class_option(C,O1),
    class_option(C,O2),
    option_window(O1,Q1),
    option_max(O1,M1),
    option_window(O2,Q2),
    option_max(O2,M2), 
    M1*Q2 < M2*Q1. 

ord(C,Q1,M1,Q2,M2) :- 
    class_option(C,O1),
    class_option(C,O2),
    option_window(O1,Q1),
    option_max(O1,M1),
    option_window(O2,Q2),
    option_max(O2,M2), 
    M1*Q2 == M2*Q1, 
    Q1 < Q2. 

smallest(C,O,Q,M) :- 
    class_option(C,O),
    option_window(O,Q),
    option_max(O,M), 
    { ord(C,_,_,Q,M) } 0. 


counter(c,C,Q,M,N) :- 
    class_count(C,N), 
    smallest(C,O,Q,M),
    last(L).
    
% Link between option and class
%%%
%
%is_car(C,P+1) :- 
%    pos(P;P+1),
%    not is_count(c,C,P,L),
%    is_count(c,C,P+1,L).
%
%is_opt(O,P) :- 
%    class_option(C,O),
%    is_car(C,P).
%
%:-  pos(P;P+1),
%    not is_count(o,O,P,K),
%    is_count(o,O,P+1,K),
%    not is_opt(O,P+1). 

%%% what is missing: exactly one each day:
%%
%%% this needs improvement, flattening the encoding?

%:- pos(P), P > 0, not 1 { is_car(_,P) } 1. 

%
%%%%%%%%%%%%%%%%%
% AtMostSeqCard %
%%%%%%%%%%%%%%%%%

% counter(I,Q,M,S), 
% with the meaning: 
% I: identifier for the counter
% Q: fixed window size
% M: maximal value per window
% S: exactly S occurences over the whole sequence

delta(Q,M,D,D) :-
    counter(_,_,Q,M,_), 
    D=0..M.

delta(Q,M,DP,M) :-
    counter(_,_,Q,M,_), 
    DP=M+1..Q-1.

lower(T,I,P-DP,N-DN) :- 
    counter(T,I,Q,M,S), 
    last(L), 
    1 <= N-DN,
    pos(P-DP),
    delta(Q,M,DP,DN),
    D = 0..L/Q,
    N = S-(D*M),
    P = L-(D*Q).

val(T,I,N) :- 
    counter(T,I,_,_,S), 
    N=1..S. 

0 { is_count(T,I,P,N) } 1 :- 
    val(T,I,N),
    pos(P). 

%%%initalisation
%
%% lower left FALSE
%:-  counter(T,I,_,_,_),
%    is_count(T,I,0,1).
%    
%% upper right TRUE
%:-  last(P), 
%    counter(T,I,_,_,S),
%    not is_count(T,I,P,S).
%
%% cumulative structure
%:-  is_count(T,I,P,N), 
%    not is_count(T,I,P+1,N), 
%    pos(P;P+1).

% presumably not needed:
%%:-  is_count(T,I,P,N+1),   
%%    not is_count(T,I,P,N), 
%%    val(T,I,N+1;N). 
%
%% propagation
%:-  counter(T,I,_,M,_), 
%    K = 1..M, 
%    is_count(T,I,P+K,N+K),
%    pos(P),val(T,I,N),
%    not is_count(T,I,P,N).
%
%:-  counter(T,I,Q,M,_), 
%    val(T,I,N), pos(P),
%    not is_count(T,I,P,N),
%    K = M+1..Q,
%    is_count(T,I,P+K,N+M). 

%for debugging
%#minimize[is_count(T,I,P,N):val(T,I,N):pos(P)].

#hide.
#show is_car/2. 

