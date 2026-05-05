# envlint

> Validate `.env` files against a schema definition and catch missing or malformed variables before deployment.

---

## Installation

```bash
go install github.com/yourusername/envlint@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envlint.git
cd envlint && go build -o envlint .
```

---

## Usage

Define a schema file (`.env.schema`) describing your expected variables:

```ini
DATABASE_URL=required,url
PORT=required,number
DEBUG=optional,bool
APP_SECRET=required,min:32
```

Then run `envlint` against your `.env` file:

```bash
envlint --schema .env.schema --env .env
```

**Example output:**

```
✔  DATABASE_URL   valid
✘  PORT           missing (required)
✔  DEBUG          valid
✘  APP_SECRET     too short (min: 32 characters)

2 error(s) found. Fix before deploying.
```

Exit code is `0` on success and `1` when validation fails, making it easy to integrate into CI pipelines.

---

## CI Integration

```yaml
- name: Lint environment
  run: envlint --schema .env.schema --env .env.ci
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)