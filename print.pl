writeS(A) :- write(A),write(',').

printDefOption(O) :- 
    option_window(O,W), 
    option_max(O,M),
    writeS(O),
    writeS(W),
    writeS(M),nl,
    OO is O+1, 
    (option(OO) ->
     printDefOption(OO);
     true
    ).

printDefClass(C) :- 
    class_count(C,M), 
    writeS(C),
    writeS(M),
    printOption(C,1),
    CC is C+1, 
    (class(CC) ->
     printDefClass(CC);
     true
    ).


printOverview :- 
    writeS(' '),nl,
    writeS(' '),nl,
    writeS(class),
    writeS(count), 
    printOptionHead(1),nl,
    printDefClass(0),
    writeS(' '),nl,
    writeS(option),
    writeS(window), 
    writeS(size),nl,
    printDefOption(1),
    writeS(' '),nl,
    writeS(index),
    writeS(class), 
    printOptionHead(1).

printOptionHead(N) :-
    (option(N) ->
     (   
         write('o/'),
         writeS(N),
         NN is N+1, 
         printOptionHead(NN)
     );true
    ).
    
start:-
    printOverview,nl,
    first(F), 
    printRow(F), 
    forall(class(O),printCounter(c,O)), 
    forall(option(O),printCounter(o,O)). 

is_car(P,C) :- 
    is(c,C,P). 

is_option(P,O) :- 
    is(o,O,P). 

printRow(F) :- 
    writeS(F), 
    is_car(F,C), 
    writeS(C),
    printOption(C,1), 
    (next(F,FF)->
     printRow(FF);
     true
    ).
    
printOption(C,O) :- 
    (class_option(C,O)->
     writeS('x'); 
     writeS('-')
    ),
    OO is O+1, 
    (option(OO) ->
     printOption(C,OO);
     nl
    ).

printPos(P) :- 
    writeS(P), 
    (next(P,PP) ->
        printPos(PP); 
        nl
    ). 

printCounter(T,I) :- 
    writeS(' '),nl,
    write(T), 
    write('/'),
    writeS(I), 
    printPos(0), 
    counter(T,I,_,_,N), 
    NN is N +1, 
    printCount(T,I,NN). 

printCount(T,I,N) :- 
    writeS(N), 
    printCounterRow(T,I,0,N), 
    NN is N-1, 
    (N > 0 ->
    printCount(T,I,NN); 
    nl). 

printCounterRow(T,I,P,N) :- 
    (lower(T,I,P,N) ->
     writeS('L');
        (upper(T,I,P,N) ->
         writeS('U');
            (in(T,I,P,N) ->
             writeS('1'); 
             (domain(T,I,P,N) ->
              writeS('0');
              writeS(' '))
          )
         )
     ),


    (next(P,PP)->
     printCounterRow(T,I,PP,N); 
     nl
    ).



