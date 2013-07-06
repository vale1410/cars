% SAT Encodings for the Car Sequencing Problem
% Valentin Mayer-Eichberger and Toby Walsh 
% 08/07/2013 at Pragmatics of SAT

Car Sequencing
======

\includegraphics[scale=0.4]{cars.jpg}

\tiny Picture from Wikipedia

Car Sequencing from the CSPLib Benchmark
======

Assemble a production line of cars such that capacity constraints on the
workstations are not exceeded. 

Notation: 

* Set of Classes $C$
* Demand $d_i$ for class $i\in C$
* Set of Options $O$
* Capacity constraint with ratio $u_l/q_l$ for option $l\in O$
* Set $O_i \subseteq O$ that is required by class $i\in C$

Car Sequencing: Example
======

* $C = \{1,2,3\}$ with demand $3,2,2$
* $O = \{a,b\}$ with capacity constraints $1/2$ and $1/5$
* Class 1 no restriction
* Class 2 requires option $a$
* Class 3 requires option $a$ and $b$

\pause

\vspace{1cm}

\begin{center}
\begin{small}
\begin{tabular}{c|cccccccc}
  & 3 & 1 & 2 & 1 & 2 & 1 & 3 \\
\hline
a & 1 & - & 1 & - & 1 & - & 1 \\
b & 1 & - & - & - & - & - & 1 \\
\end{tabular}
\end{small} 
\end{center}     

\pause

Car Sequencing is NP-Complete

PB Model 
========

* Boolean variable $c^k_i$:  car $k\in C$ is at position $i$
* Boolean variable $o^l_i$:  option $l\in O$ is at position $i$
* Demand constraints: $\forall k \in C$ $$\sum^n_{i=1} c^k_i = d_k$$                       
* Capacity constraints: $\forall l \in O$ with ratio $u_l/q_l$ $$\bigwedge_{i=0}^{n-q_l}(\sum_{j=1}^{q_l} o^l_{i+j} \leq u_l )$$

PB Model 
========

And in all positions $i \in \{1\ldots n\}$ of the sequence it must hold:                                                    

\begin{itemize}
    \item Link between classes and options: for each $k\in C$ and 
        \begin{align*}
            \forall l \in O_k :\;\; & c^k_i - o^l_i \leq 0 \\
            \forall l \in O \setminus O_k :\;\; &c^k_i + o^l_i \leq 1\\
        \end{align*}
    \item Exactly one car:  $$\sum_{k\in C} c^k_i = 1$$  
\end{itemize}

Modelling in CNF
================

* The PB model is the mostly used model in CP,IP and local search!
* This model with standard translation to CNF (minisat+,clasp ...) has bad performance
* Choose the right cardinality translation
* More redundant constraints: $$ d_l = \sum_{i=1}^n o^l_i = \sum_{k\in C_l} d_k$$                                                 
* Global constraint: Cardinality + Sequence

Sequential Counter: Variables
==================
\begin{itemize}
    \item Translation of Boolean Cardinality: $$ \sum_{i\in \{1\ldots n\}} x_{i} = d $$ 
    \item  $x_i$ is true iff the object is at position $i$
    \item  $s_{i,j}$ is true iff in the positions $0,1 \ldots i$ the object exists at least $j$ times 
\end{itemize} 

Sequential Counter
==================


$\forall i \in \{1\ldots n\}$ $\forall j \in\{0 \ldots d+1\}$: 
\begin{equation} \label{eq:1}
    \neg s_{i-1,j} \vee s_{i,j}
\end{equation}
\begin{equation} \label{eq:2}
    x_{i} \vee \neg s_{i,j} \vee s_{i-1,j}
\end{equation}

$\forall {i \in \{1\ldots n\}} \forall {j\in \{1\ldots d+1\}}$: 
\begin{equation} \label{eq:3}
    \neg s_{i,j} \vee s_{i-1,j-1}
\end{equation}
\begin{equation} \label{eq:4}
    \neg x_{i} \vee \neg s_{i-1,j-1} \vee s_{i,j}
\end{equation}
\begin{equation} \label{eq:5}
     s_{0,0} \wedge \neg s_{0,1} \wedge s_{n,d} \wedge \neg s_{n,d+1}
\end{equation}

Sequential Counter: Example
==================

\begin{center}
\include{example1}
\end{center}
Setting $x_2$ and $x_7$ to 1:
\begin{center}
\include{example2}
\end{center}


Sequential Counter: Related Work
==================

* Carsten Sinz Sequential Counter \cite{Sinz05} 
* Fahim Bacchus translation of AMONG by the Regular constraint \cite{Bacchus07}
* Translation through BDDs \cite{Een06}

Capacity Constraints: More Global
=================================

$$ (\sum_{i=1}^n x_{i} = d) \wedge \bigwedge_{i=0}^{n-q}(\sum_{l=1}^q x_{i+l} \leq u )$$


\vspace{3cm}

$\forall {i \in \{q \ldots n\}}$, $\forall {j\in\{u\ldots d+1\}}$: 

\begin{equation} \label{eq:6}
    \neg s_{i,j} \vee s_{i-q,j-u}
\end{equation}               


Capacity Constraints: Example
====================

Capacity constraint $4/8$,  demand $d=12$ on a sequence of 22 variables: 

\include{example3}

Capacity Constraints: Example
====================

Partial Assignment: $x_{1}$ and $x_{13}$ to true and $x_{12}$, $x_{14}$ and $x_{21}$ to false.

\include{example4}


A Trick for Lower Bounds (\cite{Gent98})
=======================

\begin{table}[htbp]
\tiny
    \caption{Overview of options and demands for instance 300-04}
    \centering
    \include{table_ian_1}
    \label{tab:2}
\end{table}


Results: Solved Instances
=======

\include{all}

Conclusion and Future Work
======

* SAT can be very competitive on CP benchmarks
* SAT is very strong on proving lower bounds
* Global Constraints motivate for encodings
* Choosing the right encoding of cardinality constraints is crucial

Future work: 

* Exponential encoding in the number of options? 
* Theoretical analysis of the decompositions and usage in other domains


End
======

Thank you very much

Bibliography
============

\bibliography{p}
\bibliographystyle{plain}


SAT instances
============

\include{exp11}

UNSAT instances
============

\small

\include{exp12}

lower bounds
============

\tiny

\include{lb3}

Size
============

\tiny

\include{size}

Link between Cars and Options
=============================

$\forall i\in \{1\ldots n\}$: 

\begin{equation} \label{eq:7}
     \bigwedge_{\substack{k \in C \\ l \in O_k }} \neg c^k_{i} \vee o^l_{i}
\end{equation}

\begin{equation} \label{eq:8}
    \bigwedge_{\substack{k \in C \\ l \not \in O_k}} \neg c^k_{i} \vee \neg o^l_{i}
\end{equation}

\begin{equation} \label{eq:9}
    \bigwedge_{l\in O} \left(\neg o^l_{i} \vee \bigvee_{k \in C_l} c^k_{i}\right)
\end{equation}


Example for non GAC of E2
============

    
\begin{example}\label{ex:small}
 Let $n=5$, $d=2$ with a capacity constraint of $1/2$, and let $x_3$ be true, then
     unit propagation does not force $x_2$ nor $x_4$ to false. Setting them to true will lead to a conflict through
     clauses (\ref{eq:4}) and (\ref{eq:6}) on positions 2, 3 and 4.

\begin{center}
\begin{small}
\begin{tabular}{c|cccccccccc}
3   &   &   &   &U  &U  &U  \\
2   &   &U  &U  &.  &.  &L  \\
1   &U  &.  &.  &L  &L  &   \\
0   &L  &L  &L  &   &   &   \\
\hline
$s_{i,j}$ &0  &1  &2  &3  &4  &5 \\
$x_i$     &  &.  &.  &1  &.  &.  \\
\end{tabular}
\end{small} 
\end{center}     
\end{example}

