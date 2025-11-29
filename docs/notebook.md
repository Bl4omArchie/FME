# Notebook

## Introduction of flag matching problem
Assume you have a set of flags $S$ defined as :
$$ S = \{a, b, c, d, e\} $$

From this set, we want to build a combination $C$ of flags. For example :
$$ C = \{ a, c, e\}$$

We can create any combination we want, but we now introduce rules that determine whether a combination is valid.t.

For now we define two rules :

- Requires : if flag $a$ depends on flag $b$, then any combination containing $a$ is valid only if it also contains $b$. Notation : $a \implies b$
- Interference : if flag $a$ conflicts with flag $c$, then a combination cannot contain both at the same time. Notation : $a /\implies/ b$

### Algorithm 1 : Rule implementation - Require

To implement the require rule, we use a directed, unweighted graph.
Each vertex represents a flag; each edge represents a dependency.

```
function Require(a: flag, b: flag, g: Graph)
```
This algorithm will add a dependency in graph $g$ between $a$ and $b$.

Example :
```
Require("a", "b", g)
Require("b", "c", g)
```

With the example below, the graph would look like this : **a -> b -> c**

### Algorithm 2 : Rule implementation - Interfer

For interference, we use a symmetric map, which is simpler and has a lower access cost.

```
function Interfer(a: flag, b: flag, m: Map)
```

Example :
```
Interfer("a", "b", m)
Interfer("a", "c", m)
```

Resulting map : **map[a:map[b:{} c:{}] b:map[a:{}] c:map[a:{}]]**

Another example :
```
Interfer("a", "b", m)
Interfer("b", "c", m)
```
Resulting map : **map[a:map[b:{}] b:map[a:{} c:{}] c:map[b:{}]]**

### Struct 1 : Schema

Now that rules are defined, we introduce the first struct: Schema.
A schema stores all the data structures and represents the full population of flags from which combinations can be built.

```
struct Schema:
    flags: slice
    graph: directed, unweighted graph
    interferences: symmetric map
```

### Algorithm 3 : VerifySchema

Once the graph and map are filled with rules, we can verify the schema’s consistency.
```
function VerifySchema(s: Schema) -> bool, rules
```

For instance, if I set my rules like this :
$$a \implies b$$
$$b \implies c$$
$$c /\implies/ a$$

Here $a$ requires $b$ and $b$ requires $c$. But obviously the last rule is nonsens because it says $c$ is interfering with $a$ which lead to an impossible combination.
In this case `VerifySchema` will return False and the conflicting rule.
Otherwise it will retun True and a none value


### Struct 2 : Combination

A combination $C$ is a set of flags belonging to $S$.
We add a valid field indicating whether the combination satisfies all rules.

```
struct Combination:
    flags: slice
    valid: bool
```

### Algorithm 4 : VerifyCombination 

Let's write now a function that verify the correctness of a given combination.

```
function VerifyCombination(c: Combination, s: Schema) bool
```

Given the rules :
$$a \implies b$$
$$e \implies f$$
$$d /\implies/ f$$

Is this combination valid ?
$$C = \{ a, e, f, d\} $$

This combination is invalid:
- $a \implies b$ but b is missing
- $d /\implies/ f$

A valid combination could be:
$$C = \{ a, b, d\}$$



## Algorithms

There is a summary of every algorithms we saw in this notebook.

| Algorithms | Rules | Data structures |
| :--: | :--: | :--: |
| Requirement | Requires | Directed, unweighted graph |
| Interference | Interfer | Symmetric map |

## Glossary

- Flag : a simple string. It can be a foreign key or simply a word. 
- Schema : a set of flags with defined rules
- Combination : a sub-set of flags from Schema that should respect the given rules
- Conflict : incoherencies for rules or combinations
- Rules : constraints between flags that defines if a combination is valid or not
    - Requires : a dependency between two flags where b must go with a
    - Interfere : two flags cannot be mixed together into the same combination

- Flag :
- Schema :
- Combination :
- Constraint :
- Conflict :