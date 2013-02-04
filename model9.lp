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

% options properties: 
% option_window(Option,Window)
% option_max(Option,Max)

%Idea and Model

counter(o(O),Q,M,N) :- 
    option_window(O,Q),
    option_max(O,M),
    N = #sum [class_count(C,T) = T : class_option(C,O)].

counter(c(C),1,1,N) :- 
    class_count(C,N). 
    
% Link between option and class

is_opt(o(O),P) :- 
    class_option(C,O),
    pos(P;P+1),
    not is_c(c(C),P,L),
    is_c(c(C),P+1,L). 

:-  pos(P;P+1),
    is_opt(O,P),
    is_c(O,P,K),
    not is_c(O,P+1,K+1).

:-  pos(P;P+1),
    not is_c(o(O),K,K),
    is_c(o(O),P+1,K), 
    not is_opt(o(O),P+1). 

% what is missing: exactly one each day:

ex(P,K) :- 
    pos(P).
    K=1..P. 

% at most N for each counter missing!

% do we need it?
% no because we removed the need for them by 
% omitting the variables 


%%%%%%%%%%%%%%1%%
% AtMostSeqCard %
%%%%%%%%%%%%%%%%%

% counter(I,Q,M,S), 
% with the meaning: 
% I: identifier for the counter
% Q: fixed window size
% M: maximal value per window
% S: exactly S occurences over the whole sequence

%%test
%counter(t,7,7,7).
%
%% Encoding
%val(I,N) :- 
%    counter(I,_,_,S), 
%    N=1..S. 
%
%0 { is_c(I,P,N) } 1 :- 
%    val(I,N),
%    pos(P). 
%
%%initalisation
%
%% lower left FALSE
%:-  counter(I,_,_,_),
%    is_c(I,0,1).
%    
%% upper right TRUE
%:-  last(P), 
%    counter(I,_,_,S),
%    not is_c(I,P,S).
%
%% cumulative structure
%:-  is_c(I,P,N), 
%    not is_c(I,P+1,N), 
%    pos(P;P+1).
%
%:-  is_c(I,P,N+1),   
%    not is_c(I,P,N), 
%    val(I,N+1;N). 
%
%% propagation
%:-  counter(I,_,M,_), 
%    pos(P),
%    val(I,N;N+K),
%    not is_c(I,P,N),
%    K = 1..M, 
%    is_c(I,P+K,N+K).
%
%:-  counter(I,Q,M,_), 
%    val(I,N;N+M),
%    pos(P),
%    not is_c(I,P,N),
%    K = M+1..Q,
%    is_c(I,P+K,N+M). 
%
%%for debugging
%%#minimize[is_c(I,P,N):val(I,N):pos(P)].
%
%#hide.
%#show is_c/3. 

