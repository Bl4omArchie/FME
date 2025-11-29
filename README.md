# FME - Flag Matching Engine

Flag Matching Engine, or FME, is a top-layer program based on `lvlath` graph library that verify constraints between a set of flags.

This program intend to solve backend challenges such as efficient argument parsing.

To learn more about FME, please see the notebook I've made : [click here](docs/notebook.md)


# Roadmap

v0 : this pre-release hold the fundamental aspect of the engine, the release (v1) will consist on something more reliable and production ready.
- Rules : requirement and interference
- Validation : validate a schema and a combination
- Custom error messages : ErrSchemaCycle, ErrSchemaContradiction, ErrCombinationInterfer
- Tests : basic conflict with requirement and interference + dependency cycle

Updates :
- Engine interface : Schema is now implementing the Engine interface for mor eliberty in the implementation of constraints
- PathInstance : new class to keep track on specific depedencies, ideal for error messages
- Instant validation in order to keep the schema integrity
