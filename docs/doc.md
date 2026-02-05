# FME - Flag Matching Engine

## Pitch

FME is here to respond to a specific backend challenges : matching flag (or states) with the power of graphs.

Let's say you have eleements of your backend that can be translated as flags. It can be anything like status, parameters, items etc. 
From this set of flags, you want to apply constraints such as : I don't want flag A and flag B together. Or if I havee flag A,, I must have flag B with it.
This is what FME is doing, giving a set of flags, you'll be able to efficiently apply those kind of rules and maintains the integrity of your schema.

## Theorical example

From a set of flag called S = {a, b, c, d}
I apply the following constraints : a requires b, c interfer with d.

Giving the folllowing commbination : C = {a, b}

FME will be able to accept the combination.


## Practical example

From a set of parameters : S = {"-filepath", "-extension", "T"}
I apply the following constraints : "-filepath" requires "-extension", "T" interfer with "-extension".

Then we try the following combinations into FME engine :

c1 = {"-filepath", "-extension"} => True
c2 = {"-filepath", "-T"} => False
c3 = {"-extension", "T"} => False
