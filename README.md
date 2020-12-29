# TODO-CLI

A cli tool for tracking todo items.

## Design

* CLI boilerplate: `./cmd` module
* Business logic: `./notebook` module

## Build

`go build -o todo`

## Run

#### Add a new todo item "call mum"

Input: `./todo add "call mum"`

Output: `Item 1 added` 

#### Complete a todo item

Input: `./todo done <itemId>`

Output `Item <itemId> done.`

#### List all items

Input: `./todo list --all`

Output: 
```
1. call mum
2. [Done] call dad
Total: 2 items, 1 items done
```
#### List undone items

Input: `./todo list`

Output: 
```
1. call mum
Total: 1 items
```