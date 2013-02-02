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

% For each Option we need a total or of
% classes and values for the channeling

%option_succ(Option,C1,N1,C2,N2). 

class_value(C,V) :- 
    class_count(C,N), 
    V = 1..N. 

option_first(O,C,1) :-  
    option(O), 
    C = #min[class_option(X,O) = X].

option_succ(O,C,V1,C,V1+1) :- 
    class_option(C,O), 
    class_value(C,V1;V1+1). 

option_succ(O,C1,V,C2,1) :- 
    class_option(C1,O), 
    class_option(C2,O), 
    C1 < C2,
    class_value(C1,V), 
    not class_value(C1,V+1), 
    { class_option(C3,O):C3=C1+1..(C2-1) } 0. 

option_node(O,C,V) :- 
    option_first(O,C,V). 

option_node(O,C,V) :- 
    option_succ(O,_,_,C,V). 

option_last(O,C,N) :- 
    option_node(O,C,N), 
    { option_succ(O,C,N,_,_) } 0. 

class_option(C,e) :- class(C). 

% e stands for exactly
% Exactly one car per position
% which is modelled by aux option
option_window(e,1). 
option_max(e,1). 
class_option(C,e) :- class(C). 

%counter(c(C),1,1,N) :- 
%    class_count(C,N). 
%
%counter(o(O),Q,M,N) :- 
%    option_window(O,Q),
%    option_max(O,M),
%    N = #sum [class_count(C,T) = T : class_option(C,O)].

% Now we need to connect
% con(Option,Position,Counter,Value). 




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

