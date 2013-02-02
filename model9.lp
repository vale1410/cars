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

% We introduce an auxiliary option e
% Exactly one car per position

option_window(e,1). 
option_max(e,1). 
class_option(C,e) :- class(C). 

counter(o(O),Q,M,N) :- 
    option_window(O,Q),
    option_max(O,M),
    N = #sum [class_count(C,T) = T : class_option(C,O)].

counter(c(C),1,1,N) :- 
    class_count(C,N). 
    
% Link between option and class

o(P,o(O)) :- 
    class_option(C,O),
    pos(P;P+1),
    not c(c(C),P,L),
    c(c(C),P+1,L). 

:-  pos(P;P+1),
    o(P,O),
    c(O,P,K),
    not c(O,P+1,K+1).

:-  pos(P;P+1),
    not c(o(O),K,K),
    c(o(O),P+1,K), 
    not o(P+1,o(O)). 


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
%0 { c(I,P,N) } 1 :- 
%    val(I,N),
%    pos(P). 
%
%%initalisation
%
%% lower left FALSE
%:-  counter(I,_,_,_),
%    c(I,0,1).
%    
%% upper right TRUE
%:-  last(P), 
%    counter(I,_,_,S),
%    not c(I,P,S).
%
%% cumulative structure
%:-  c(I,P,N), 
%    not c(I,P+1,N), 
%    pos(P;P+1).
%
%:-  c(I,P,N+1),   
%    not c(I,P,N), 
%    val(I,N+1;N). 
%
%% propagation
%:-  counter(I,_,M,_), 
%    pos(P),
%    val(I,N;N+K),
%    not c(I,P,N),
%    K = 1..M, 
%    c(I,P+K,N+K).
%
%:-  counter(I,Q,M,_), 
%    val(I,N;N+M),
%    pos(P),
%    not c(I,P,N),
%    K = M+1..Q,
%    c(I,P+K,N+M). 
%
%%for debugging
%%#minimize[c(I,P,N):val(I,N):pos(P)].
%
%#hide.
%#show c/3. 

