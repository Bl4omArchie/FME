## New rules

We now introduce additional rules.
The first three relate to positional constraints, which matter in systems like OTO where a binary can have positional arguments.
The last introduces a scale of intensity.

- Behind to : if flag $a$ is BehindTo flag $b$, a combination containing $a$ must be positioned behind to $b$ if $b$ appears.
- In front of if flag $a$ is InFrontOf flag $b$, a combination containing $a$ must be positioned in front of $b$ if $b$ appears.
- Position at $i$: a flag $a$ must appear at index $i$ (0-based).
- Scale at : some flags express the same idea at different intensities; we group them and assign weights (0–10, for example).

### Struct 1 - Modification

In those new cases, what data structure would could use ?

For `BehindTo`, `InFrontOf` and `Position at`, my personal proposition is to use either a doubly linked list or a sorted slice.

A DLL can be simplier for modifying only flags that are behind or in front of a given node but with only a sorted slice if could make thing easier and we'll then need one data structure for three rules.

To support the new rules:
```
struct Schema:
    flags: slice
    graph: directed, unweighted graph
    conflicts: symmetric map
    ddl: doubly linked list
    position: []FlagID
    scales: []Scale
```


### Struct 3 : Scale

For the scale rule, we use a weighted map: (map[int]FlagID)

```
struct Scale:
    name: string
    description: string
    values: map[FlagID]int
```

## Algorithms

There is a summary of every algorithms we saw in this notebook.

| Algorithms | Rules | Data structures |
| :--: | :--: | :--: |
| Requirement | Requires | Directed, unweighted graph |
| Interference | Interfer | Symmetric map |
| Position | Behind to & In front of| Doubly linked list |
| Position at| Place at $N$ | Ordered slice |
| Scale at | Scale 0 .. $N$ | Map|


## Conclusion 

More ideas.
- each rule implement its own verifcation function and ValidateCombination() only task is to call those functions.
- Removes ValidateSchema : each new added rule is instantly verified so the validation can be instant and automatic.
- constraints visualization : export the `core.Graph` to a DOT/GraphML representation for documentation or debugging.
- group rules
- optimization : is there an algorithm that can optimize the rules ? use less data structure, reduce computation.

## Glossary

- Flag : a simple string that represent your data.
- Schema : actual engine that handle the constraints between your flags.
- Combination : a set of flags respecting the rules from the Schema.
- Conflict : incorrect rule or combination.
- Rules : constraints between flags that defines if a combination is valid or not
    - Requires : a dependency between two flags where b must go with a
    - Interferes : two flags cannot be mixed together into the same combination
    - Position: relative ordering constraints.
    - Position At: fixed index.
    - Scale At: weighted intensity grouping.
