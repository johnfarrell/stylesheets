# stylesheets
Style guides and reference for use with Go, Templ, HTMX, Alpine.js, and Tailwind.css (`GHATT` stack?)

---
## About

This project serves two main goals of mine:
 - Build a library of references that I can use when working in Golang and want a specific style to use with this tech stack.
 - Provide a safe sandbox for me to try out Claude code.

## Building

This will compile any template updates and build/run the final Go binary locally. 
```shell
make run
```

The server port can be defined with the `PORT` env variable; by default it is `8080`.


Additionally, `make help` can be used to show all of the available commands. 

## AI Disclaimer

Almost the entirety of this project is created by Claude Code as a personal practice/demo project for me to play around with it.

I'm a backend engineer and don't have much experience at this time with using Templ+HTMX to build frontends and I found it difficult to find good references I can
use to learn from. I figured this is a pretty good use-case for AI code generation, and since I hadn't used that much either it was a good
time to try.

I have committed the `docs/plans` repo deliberately to show the different steps and planning phases Claude Code went through
in case anyone finds it useful.


All that being said, I have read through most of the code as it's generated and it passes the initial sniff test.