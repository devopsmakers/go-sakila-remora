# go-sakila-remora
A sidecar process written in Go to monitor MySQL slaves / masters for replication errors or lag.

## Welcome!

Hi! and thanks for stopping by this project in it's early days.

Inspired by GitHub's interesting article:
https://githubengineering.com/context-aware-mysql-pools-via-haproxy/

It's similar to methods I've used previously to manage traffic being sent
to Asterisk VoIP servers based on several key metrics.

## Project goals

Ultimately, I'm pretty new to Go and thought this would be an interesting
challenge and an awesome product if everything turns out well.

### Phase 1:

The initial iteration will handle reading YAML configuration, executing some
basic checks against MySQL and presenting a JSON based endpoint for health
checking against. After this iteration I should have a decent BVP - Barely
Viable Product.

### Phase 2:

Let's not get ahead of ourselves.
