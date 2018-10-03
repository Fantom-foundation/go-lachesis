## UNRELEASED

SECURITY:

FEATURES:

IMPROVEMENTS:

BUG FIXES:

## v0.3.0 (September 4, 2018)

FEATURES:

* poset: Replaced Leemon Baird's original "Fair" ordering method with 
Lamport timestamps.
* poset: Introduced the concept of Frames and Roots to enable initializing a
poset from a "non-zero" state.
* node: Added FastSync protocol to enable nodes to catch up with other nodes 
without downloading the entire poset. 
* proxy: Introduce Snapshot/Restore functionality.

IMPROVEMENTS:

* poset: Refactored the consensus methods around the concept of Frames.
* poset: Removed special case for "initial" Events, and make use of Roots 
instead. 
* docs: Added sections on Lachesis and FastSync.