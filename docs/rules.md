# Categorization rules

Ledgerly assigns a category to each transaction by running it through an ordered
list of rules. The **first rule that matches wins**, so put specific rules above
broad ones. Rules live in a YAML file; the built-in default set is in
[`internal/categorize/rules.yaml`](../internal/categorize/rules.yaml).

## Rule shape

```yaml
rules:
  - { match: "REWE", category: "Groceries" }
  - { regex: "amazon.*(prime|video)", category: "Subscriptions" }
  - { match: "Gehalt", field: purpose, category: "Income" }
```

| Key        | Meaning                                                                 |
|------------|-------------------------------------------------------------------------|
| `match`    | Case-insensitive substring to look for.                                 |
| `regex`    | Regular expression (case-insensitive). Use instead of `match`.          |
| `field`    | Where to look: `payee`, `purpose` or `any` (default — both, joined).    |
| `category` | The category to assign. Required. Any string; the built-ins are a guide.|

Each rule needs exactly one of `match` or `regex`, and a `category`.

## Precedence and fallback

- Rules are evaluated top to bottom; the first match decides the category.
- A transaction that matches no rule is categorized as `Uncategorized`.
- With `--rules myfile.yaml`, your rules are checked **before** the defaults, so
  you can override how a merchant is categorized without restating the whole
  list. Everything you don't cover still falls through to the defaults.

```sh
ledgerly convert --profile ING --rules my-rules.yaml statement.csv
```

## Built-in categories

The default rules map to these buckets:

`Income`, `Groceries`, `Dining Out`, `Rent`, `Utilities`, `Subscriptions`,
`Transport`, `Shopping`, `Health`, `Insurance`, `Fees`, `Cash`, and the
`Uncategorized` fallback.

You are not limited to these — a rule can emit any category string, and it will
show up in the CSV/JSON output and the terminal summary as-is.

## Tips

- Order matters: list `Amazon Prime` (Subscriptions) above a bare `Amazon`
  (Shopping), otherwise the broad rule swallows it.
- Scope with `field: purpose` when a word is too generic to match on the payee
  (e.g. `Gehalt` in the purpose line).
- Prefer `match` for plain merchant names; reach for `regex` only when you need
  alternation or word boundaries.
