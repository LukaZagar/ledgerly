# Ledgerly

Every bank exports its CSV differently — different column names, different date
formats, comma vs. dot as the decimal separator, debit and credit split across
two columns or folded into one signed amount. The moment you try to merge
several accounts or feed them into a budgeting tool, you drown in that mess.

**Ledgerly turns those wildly different bank CSVs into one clean, unified and
automatically categorized transaction list** — as a single, dependency-light
command-line binary. No bank login, no cloud: it only ever touches the files you
hand it.

### Before → after

A raw ING export:

```
Buchung;Valuta;Auftraggeber/Empfänger;Buchungstext;Verwendungszweck;Betrag;Währung
02.01.2026;02.01.2026;REWE SAGT DANKE;Lastschrift;Einkauf Lebensmittel;-43,12;EUR
```

After `ledgerly convert`:

```
date,value_date,amount,currency,payee,purpose,account,category,reference
2026-01-02,2026-01-02,-43.12,EUR,REWE SAGT DANKE,Einkauf Lebensmittel,ing,Groceries,
```

Same row, but now with an ISO date, a dot-decimal amount, a stable column order
and a category — ready to merge with any other account.

## Features

- **Bank profiles** (YAML): map a bank's columns onto one unified schema,
  including date and number formats. Ships with ING, DKB, Sparkasse and PayPal.
- **Rule-based categorization** ("REWE"/"EDEKA" → Groceries, "Netflix" →
  Subscriptions), with a sensible default rule set you can extend or override.
- **Merge multiple accounts** into a single chronological timeline.
- **Duplicate detection** for overlapping export ranges.
- Output as unified **CSV or JSON**.
- A short **terminal summary** (totals per category and per month).

## Installation

With a Go toolchain (1.22+):

```sh
go install github.com/SciTee/ledgerly/cmd/ledgerly@latest
```

Or grab a prebuilt binary from the [releases page](https://github.com/SciTee/ledgerly/releases)
and drop it on your `PATH`. Ledgerly is a single static binary with no runtime
dependencies.

## Quickstart

Convert one statement, auto-detecting the bank from its header:

```sh
ledgerly convert --auto-detect statement.csv > unified.csv
```

Or name the profile explicitly and write JSON:

```sh
ledgerly convert --profile ING --format json -o january.json ing-january.csv
```

See which profiles are built in:

```sh
ledgerly profiles list
```

## Merging multiple accounts

Pass several files; ledgerly reads each with the chosen profile, stamps the
account label (the file name by default), merges everything into one timeline
and drops duplicate rows from overlapping exports:

```sh
ledgerly convert --profile Sparkasse giro-jan.csv giro-feb.csv > giro.csv
```

Use `--account` to set an explicit label instead of the file name.

> Tip: convert each bank separately (different profiles), then concatenate the
> unified CSVs — they all share the same schema.

## Categorization rules

Categories come from an ordered list of rules; the first match wins. The
built-in set covers common German merchants. To add your own, write a rules
file and stack it on top of the defaults with `--rules`:

```yaml
# my-rules.yaml — checked before the defaults
rules:
  - { match: "Lieblingsbäcker", category: "Dining Out" }
  - { regex: "miete|kaltmiete", field: purpose, category: "Rent" }
```

```sh
ledgerly convert --profile DKB --rules my-rules.yaml statement.csv
```

A rule matches a case-insensitive substring (`match`) or a regular expression
(`regex`) against the `payee`, `purpose`, or both (`field: any`, the default).
Anything no rule matches is left as `Uncategorized`. Pass `--no-categorize` to
skip the step entirely.

## Output format

The unified schema has a fixed column order in CSV and the same fields in JSON:

| Field        | Notes                                           |
|--------------|-------------------------------------------------|
| `date`       | Booking date, ISO `YYYY-MM-DD`.                 |
| `value_date` | Value date (Wertstellung), if the bank has one. |
| `amount`     | Signed; negative = outflow. Dot decimal.        |
| `currency`   | ISO 4217 code from the profile.                 |
| `payee`      | Counterparty.                                   |
| `purpose`    | Description / Verwendungszweck.                 |
| `account`    | Source account label.                           |
| `category`   | Assigned category or `Uncategorized`.           |
| `reference`  | Bank / end-to-end reference, if present.        |

## Bank profiles

Built in: **ING**, **DKB**, **Sparkasse**, **PayPal**. A profile is just a small
YAML file, so you can point at your own with `--profile ./mybank.yaml`. See
[docs/profiles.md](docs/profiles.md) for the full schema and a walkthrough.

## Contributing

Adding a profile for your bank is the most useful contribution. Drop the YAML in
[`internal/profile/builtin`](internal/profile/builtin), add a small sample
export and a golden test, and open a PR — details in
[docs/profiles.md](docs/profiles.md).

## Non-goals

- No bank API / online-banking login — strictly file-based.
- Not a full accounting or budgeting app — Ledgerly is the clean preprocessing
  step in front of one.

## License

MIT — see [LICENSE](LICENSE).
