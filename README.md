<h1 align="center">
  One 🐋
</h1>

<p align="center">
  <strong>Utilities for <a href="https://jeffhuang.com/productivity_text_file/">One Big Text File-style</a> note taking</strong>
</p>

- 🌿 All your notes live in a single, plain-text markdown file
- 🍱 Automatic organization of daily and nested notes
- 🌍 Browse your notes with a minimal, keyboard-friendly web UI
- 🔖 Organize freely with `#tags`, `[[wiki links]]` and emoji icons
- ✨ CLI tools to list, sort, lint, and format your notes file

## Usage

Here's an example of what a valid notes file looks like:

````md
# 01.08.2026

This heading starts a new note. Since its title matches the DD.MM.YYYY format,
it's recognized as a daily note.

## Standup 📋

Level 2 headings inside a daily note become their own child notes, so you can
group different things from the same day. The first emoji found in a note (like
the one above) is used as its icon.

Discussed the roadmap for #project-x, see [[Some other note]] for the details.
Wiki links point at the note with that title, including daily notes like
[[01.08.2026]] and child notes like [[Groceries]].

## Groceries

- Milk
- Eggs
- Coffee

# Some other note

Notes don't need a date at all. Add #tags anywhere in the text to make notes
easier to find later, and use fenced code blocks for anything that shouldn't
be parsed as markdown:

```sh
echo "this won't be touched by the parser"
```
````

Run `one web --input <file>` to browse a notes file like this one in the browser, or `one --help` to see everything else `one` can do, such as sorting, linting, and formatting notes from the command line.

> [!NOTE]
> 
> If you want to format notes, [oxlint](https://oxc.rs/docs/guide/usage/linter) has to be in your `PATH`.

## Development

One is written in [Go](https://go.dev), with some tooling (formatting) managed through [pnpm](https://pnpm.io). Tasks are run with [mise](https://mise.jdx.dev):

```sh
mise run dev    # Start the web server with hot reload
mise run fmt    # Format Go code and templates
mise run test   # Run tests
```

## Credits

This app uses a number of open source packages listed in [go.mod](go.mod). Icons are from [Lucide](https://lucide.dev). It was inspired by Jeff Huang's [My productivity app is a never-ending .txt file](https://jeffhuang.com/productivity_text_file/) and Andrej Karpathy's [The Append-and-Review Note](https://karpathy.bearblog.dev/the-append-and-review-note/).

Thanks 🙏
</content>
