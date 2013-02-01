% revolutionary idea
%
% Probably the GAC encoding of the AtMostSequcenConstraint

%Input: 

%class(C) :- class_count(C,_). 
%last(N) :- N = #sum[class_count(_,K) = K ].
%pos(1..L) :- first(F), last(L). 

%     
% option_window(Option,WMax),
% option_max(Option,MMax),
% class_option(Car,O).


%Idea and Model
%the problem is given as window constraints
%of the general form 
%
% counter(I,Q,M,S), 
% with the meaning: 
% I : identifier for the counter
% Q: fixed window size
% M: maximal value per window
% S: exactly S occurences over the whole sequence

%test
last(4). 
pos(1..L) :- last(L). 
counter(t,1,2,2). 

0 { c(I,P,N) } 1 :- 
    counter(I,_,_,S), 
    N = 1..S+1, 
    pos(P). 

%initalisation

c(I,0,1) :- 
    counter(I,_,_,_). 

:-  last(P), 
    counter(I,_,_,S),
    not c(I,P,S).

:-  counter(I,_,_,S),
    c(I,P,S+1). 

% cumulative structure
:-  c(I,P,N), 
    not c(I,P+1,N). 

:-  c(I,P,N),   
    not c(I,P,N-1). 

%negative propagation
:-  counter(I,_,M,R), 
    pos(P),
    N = 1..R, 
    not c(I,P,N),
    K = 1..M, 
    c(I,P+K,N+K). 

:-  counter(I,Q,M,R), 
    pos(P),
    N = 1..R, 
    not c(I,P,N),
    K = Q-M..Q+1,
    c(I,P+K,N+M). 

%positive propagation
:-  counter(I,_,M,_), 
    c(I,P,N),
    K = 1..M, 
    not c(I,P-K,N-K). 

:-  counter(I,Q,M,_), 
    c(I,P,N),
    K = Q-M..Q+1,
    not c(I,P-K,N-M). 


