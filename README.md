Different SAT models for the car sequencing problem
====

Modelling attempts for the car sequencing problem from the CSPLIB

Description of the models
-------------------------

SAT Models: (subfolder gen)

e1 - e5) see paper description

ASP Models (subfolder asp)

1) ASP naive (for each window and option (for each window and option

2) as 1) but one abstraction (variable for option)

3) as 2) but with symmetry + redundant

7) as 3) but all cardinalites replaced by explicit counters

8) attempt to translate to a pure graphical representation [failed]

9) attempt to model exclusively by counters that should archive GAC on
the AtMostSeqCar constraint [to be continued]

10) Model that enforces GAC on AtMostSeqCard with FL (failedLiteral) and
is so far the best model in comparison. There are translations to SAT
and a paper discribing the encoding (ex/p.pdf)
