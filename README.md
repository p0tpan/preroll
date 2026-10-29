## Grammar
Valid dice rolls are defined by an [EBNF](https://en.wikipedia.org/wiki/Extended_Backus%E2%80%93Naur_form) specification which can be found in the [docs](docs/preroll.iso-ebnf). You can view the railroad diagram [here](docs/dice_roll_ebnf_railroad_solid_bg.png).

### Examples
Here are some valid roll inputs:
- `1d20`
- `13d3`
- `2d5 - 4`
- `5d{red,blue,green}`
- `2d{1,3,17}`
- `kh2d20 + 3`
- `dl3d13 - 4`

## Generating tables
`preroll gen` rolls an expression `-n` times and prints the results:

```
preroll gen -roll "1d20+3" -n 100 -seed 42 -format csv|json|md
```

Output is deterministic for a given expression, count and seed. If `-seed` is omitted a random one is chosen and printed in the output header (a top-level `seed` field in JSON, a `# roll=… seed=… n=…` comment line in CSV, a header list in Markdown), so any table can be regenerated.

`kh`/`kl`/`dh`/`dl` keep or drop exactly one die, so `kh2d20` is "advantage" and `dl4d6` drops the lowest of four. `!{N}` allows at most N extra rolls per die; a bare `!` allows up to 100. Dropped dice appear in parentheses in CSV and Markdown output.
