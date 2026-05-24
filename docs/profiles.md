# Writing a bank profile

A profile teaches ledgerly how to read one bank's CSV export. It is a small
YAML file that maps the bank's column headers onto ledgerly's unified fields,
plus a few hints about date and number formatting. Profiles are data, not code:
drop a `.yaml` file anywhere and pass it with `--profile ./mybank.yaml`, or open
a pull request to add it to the built-in set.

## A minimal example

```yaml
name: MyBank
delimiter: ";"        # field separator; omit to auto-detect
decimal: comma         # comma (1.234,56), dot (1,234.56) or auto
currency: EUR          # ISO 4217 code stamped on every row
date_layouts:          # Go time layouts, tried in order
  - "02.01.2006"
columns:
  date: Buchungstag
  value_date: Wertstellung
  amount: Betrag
  payee: "Empfänger"
  purpose: Verwendungszweck
```

Apply it:

```sh
ledgerly convert --profile ./mybank.yaml statement.csv
```

## Fields

| Key            | Required | Meaning                                                        |
|----------------|----------|----------------------------------------------------------------|
| `name`         | yes      | Profile name, also used by `--profile <name>` for built-ins.   |
| `delimiter`    | no       | Single character field separator. Omitted ⇒ auto-detected.     |
| `decimal`      | no       | `comma`, `dot` or `auto` (default). How to read amounts.       |
| `currency`     | no       | ISO 4217 code written into every transaction.                  |
| `date_layouts` | no       | List of Go date layouts; falls back to common defaults.        |
| `skip_rows`    | no       | Number of preamble lines above the header to drop.             |
| `signature`    | no       | Header columns that identify this bank for `--auto-detect`.    |
| `columns`      | yes      | Mapping of unified fields to source column names (see below).  |

### Columns

```yaml
columns:
  date: ...          # required: booking date column
  value_date: ...    # optional: value date (Wertstellung)
  amount: ...        # signed amount column
  debit: ...         # OR a debit column (outflow)...
  credit: ...        # ...paired with a credit column (inflow)
  payee: ...         # counterparty name
  purpose: ...       # description / Verwendungszweck
  reference: ...     # bank or end-to-end reference
```

You must map either a single signed `amount` column **or** a `debit`/`credit`
pair. With a pair, the debit side is treated as an outflow (made negative) and
the credit side as an inflow.

## Date layouts

Layouts use Go's reference date `Mon Jan 2 15:04:05 MST 2006`. The day is `02`,
the month `01`, a four-digit year `2006` and a two-digit year `06`. A few
examples:

| Sample        | Layout         |
|---------------|----------------|
| `2026-01-31`  | `2006-01-02`   |
| `31.01.2026`  | `02.01.2006`   |
| `31.01.26`    | `02.01.06`     |
| `01/31/2026`  | `01/02/2006`   |

List several layouts if a bank is inconsistent; the first one that parses wins.

## Skipping a preamble

Some banks (DKB, for example) print a few account-info lines above the real
header. `skip_rows: 4` drops them before the header is read.

## Auto-detection

If you list a `signature` — a handful of header columns unique to the bank —
ledgerly can pick the profile for you:

```sh
ledgerly convert --auto-detect statement.csv
```

The profile whose signature columns are all present in the header wins; the
most specific signature breaks ties. (Auto-detection reads the first line as the
header, so it is not used together with `skip_rows`.)

## Contributing a profile

Built-in profiles live in [`internal/profile/builtin`](../internal/profile/builtin).
To add one, drop the YAML there, add a small sample export under
`internal/profile/testdata`, wire it into `TestBuiltinProfilesGolden` and run
the tests with `-update` to generate the golden file. PRs welcome.
