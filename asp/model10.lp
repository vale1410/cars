% revolutionary idea
%
% no car variables!

%Input: 

class(C) :- class_count(C,_). 
last(N) :- N = #sum[class_count(_,K) = K ].
pos(0..L) :- last(L). 

% mapping class to options

option_max(e,1). 
option_window(e,1). 


option(O) :- class_option(_,O). 

% model

counter(o,O,Q,M,N) :- 
    option_window(O,Q),
    option_max(O,M),
    N = #sum [class_count(C,T) = T : class_option(C,O)].

% all classes have option 0 

counter(o,0,1,1,N) :- 
    last(N). 

class_option(C,0) :- 
    class(C). 

% find the counter with the slowest decent and by 
% that determine the slope for the classes

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

% for each class take the smallest slope of its options. 
counter(c,C,Q,M,N) :- 
    class_count(C,N), 
    smallest(C,O,Q,M),
    last(L).
    
% Link between option and class
    
:-  option(O),
    pos(P), 
    not 0 [
    lower(o,O,P,LN) = -LN+1,
    is_count(o,O,P,_) = -1,
    lower(c,C,P,LN) = LN-1 : class_option(C,O), 
    is_count(c,C,P,_) = 1 : class_option(C,O) ] 0. 


%% symmetry breaking
%:-  last(L), 
%    is_car(C1,1),
%    is_car(C2,L),
%    C2 < C1. 
%
%%%%% what is missing: exactly one each day:
%%%%
%%%%% this needs improvement, flattening the encoding?
%%
%:- pos(P), P > 0, not 1 { is_car(_,P) } 1. 


%%%%%%%%%%%%%%%%%
% AtMostSeqCard %
%%%%%%%%%%%%%%%%%

% counter(T,I,Q,M,S), 
% with the meaning: 
% T: type of counter, car or option (c/o)
% I: identifier for the counter
% Q: fixed window size
% M: maximal hits per window
% S: exactly S occurences over the whole sequence

delta(Q,M,D,D) :-
    counter(_,_,Q,M,_), 
    D=0..M.

delta(Q,M,DP,M) :-
    counter(_,_,Q,M,_), 
    DP=M+1..Q-1.

lower(T,I,P,0) :- 
    pos(P),last(L),
    counter(T,I,Q,M,S),
    P <= L+Q-1-(Q*S)/M.

lower(T,I,P-DP,N-DN) :- 
    counter(T,I,Q,M,S), 
    last(L), 
    1 <= N-DN,
    pos(P-DP),
    delta(Q,M,DP,DN),
    D = 0..L/Q,
    N = S-(D*M),
    P = L-(D*Q).

upper(T,I,L-P,S-N+1) :- 
    counter(T,I,_,_,S), 
    last(L), 
    lower(T,I,P,N).

val(T,I,N) :- 
    counter(T,I,_,_,S), 
    N=1..S. 

domain(T,I,P,N) :- 
    lower(T,I,P,N).

domain(T,I,P,N) :- 
    upper(T,I,P,N). 

domain(T,I,P,N) :- 
    lower(T,I,P,LN), 
    upper(T,I,P,UN), 
    val(T,I,N), 
    LN < N, N < UN. 

0 { is_count(T,I,P,N) } 1 :- 
    domain(T,I,P,N). 

:-  not is_count(T,I,P,N), 
    lower(T,I,P,N). 

:-  is_count(T,I,P,N), 
    upper(T,I,P,N). 

% propagation

% right reasoning
:-  is_count(T,I,P,N), 
    domain(T,I,P+1,N),
    not is_count(T,I,P+1,N).

% left/down reasoning
:-  is_count(T,I,P+1,N+1),   
    domain(T,I,P,N),
    not is_count(T,I,P,N).

% steepness
:-  counter(T,I,Q,M,_), 
    is_count(T,I,P,N),
    domain(T,I,P-Q,N-M),
    not is_count(T,I,P-Q,N-M). 

#hide.
#show is_count/4. 

