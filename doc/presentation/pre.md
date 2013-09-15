% SAT Encodings for the Car Sequencing Problem
% \underline{Valentin Mayer-Eichberger} and Toby Walsh 
% 16/09/2013 Doctoral Program at CP

Car Sequencing
======

\begin{multicols}{2}

\includegraphics[scale=0.23]{cars.jpg}

{\tiny Source Wikipedia }
\newpage

\begin{itemize}
\item Cars require different options (air-conditioning, sun-roof, etc.)
\item Is there a production sequence for cars on the assembly line satisfying the sliding capacity constraints?
\item CSPLib Benchmark Nr. 1
\item CP, IP, local search
\end{itemize}

\end{multicols}


Car Sequencing: Example
======

* Classes $C = \{1,2,3\}$ with demand $d_1=3, d_2=2,d_ 3=2$
* Options $O = \{a,b\}$ with capacity constraints $1/2$ and $1/5$
* Class 1 no restriction
* Class 2 requires option $a$
* Class 3 requires option $a$ and $b$

\vspace{0.5cm}

\begin{center}
\begin{small}
\begin{tabular}{c|cccccccc}
Sequence of cars  & 3 & 1 & 2 & 1 & 2 & 1 & 3 \\
\hline
Option a & 1 & - & 1 & - & 1 & - & 1 \\
Option b & 1 & - & - & - & - & - & 1 \\
\end{tabular}
\end{small} 
\end{center}     

\pause

\begin{center}

\todo{Car Sequencing is NP-Complete}
\end{center}

The SAT Approach: The Ultimate Decomposition
===============

* ONE constraint, the clause (e.g. $a \vee b \vee \neg c$).
* ONE propagator, unit propagation: (e.g. $a$ and $\neg a \vee b$ then propagate $b$).
* Using SAT solvers as blackboxes. 
* \todo{Challenge: Finding good CNF representations.}
\pause
* Global constraint: $(\sum_{i=1}^n x_{i} = d) \wedge \bigwedge_{i=0}^{n-q}(\sum_{l=1}^q x_{i+l} \leq u )$
* Use cumulative sums:

$$ s_{i,j} \iff (j \leq \sum_{l=1}^{i} x_{l}) $$

Sequential Counter
==================
\begin{center}

\begin{tikzpicture}[node distance=1cm, auto,]

\coordinate (A) at (0.5,1.1);
\coordinate (B) at (1.5,1.1);
\coordinate (C) at (0.5,-0.1);
\coordinate (D) at (1.5,-0.1);

\draw (0, 0) rectangle (1, 1);
\draw (1, 0) rectangle (2, 1);
\draw[->,thick] (A) to [bend left] (B);
\draw[->,thick] (D) to [bend left] (C);


\node at (0.5,0.5) {$s_{i-1,j}$};
\node at (1.5,0.5) {$s_{i,j}$};
%\node at (1,1.5) {test};
\node at (1,-0.5) {$\neg x_i$};

\node at (5,1) {$\neg s_{i-1,j} \vee s_{i,j}$};
\node at (5,0) {$x_{i} \vee \neg s_{i,j} \vee s_{i-1,j}$};

\end{tikzpicture}

\vspace{1cm}

\begin{tikzpicture}[node distance=1cm, auto,]

\coordinate (A) at (0.1,1.1);
\coordinate (B) at (0.9,1.9);
\coordinate (C) at (1.1,0.1);
\coordinate (D) at (1.9,0.9);

\draw (0, 0) rectangle (1, 1);
\draw (1, 1) rectangle (2, 2);
\draw[->,thick] (A) to [bend left] (B);
\draw[->,thick] (D) to [bend left] (C);

\node at (0.7,0.5) {$s_{i-1,j-1}$};
\node at (1.5,1.5) {$s_{i,j}$};
\node at (0,1.8) {$x_i$};
%\node at (2,0) {test};

\node at (5,1) {$\neg x_{i} \vee \neg s_{i-1,j-1} \vee s_{i,j}$};
\node at (5,0) {$\neg s_{i,j} \vee s_{i-1,j-1}$};

\end{tikzpicture}
\end{center}

* This idea can translate all cardinality constraints

Demand Constraint + Capacity Constraint
=======================================

$$ (\sum_{i=1}^n x_{i} = d) \wedge \bigwedge_{i=0}^{n-q}(\sum_{l=1}^q x_{i+l} \leq u )$$

\pause
\vspace{1cm}

\begin{center}

\begin{tikzpicture}[node distance=1cm, auto,]

\coordinate (A) at (0.5,1.1);
\coordinate (B) at (3.9,3.5);

\draw (0, 0) rectangle (1, 1);
\draw (4, 3) rectangle (5, 4);
\draw[->,thick] (B) to [bend right] (A);

\draw[dashed] (4.5,0.5) -- (4.5,3);
\draw[dashed] (1.5,0.5) -- (4.5,0.5);

\node at (0.7,0.5) {$s_{i-q,j-u}$};
\node at (4.5,3.5) {$s_{i,j}$};

\node at (4.1,2) {$u$};
\node at (2.5,0) {$q$};

\node at (8,2) {$\neg s_{i,j} \vee s_{i-q,j-u}$};

\end{tikzpicture}

\end{center}

Results on CSPLib
====== 

* 30+9 hard solved 28 within 20min
* Largest: 400 cars, 5 options, 23 classes: 200K Var, 1M Clauses
* Several variations of this encoding

\begin{center}
\begin{tabular}{ l|ccc}
	&E1	&E2	&E3	\\
\hline
\#solved UNSAT	& 17	&15	& 17 \\
\#fastest UNSAT	& 5	&4	& 4 \\
\#solved SAT	& 11	& 11	& 11  \\
\#fastest SAT	& 0	&4	& 7  \\
\hline
\end{tabular}
\end{center}

* With another trick 36 decision problems in the CSPLib can be solved (3 left
  open)
* Decoder and encoder available \verb+github.com/vale1410/car-sequencing+

Related Work
==================

* Sinz: Sequential Counter CNF \cite{Sinz05} 
* Een and Soerensson: Translation through BDDs to CNF \cite{Een06}
* Bacchus: Decomposition through DFAs to CNF \cite{Bacchus07}
* Brand et al: Decomposition to cumulative sums for CP \cite{Brand07}
* Siala et al: Linear time propagator for CP \cite{Siala12}

Conclusions 
======

Conclusion

* SAT is strong on instances of the CSPLib
* Global Constraints motivate for encodings
* Choosing the right encoding of cardinality constraints is crucial
* SAT can be very competitive on CP benchmarks

Current and Future work

* Fair Comparison to CP, IP, ASP, LS \ldots
* Elegant proof of GAC and lower bound on size. 
* Idea useful in rostering, planning, scheduling?
* Exponential encoding in the number of options? 
* Treat the optimization problem. 

\appendix
\newcounter{finalframe}
\setcounter{finalframe}{\value{framenumber}}

Bibliography
============

\bibliography{p}
\bibliographystyle{plain}

Backupslides
======

Capacity Constraints: Example
====================

Capacity constraint $4/8$,  demand $d=12$ on a sequence of 22 variables: 

\include{example3}

Capacity Constraints: Example
====================

Partial Assignment: $x_{1}$ and $x_{13}$ to true and $x_{12}$, $x_{14}$ and $x_{21}$ to false.

\include{example4}

Sequential Counter: Comparison to \todo{Sinzs AtMost}
==================


$\forall i \in \{1\ldots n\}$ $\forall j \in\{0 \ldots d+1\}$: 
\todo{
\begin{equation} \label{eq:1}
    \neg s_{i-1,j} \vee s_{i,j}
\end{equation}
}
\begin{equation} \label{eq:2}
    x_{i} \vee \neg s_{i,j} \vee s_{i-1,j}
\end{equation}

$\forall {i \in \{1\ldots n\}} \forall {j\in \{1\ldots d+1\}}$: 
\begin{equation} \label{eq:3}
    \neg s_{i,j} \vee s_{i-1,j-1}
\end{equation}
\todo{
\begin{equation} \label{eq:4}
    \neg x_{i} \vee \neg s_{i-1,j-1} \vee s_{i,j}
\end{equation}

\begin{equation} \label{eq:5}
     s_{0,0} \wedge \neg s_{0,1} \wedge \neg s_{n,d+1}
\end{equation}
}


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

A Trick for Lower Bounds (\cite{Gent98})
=======================

\begin{table}
\tiny
    \include{table_ian_1}
    \label{tab:2}
\end{table}

* Class 21 and 23 have option 0,1,2,4 with a total demand of 9
* All other classes share at least one option with 21 and 23
* Potential neighbours of 21 and 23?
