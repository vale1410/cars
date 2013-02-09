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
    printRow(F). 

option(O) :- 
    is_option(_,O).

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
